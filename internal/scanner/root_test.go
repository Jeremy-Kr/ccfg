package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindProjectRoot(t *testing.T) {
	// Create .git structure in temp directory
	tmp := t.TempDir()
	projectDir := filepath.Join(tmp, "myproject")
	subDir := filepath.Join(projectDir, "src", "pkg")

	if err := os.MkdirAll(filepath.Join(projectDir, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(subDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Should find project root from a subdirectory
	root := FindProjectRoot(subDir)
	if root != projectDir {
		t.Errorf("FindProjectRoot(%q) = %q, want %q", subDir, root, projectDir)
	}

	// Should return itself when called from the project root
	root = FindProjectRoot(projectDir)
	if root != projectDir {
		t.Errorf("FindProjectRoot(%q) = %q, want %q", projectDir, root, projectDir)
	}

	// Should return empty string when no .git directory exists
	noGitDir := filepath.Join(tmp, "nogit")
	if err := os.MkdirAll(noGitDir, 0o755); err != nil {
		t.Fatal(err)
	}
	root = FindProjectRoot(noGitDir)
	if root != "" {
		t.Errorf("FindProjectRoot(%q) = %q, want empty", noGitDir, root)
	}
}
