package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jeremy-kr/ccfg/internal/model"
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
		}

		files = append(files, cf)
	}
	return files
}

// scanDir은 디렉토리 내의 파일들을 스캔한다.
func scanDir(dir string, scope model.Scope, category model.ConfigCategory) []model.ConfigFile {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	var children []model.ConfigFile
	for _, entry := range entries {
		if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		absPath := filepath.Join(dir, entry.Name())
		info, err := entry.Info()
		if err != nil {
			continue
		}
		children = append(children, model.ConfigFile{
			Path:        absPath,
			Scope:       scope,
			FileType:    detectFileType(entry.Name()),
			Category:    category,
			Exists:      true,
			Size:        info.Size(),
			ModTime:     info.ModTime(),
			Description: entry.Name(),
		})
	}
	return children
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
