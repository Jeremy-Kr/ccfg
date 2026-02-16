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

// Scanner는 Claude Code 설정 파일을 탐색한다.
type Scanner struct {
	// WorkDir은 프로젝트 루트 감지의 시작점이다. 비어있으면 CWD 사용.
	WorkDir string
}

// New는 새로운 Scanner를 생성한다.
func New(workDir string) *Scanner {
	return &Scanner{WorkDir: workDir}
}

// Scan은 모든 Scope의 설정 파일을 탐색하고 메타데이터를 수집한다.
func (s *Scanner) Scan() (*model.ScanResult, error) {
	workDir := s.WorkDir
	if workDir == "" {
		wd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("작업 디렉토리 조회 실패: %w", err)
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

// scanEntries는 FileEntry 목록을 순회하며 각 파일의 메타데이터를 수집한다.
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

			// 디렉토리면 하위 파일을 스캔
			if e.IsDir && info.IsDir() {
				cf.Children = scanDir(absPath, scope, e.Category)
			}

			// settings.json → hooks + mcpServers 가상 Children
			if cf.Category == model.CategorySettings && !cf.IsDir {
				cf.Children = parseSettingsSections(absPath, scope)
			}

			// .mcp.json → 서버 목록 가상 Children
			if cf.Category == model.CategoryMCP && !cf.IsDir {
				cf.Children = parseMCPSections(absPath, scope)
			}
		}

		files = append(files, cf)
	}
	return files
}

// scanDir은 디렉토리 내의 파일들을 스캔한다.
// 심볼릭 링크가 디렉토리를 가리키는 경우도 포함한다.
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

		// 심볼릭 링크 해석
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

		// 하위 디렉토리 재귀 스캔
		if info.IsDir() {
			cf.Children = scanDir(absPath, scope, category)
		}

		children = append(children, cf)
	}
	return children
}

// parseSettingsSections는 settings.json에서 hooks와 mcpServers를 파싱하여 가상 Children을 생성한다.
func parseSettingsSections(path string, scope model.Scope) []model.ConfigFile {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	raw := string(data)

	var children []model.ConfigFile

	// hooks 파싱
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

	// mcpServers 파싱
	if mcpGroup := buildMCPServerGroup(path, scope, raw); mcpGroup != nil {
		children = append(children, *mcpGroup)
	}

	return children
}

// parseMCPSections는 .mcp.json에서 서버 목록을 파싱하여 가상 Children을 생성한다.
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

// buildMCPServerGroup은 JSON 원본에서 mcpServers를 파싱하여 가상 그룹 노드를 생성한다.
// 서버가 없으면 nil을 반환한다.
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

// detectFileType은 파일 경로의 확장자로 FileType을 판별한다.
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
