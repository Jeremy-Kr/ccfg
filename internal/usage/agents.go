package usage

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// transcriptLine는 transcript JSONL에서 필요한 필드만 파싱한다.
type transcriptLine struct {
	ToolName  string          `json:"tool_name"`
	ToolInput json.RawMessage `json:"tool_input"`
}

// agentInput은 task/delegate_task의 tool_input에서 subagent_type을 추출한다.
type agentInput struct {
	SubagentType string `json:"subagent_type"`
}

// collectAgents는 transcript에서 에이전트별 호출 횟수를 집계한다.
func collectAgents(homeDir, projectFilter string) (map[string]int, error) {
	dirs := transcriptDirs(homeDir, projectFilter)
	counts := make(map[string]int)

	for _, dir := range dirs {
		if err := scanTranscripts(dir, counts, extractAgent); err != nil {
			continue
		}
	}
	return counts, nil
}

func extractAgent(line []byte) (name string, ok bool) {
	if !bytes.Contains(line, []byte(`"subagent_type"`)) {
		return "", false
	}
	var tl transcriptLine
	if err := json.Unmarshal(line, &tl); err != nil {
		return "", false
	}
	if tl.ToolName != "task" && tl.ToolName != "delegate_task" {
		return "", false
	}
	var input agentInput
	if err := json.Unmarshal(tl.ToolInput, &input); err != nil || input.SubagentType == "" {
		return "", false
	}
	return input.SubagentType, true
}

// transcriptDirs는 스캔 대상 transcript 디렉토리 목록을 반환한다.
func transcriptDirs(homeDir, projectFilter string) []string {
	if projectFilter != "" {
		// 프로젝트별 transcript: ~/.claude/projects/{encoded-path}/*.jsonl
		encoded := encodeProjectPath(projectFilter)
		return []string{filepath.Join(homeDir, ".claude", "projects", encoded)}
	}
	// 전체: ~/.claude/transcripts/
	return []string{filepath.Join(homeDir, ".claude", "transcripts")}
}

// encodeProjectPath는 프로젝트 경로를 Claude의 디렉토리 인코딩 형식으로 변환한다.
// 예: /Users/jeremy/code/ccfg → -Users-jeremy-code-ccfg
func encodeProjectPath(path string) string {
	result := make([]byte, 0, len(path))
	for _, c := range []byte(path) {
		if c == '/' {
			result = append(result, '-')
		} else {
			result = append(result, c)
		}
	}
	return string(result)
}

type extractFunc func(line []byte) (name string, ok bool)

// scanTranscripts는 디렉토리 내 JSONL 파일을 스캔하며 extract 함수로 데이터를 추출한다.
func scanTranscripts(dir string, counts map[string]int, extract extractFunc) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("transcript 디렉토리 읽기 실패: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".jsonl" {
			continue
		}
		if err := scanFile(filepath.Join(dir, entry.Name()), counts, extract); err != nil {
			continue
		}
	}
	return nil
}

func scanFile(path string, counts map[string]int, extract extractFunc) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		if name, ok := extract(scanner.Bytes()); ok {
			counts[name]++
		}
	}
	return scanner.Err()
}
