package usage

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// sessionMeta parses only the required fields from a session-meta JSON file.
type sessionMeta struct {
	ProjectPath string         `json:"project_path"`
	ToolCounts  map[string]int `json:"tool_counts"`
}

// collectTools tallies tool invocation counts from both session-meta and transcripts.
func collectTools(homeDir, projectFilter string) (map[string]int, error) {
	counts := make(map[string]int)

	// 1) Collect from session-meta
	if err := collectToolsFromSessionMeta(homeDir, projectFilter, counts); err != nil {
		return nil, err
	}

	// 2) Collect from transcripts (a single line may contain multiple tool_use blocks)
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
		return fmt.Errorf("failed to read session-meta directory: %w", err)
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

// collectToolsFromTranscripts extracts tool usage directly from transcript files.
// Uses a dedicated scanner instead of extractFunc because a single line may contain multiple tool_use blocks.
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

// extractToolsFromLine extracts all tool names from a single line and adds them to counts.
func extractToolsFromLine(line []byte, counts map[string]int) {
	if !bytes.Contains(line, []byte(`"tool_use"`)) {
		return
	}

	// opencode format: {"type":"tool_use", "tool_name":"Read"}
	if bytes.Contains(line, []byte(`"tool_name"`)) {
		var ol opencodeLine
		if err := json.Unmarshal(line, &ol); err == nil && ol.ToolName != "" {
			counts[ol.ToolName]++
			return
		}
	}

	// Claude Code format: tool_use blocks inside an assistant message
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
