package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindProjectRoot(t *testing.T) {
	// 임시 디렉토리에 .git 구조 생성
	tmp := t.TempDir()
	projectDir := filepath.Join(tmp, "myproject")
	subDir := filepath.Join(projectDir, "src", "pkg")

	if err := os.MkdirAll(filepath.Join(projectDir, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(subDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// 하위 디렉토리에서 프로젝트 루트를 찾아야 함
	root := FindProjectRoot(subDir)
	if root != projectDir {
		t.Errorf("FindProjectRoot(%q) = %q, want %q", subDir, root, projectDir)
	}

	// 프로젝트 루트에서 자기 자신을 반환
	root = FindProjectRoot(projectDir)
	if root != projectDir {
		t.Errorf("FindProjectRoot(%q) = %q, want %q", projectDir, root, projectDir)
	}

	// .git이 없는 디렉토리에서는 빈 문자열
	noGitDir := filepath.Join(tmp, "nogit")
	if err := os.MkdirAll(noGitDir, 0o755); err != nil {
		t.Fatal(err)
	}
	root = FindProjectRoot(noGitDir)
	if root != "" {
		t.Errorf("FindProjectRoot(%q) = %q, want empty", noGitDir, root)
	}
}
