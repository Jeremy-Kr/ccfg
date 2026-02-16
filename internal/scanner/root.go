package scanner

import (
	"os"
	"path/filepath"
)

// FindProjectRoot walks up from startDir looking for a .git directory.
// It returns the directory containing .git, or an empty string if none is found.
func FindProjectRoot(startDir string) string {
	dir := startDir
	for {
		if isGitRoot(dir) {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached filesystem root without finding .git
			return ""
		}
		dir = parent
	}
}

func isGitRoot(dir string) bool {
	info, err := os.Stat(filepath.Join(dir, ".git"))
	return err == nil && info.IsDir()
}
