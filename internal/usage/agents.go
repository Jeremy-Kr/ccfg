package usage

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// opencodeLine는 opencode 형식의 transcript JSONL을 파싱한다.
// 형식: {"type":"tool_use", "tool_name":"task", "tool_input":{...}}
type opencodeLine struct {
	ToolName  string          `json:"tool_name"`
	ToolInput json.RawMessage `json:"tool_input"`
}

// claudeCodeLine는 Claude Code 형식의 transcript JSONL을 파싱한다.
// 형식: {"type":"assistant", "message":{"content":[{"type":"tool_use","name":"Task","input":{...}}]}}
type claudeCodeLine struct {
	Type    string `json:"type"`
	Message struct {
		Content json.RawMessage `json:"content"`
	} `json:"message"`
}

// contentBlock는 Claude Code의 content 배열 요소이다.
type contentBlock struct {
	Type  string          `json:"type"`
	Name  string          `json:"name"`
	Input json.RawMessage `json:"input"`
}

// agentInput은 task/delegate_task의 tool_input에서 subagent_type을 추출한다.
type agentInput struct {
	SubagentType string `json:"subagent_type"`
	Description  string `json:"description"`
}

// collectAgents는 transcript에서 에이전트별 호출 횟수를 집계한다.
func collectAgents(homeDir, projectFilter string) (map[string]int, error) {
	return collectFromTranscripts(homeDir, projectFilter, extractAgent)
}

func extractAgent(line []byte) (name string, ok bool) {
	if !bytes.Contains(line, []byte(`subagent_type`)) {
		return "", false
	}

	// opencode 형식
	if bytes.Contains(line, []byte(`"tool_name"`)) {
		var ol opencodeLine
		if err := json.Unmarshal(line, &ol); err == nil {
			if ol.ToolName == "task" || ol.ToolName == "delegate_task" {
				var input agentInput
				if err := json.Unmarshal(ol.ToolInput, &input); err == nil && input.SubagentType != "" {
					if resolved := resolveAgentName(input); resolved != "" {
						return resolved, true
					}
				}
			}
		}
	}

	// Claude Code 형식
	if bytes.Contains(line, []byte(`"assistant"`)) {
		return extractFromClaudeCode(line, func(block contentBlock) (string, bool) {
			if !strings.EqualFold(block.Name, "task") {
				return "", false
			}
			var input agentInput
			if err := json.Unmarshal(block.Input, &input); err == nil && input.SubagentType != "" {
				if resolved := resolveAgentName(input); resolved != "" {
					return resolved, true
				}
			}
			return "", false
		})
	}

	return "", false
}

// toolAsAgent는 빌트인 인프라 서브에이전트로, 에이전트 랭킹에서 제외한다.
// 도구와 이름이 겹치거나 범용적이라 커스텀 에이전트와 구분이 어려운 타입들.
var toolAsAgent = map[string]bool{
	"Bash": true, "bash": true,
	"Explore": true, "explore": true,
}

// resolveAgentName은 서브에이전트 타입을 의미 있는 이름으로 변환한다.
// - 도구명과 겹치는 에이전트는 빈 문자열 반환 (제외)
// - general-purpose의 description에서 커스텀 에이전트 이름 추출
func resolveAgentName(input agentInput) string {
	if toolAsAgent[input.SubagentType] {
		return ""
	}

	if input.SubagentType != "general-purpose" || input.Description == "" {
		return input.SubagentType
	}

	// "horner: coding-guidelines로 평가" → "horner"
	idx := strings.Index(input.Description, ":")
	if idx <= 0 || idx > 30 {
		return input.SubagentType
	}

	candidate := strings.TrimSpace(input.Description[:idx])
	// 공백이 포함되면 이름이 아님 ("Agent #1: xxx" 같은 일회성 라벨 제외)
	if strings.Contains(candidate, " ") {
		return input.SubagentType
	}
	return candidate
}

// extractFromClaudeCode는 Claude Code 형식의 assistant 메시지에서 tool_use 블록을 순회한다.
func extractFromClaudeCode(line []byte, match func(contentBlock) (string, bool)) (string, bool) {
	var cl claudeCodeLine
	if err := json.Unmarshal(line, &cl); err != nil || cl.Type != "assistant" {
		return "", false
	}

	var blocks []contentBlock
	if err := json.Unmarshal(cl.Message.Content, &blocks); err != nil {
		return "", false
	}

	for _, block := range blocks {
		if block.Type != "tool_use" {
			continue
		}
		if name, ok := match(block); ok {
			return name, true
		}
	}
	return "", false
}

// transcriptDirs는 스캔 대상 transcript 디렉토리 목록을 반환한다.
func transcriptDirs(homeDir, projectFilter string) []string {
	if projectFilter != "" {
		encoded := encodeProjectPath(projectFilter)
		return []string{filepath.Join(homeDir, ".claude", "projects", encoded)}
	}

	// 전체 범위: transcripts/ + projects/의 모든 하위 디렉토리
	dirs := []string{filepath.Join(homeDir, ".claude", "transcripts")}

	projectsBase := filepath.Join(homeDir, ".claude", "projects")
	entries, err := os.ReadDir(projectsBase)
	if err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				dirs = append(dirs, filepath.Join(projectsBase, entry.Name()))
			}
		}
	}
	return dirs
}

// encodeProjectPath는 프로젝트 경로를 Claude의 디렉토리 인코딩 형식으로 변환한다.
// 예: /Users/jeremy/code/ccfg → -Users-jeremy-code-ccfg
func encodeProjectPath(path string) string {
	return strings.ReplaceAll(path, "/", "-")
}

type extractFunc func(line []byte) (name string, ok bool)

// collectFromTranscripts는 transcript 파일들에서 extract 함수로 데이터를 수집한다.
func collectFromTranscripts(homeDir, projectFilter string, extract extractFunc) (map[string]int, error) {
	dirs := transcriptDirs(homeDir, projectFilter)
	counts := make(map[string]int)
	for _, dir := range dirs {
		if err := scanTranscripts(dir, counts, extract); err != nil {
			continue
		}
	}
	return counts, nil
}

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
