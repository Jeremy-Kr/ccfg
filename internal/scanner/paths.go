package scanner

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/jeremy-kr/ccfg/internal/model"
)

// WatchPaths collects file and directory paths across all scopes for fsnotify watching.
// For files, it adds both the parent directory (to detect creation/deletion) and the file
// itself (to detect content changes, if it exists). Directories are added as-is.
func WatchPaths(projectRoot string) []string {
	seen := make(map[string]bool)
	var paths []string

	add := func(p string) {
		if seen[p] {
			return
		}
		seen[p] = true
		paths = append(paths, p)
	}

	collect := func(base string, entries []FileEntry) {
		if base == "" {
			return
		}
		for _, e := range entries {
			abs := filepath.Join(base, e.RelPath)
			if e.IsDir {
				add(abs)
			} else {
				// Parent directory (to detect file creation/deletion)
				add(filepath.Dir(abs))
				// File itself (to detect content changes, only if it exists)
				if _, err := os.Stat(abs); err == nil {
					add(abs)
				}
			}
		}
	}

	if base, entries := ManagedPaths(); base != "" {
		collect(base, entries)
	}
	if base, entries := UserPaths(); base != "" {
		collect(base, entries)
	}
	if base, entries := ProjectPaths(projectRoot); base != "" {
		collect(base, entries)
	}

	return paths
}

// FileEntry defines a single file to scan.
type FileEntry struct {
	RelPath     string               // Relative path from the base directory.
	Description string               // Human-readable description.
	Category    model.ConfigCategory // Functional category.
	IsDir       bool                 // Whether to scan as a directory.
}

// GetUserHomeDir returns the current user's home directory.
func GetUserHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return home
}

// ManagedPaths returns paths for system-administered configuration files.
func ManagedPaths() (string, []FileEntry) {
	var base string
	switch runtime.GOOS {
	case "darwin":
		base = "/Library/Application Support/ClaudeCode"
	case "linux":
		base = "/etc/claude-code"
	default:
		return "", nil
	}

	return base, []FileEntry{
		{RelPath: "managed_settings.json", Description: "Managed settings", Category: model.CategorySettings},
		{RelPath: "policies.json", Description: "Policy file", Category: model.CategoryPolicy},
	}
}

// UserPaths returns paths for user-level global configuration files.
func UserPaths() (string, []FileEntry) {
	home := GetUserHomeDir()
	if home == "" {
		return "", nil
	}

	return home, []FileEntry{
		{RelPath: filepath.Join(".claude", "settings.json"), Description: "User global settings", Category: model.CategorySettings},
		{RelPath: filepath.Join(".claude", "settings.local.json"), Description: "User local settings", Category: model.CategorySettings},
		{RelPath: ".claude.json", Description: "Legacy global settings", Category: model.CategorySettings},
		{RelPath: filepath.Join(".claude", "CLAUDE.md"), Description: "User global instructions", Category: model.CategoryInstructions},
		{RelPath: ".mcp.json", Description: "MCP server global settings", Category: model.CategoryMCP},
		{RelPath: filepath.Join(".claude", "commands"), Description: "Custom commands", Category: model.CategoryCommands, IsDir: true},
		{RelPath: filepath.Join(".claude", "skills"), Description: "Agent skills", Category: model.CategorySkills, IsDir: true},
		{RelPath: filepath.Join(".claude", "agents"), Description: "Custom agents", Category: model.CategoryAgents, IsDir: true},
		{RelPath: filepath.Join(".claude", "keybindings.json"), Description: "Keybinding settings", Category: model.CategoryKeybindings},
	}
}

// ProjectPaths returns paths for project-level configuration files.
func ProjectPaths(root string) (string, []FileEntry) {
	if root == "" {
		return "", nil
	}

	return root, []FileEntry{
		{RelPath: filepath.Join(".claude", "settings.json"), Description: "Project settings", Category: model.CategorySettings},
		{RelPath: filepath.Join(".claude", "settings.local.json"), Description: "Project local settings", Category: model.CategorySettings},
		{RelPath: "CLAUDE.md", Description: "Project instructions", Category: model.CategoryInstructions},
		{RelPath: filepath.Join(".claude", "CLAUDE.md"), Description: "Project instructions (alternate location)", Category: model.CategoryInstructions},
		{RelPath: ".mcp.json", Description: "MCP server project settings", Category: model.CategoryMCP},
		{RelPath: filepath.Join(".claude", "commands"), Description: "Project commands", Category: model.CategoryCommands, IsDir: true},
		{RelPath: filepath.Join(".claude", "skills"), Description: "Project skills", Category: model.CategorySkills, IsDir: true},
		{RelPath: filepath.Join(".claude", "agents"), Description: "Project agents", Category: model.CategoryAgents, IsDir: true},
	}
}
