package usage

import (
	"bytes"
	"encoding/json"
)

// skillInput은 skill tool_input에서 name을 추출한다.
type skillInput struct {
	Name string `json:"name"`
}

// collectSkills는 transcript에서 스킬별 호출 횟수를 집계한다.
func collectSkills(homeDir, projectFilter string) (map[string]int, error) {
	return collectFromTranscripts(homeDir, projectFilter, extractSkill)
}

func extractSkill(line []byte) (name string, ok bool) {
	if !bytes.Contains(line, []byte(`"tool_name":"skill"`)) &&
		!bytes.Contains(line, []byte(`"tool_name": "skill"`)) {
		return "", false
	}
	var tl transcriptLine
	if err := json.Unmarshal(line, &tl); err != nil {
		return "", false
	}
	if tl.ToolName != "skill" {
		return "", false
	}
	var input skillInput
	if err := json.Unmarshal(tl.ToolInput, &input); err != nil || input.Name == "" {
		return "", false
	}
	return input.Name, true
}
