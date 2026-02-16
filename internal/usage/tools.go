package usage

import (
	"bufio"
	"bytes"
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

// collectTools는 session-meta와 transcript 양쪽에서 도구별 호출 횟수를 집계한다.
func collectTools(homeDir, projectFilter string) (map[string]int, error) {
	counts := make(map[string]int)

	// 1) session-meta에서 수집
	if err := collectToolsFromSessionMeta(homeDir, projectFilter, counts); err != nil {
		return nil, err
	}

	// 2) transcript에서 수집 (한 줄에 여러 tool_use가 있을 수 있어 별도 스캐너 사용)
	if err := collectToolsFromTranscripts(homeDir, projectFilter, counts); err != nil {
		return nil, err
	}

	return counts, nil
}

func collectToolsFromSessionMeta(homeDir, projectFilter string, counts map[string]int) error {
	dir := filepath.Join(homeDir, ".claude", "usage-data", "session-meta")
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("session-meta 디렉토리 읽기 실패: %w", err)
	}

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
	return nil
}

// collectToolsFromTranscripts는 transcript 파일들에서 도구 사용을 직접 추출한다.
// 한 줄에 여러 tool_use 블록이 있을 수 있어 extractFunc 대신 직접 스캔한다.
func collectToolsFromTranscripts(homeDir, projectFilter string, counts map[string]int) error {
	dirs := transcriptDirs(homeDir, projectFilter)
	for _, dir := range dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if entry.IsDir() || filepath.Ext(entry.Name()) != ".jsonl" {
				continue
			}
			scanFileMultiTool(filepath.Join(dir, entry.Name()), counts)
		}
	}
	return nil
}

func scanFileMultiTool(path string, counts map[string]int) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		extractToolsFromLine(scanner.Bytes(), counts)
	}
}

// extractToolsFromLine는 한 줄에서 모든 도구 이름을 추출하여 counts에 추가한다.
func extractToolsFromLine(line []byte, counts map[string]int) {
	if !bytes.Contains(line, []byte(`"tool_use"`)) {
		return
	}

	// opencode 형식: {"type":"tool_use", "tool_name":"Read"}
	if bytes.Contains(line, []byte(`"tool_name"`)) {
		var ol opencodeLine
		if err := json.Unmarshal(line, &ol); err == nil && ol.ToolName != "" {
			counts[ol.ToolName]++
			return
		}
	}

	// Claude Code 형식: assistant 메시지 안의 tool_use 블록들
	if !bytes.Contains(line, []byte(`"assistant"`)) {
		return
	}

	var cl claudeCodeLine
	if err := json.Unmarshal(line, &cl); err != nil || cl.Type != "assistant" {
		return
	}

	var blocks []contentBlock
	if err := json.Unmarshal(cl.Message.Content, &blocks); err != nil {
		return
	}

	for _, block := range blocks {
		if block.Type == "tool_use" && block.Name != "" {
			counts[block.Name]++
		}
	}
}
