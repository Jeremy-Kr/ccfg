package usage

import (
	"bytes"
	"encoding/json"
	"strings"
)

// skillInput extracts the name from a skill tool_input.
type skillInput struct {
	Name  string `json:"name"`
	Skill string `json:"skill"` // Claude Code format
}

// collectSkills tallies skill invocations from transcripts.
func collectSkills(homeDir, projectFilter string) (map[string]int, error) {
	return collectFromTranscripts(homeDir, projectFilter, extractSkill)
}

func extractSkill(line []byte) (name string, ok bool) {
	if !bytes.Contains(line, []byte(`skill`)) {
		return "", false
	}

	// opencode format: {"tool_name":"skill", "tool_input":{"name":"git-master"}}
	if bytes.Contains(line, []byte(`"tool_name"`)) {
		var ol opencodeLine
		if err := json.Unmarshal(line, &ol); err == nil && ol.ToolName == "skill" {
			var input skillInput
			if err := json.Unmarshal(ol.ToolInput, &input); err == nil && input.Name != "" {
				return input.Name, true
			}
		}
	}

	// Claude Code format: {"type":"assistant","message":{"content":[{"type":"tool_use","name":"Skill","input":{"skill":"commit"}}]}}
	if bytes.Contains(line, []byte(`"assistant"`)) {
		return extractFromClaudeCode(line, func(block contentBlock) (string, bool) {
			if !strings.EqualFold(block.Name, "skill") {
				return "", false
			}
			var input skillInput
			if err := json.Unmarshal(block.Input, &input); err == nil {
				// Claude Code uses the "skill" field, opencode uses the "name" field
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
