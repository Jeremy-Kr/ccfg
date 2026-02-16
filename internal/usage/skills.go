package usage

import (
	"bytes"
	"encoding/json"
	"strings"
)

// skillInput은 skill tool_input에서 name을 추출한다.
type skillInput struct {
	Name  string `json:"name"`
	Skill string `json:"skill"` // Claude Code 형식
}

// collectSkills는 transcript에서 스킬별 호출 횟수를 집계한다.
func collectSkills(homeDir, projectFilter string) (map[string]int, error) {
	return collectFromTranscripts(homeDir, projectFilter, extractSkill)
}

func extractSkill(line []byte) (name string, ok bool) {
	if !bytes.Contains(line, []byte(`skill`)) {
		return "", false
	}

	// opencode 형식: {"tool_name":"skill", "tool_input":{"name":"git-master"}}
	if bytes.Contains(line, []byte(`"tool_name"`)) {
		var ol opencodeLine
		if err := json.Unmarshal(line, &ol); err == nil && ol.ToolName == "skill" {
			var input skillInput
			if err := json.Unmarshal(ol.ToolInput, &input); err == nil && input.Name != "" {
				return input.Name, true
			}
		}
	}

	// Claude Code 형식: {"type":"assistant","message":{"content":[{"type":"tool_use","name":"Skill","input":{"skill":"commit"}}]}}
	if bytes.Contains(line, []byte(`"assistant"`)) {
		return extractFromClaudeCode(line, func(block contentBlock) (string, bool) {
			if !strings.EqualFold(block.Name, "skill") {
				return "", false
			}
			var input skillInput
			if err := json.Unmarshal(block.Input, &input); err == nil {
				// Claude Code는 "skill" 필드, opencode는 "name" 필드
				if input.Skill != "" {
					return input.Skill, true
				}
				if input.Name != "" {
					return input.Name, true
				}
			}
			return "", false
		})
	}

	return "", false
}
