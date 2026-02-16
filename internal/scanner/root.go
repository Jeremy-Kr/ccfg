package scanner

import (
	"os"
	"path/filepath"
)

// FindProjectRoot는 startDir부터 상위로 올라가며 .git 디렉토리를 찾는다.
// 찾으면 해당 디렉토리를, 못 찾으면 빈 문자열을 반환한다.
func FindProjectRoot(startDir string) string {
	dir := startDir
	for {
		if isGitRoot(dir) {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			// 루트까지 올라감, 찾지 못함
			return ""
		}
		dir = parent
	}
}

func isGitRoot(dir string) bool {
	info, err := os.Stat(filepath.Join(dir, ".git"))
	return err == nil && info.IsDir()
}
