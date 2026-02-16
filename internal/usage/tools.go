package usage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// sessionMeta는 session-meta JSON의 필요한 필드만 파싱한다.
type sessionMeta struct {
	ProjectPath string         `json:"project_path"`
	ToolCounts  map[string]int `json:"tool_counts"`
}

// collectTools는 session-meta에서 도구별 호출 횟수를 집계한다.
func collectTools(homeDir, projectFilter string) (map[string]int, error) {
	dir := filepath.Join(homeDir, ".claude", "usage-data", "session-meta")
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("session-meta 디렉토리 읽기 실패: %w", err)
	}

	counts := make(map[string]int)
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			continue
		}
		var meta sessionMeta
		if err := json.Unmarshal(data, &meta); err != nil {
			continue
		}
		if projectFilter != "" && meta.ProjectPath != projectFilter {
			continue
		}
		for tool, count := range meta.ToolCounts {
			counts[tool] += count
		}
	}
	return counts, nil
}
