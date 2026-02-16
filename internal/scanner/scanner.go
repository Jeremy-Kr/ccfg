package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jeremy-kr/ccfg/internal/model"
	"github.com/jeremy-kr/ccfg/internal/parser"
)

// Scanner discovers Claude Code configuration files.
type Scanner struct {
	// WorkDir is the starting directory for project root detection. If empty, CWD is used.
	WorkDir string
}

// New creates a new Scanner.
func New(workDir string) *Scanner {
	return &Scanner{WorkDir: workDir}
}

// Scan discovers configuration files across all scopes and collects their metadata.
func (s *Scanner) Scan() (*model.ScanResult, error) {
	workDir := s.WorkDir
	if workDir == "" {
		wd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get working directory: %w", err)
		}
		workDir = wd
	}

	result := &model.ScanResult{}

	// Managed scope
	if base, entries := ManagedPaths(); base != "" {
		result.Managed = scanEntries(base, entries, model.ScopeManaged)
	}

	// User scope
	if base, entries := UserPaths(); base != "" {
		result.User = scanEntries(base, entries, model.ScopeUser)
	}

	// Project scope
	rootDir := FindProjectRoot(workDir)
	result.RootDir = rootDir
	if rootDir != "" {
		if base, entries := ProjectPaths(rootDir); base != "" {
			result.Project = scanEntries(base, entries, model.ScopeProject)
		}
	}

	return result, nil
}

// scanEntries iterates over FileEntry items and collects metadata for each file.
func scanEntries(base string, entries []FileEntry, scope model.Scope) []model.ConfigFile {
	files := make([]model.ConfigFile, 0, len(entries))
	for _, e := range entries {
		absPath := filepath.Join(base, e.RelPath)
		cf := model.ConfigFile{
			Path:        absPath,
			Scope:       scope,
			FileType:    detectFileType(e.RelPath),
			Category:    e.Category,
			Description: e.Description,
		}

		if info, err := os.Stat(absPath); err == nil {
			cf.Exists = true
			cf.Size = info.Size()
			cf.ModTime = info.ModTime()
			cf.IsDir = info.IsDir()

			// Scan children if it is a directory
			if e.IsDir && info.IsDir() {
				cf.Children = scanDir(absPath, scope, e.Category)
			}

			// settings.json -> virtual children for hooks + mcpServers
			if cf.Category == model.CategorySettings && !cf.IsDir {
				cf.Children = parseSettingsSections(absPath, scope)
			}

			// .mcp.json -> virtual children for server list
			if cf.Category == model.CategoryMCP && !cf.IsDir {
				cf.Children = parseMCPSections(absPath, scope)
			}
		}

		files = append(files, cf)
	}
	return files
}

// scanDir scans files within a directory.
// It also follows symbolic links that point to directories.
func scanDir(dir string, scope model.Scope, category model.ConfigCategory) []model.ConfigFile {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	var children []model.ConfigFile
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		absPath := filepath.Join(dir, entry.Name())

		// Resolve symbolic links
		info, err := os.Stat(absPath)
		if err != nil {
			continue
		}

		cf := model.ConfigFile{
			Path:        absPath,
			Scope:       scope,
			FileType:    detectFileType(entry.Name()),
			Category:    category,
			Exists:      true,
			IsDir:       info.IsDir(),
			Size:        info.Size(),
			ModTime:     info.ModTime(),
			Description: entry.Name(),
		}

		// Recursively scan subdirectories
		if info.IsDir() {
			cf.Children = scanDir(absPath, scope, category)
		}

		children = append(children, cf)
	}
	return children
}

// parseSettingsSections parses hooks and mcpServers from settings.json and creates virtual children.
func parseSettingsSections(path string, scope model.Scope) []model.ConfigFile {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	raw := string(data)

	var children []model.ConfigFile

	// Parse hooks
	hooks := parser.ParseSettingsHooks(raw)
	if len(hooks) > 0 {
		sort.Slice(hooks, func(i, j int) bool {
			return hooks[i].Event < hooks[j].Event
		})
		var hookChildren []model.ConfigFile
		for _, h := range hooks {
			hookChildren = append(hookChildren, model.ConfigFile{
				Path:        path + "#hooks." + h.Event,
				Scope:       scope,
				FileType:    model.FileTypeJSON,
				Category:    model.CategoryHooks,
				Exists:      true,
				IsVirtual:   true,
				Description: h.Event,
			})
		}
		children = append(children, model.ConfigFile{
			Path:        path + "#hooks",
			Scope:       scope,
			FileType:    model.FileTypeJSON,
			Category:    model.CategoryHooks,
			Exists:      true,
			IsDir:       true,
			IsVirtual:   true,
			Description: fmt.Sprintf("Hooks (%d)", len(hooks)),
			Children:    hookChildren,
		})
	}

	// Parse mcpServers
	if mcpGroup := buildMCPServerGroup(path, scope, raw); mcpGroup != nil {
		children = append(children, *mcpGroup)
	}

	return children
}

// parseMCPSections parses the server list from .mcp.json and creates virtual children.
func parseMCPSections(path string, scope model.Scope) []model.ConfigFile {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	group := buildMCPServerGroup(path, scope, string(data))
	if group == nil {
		return nil
	}
	return group.Children
}

// buildMCPServerGroup parses mcpServers from raw JSON and creates a virtual group node.
// It returns nil if no servers are found.
func buildMCPServerGroup(path string, scope model.Scope, raw string) *model.ConfigFile {
	servers := parser.ParseMCPServers(raw)
	if len(servers) == 0 {
		return nil
	}

	sort.Slice(servers, func(i, j int) bool {
		return servers[i].Name < servers[j].Name
	})

	children := make([]model.ConfigFile, 0, len(servers))
	for _, s := range servers {
		children = append(children, model.ConfigFile{
			Path:        path + "#mcpServers." + s.Name,
			Scope:       scope,
			FileType:    model.FileTypeJSON,
			Category:    model.CategoryMCP,
			Exists:      true,
			IsVirtual:   true,
			Description: s.Name,
		})
	}

	return &model.ConfigFile{
		Path:        path + "#mcpServers",
		Scope:       scope,
		FileType:    model.FileTypeJSON,
		Category:    model.CategoryMCP,
		Exists:      true,
		IsDir:       true,
		IsVirtual:   true,
		Description: fmt.Sprintf("MCP Servers (%d)", len(servers)),
		Children:    children,
	}
}

// detectFileType determines the FileType based on the file path extension.
func detectFileType(path string) model.FileType {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".md":
		return model.FileTypeMarkdown
	case ".jsonc":
		return model.FileTypeJSONC
	case ".json":
		return model.FileTypeJSON
	default:
		return model.FileTypeJSON
	}
}
