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

// opencodeLine parses an opencode-format transcript JSONL line.
// Format: {"type":"tool_use", "tool_name":"task", "tool_input":{...}}
type opencodeLine struct {
	ToolName  string          `json:"tool_name"`
	ToolInput json.RawMessage `json:"tool_input"`
}

// claudeCodeLine parses a Claude Code-format transcript JSONL line.
// Format: {"type":"assistant", "message":{"content":[{"type":"tool_use","name":"Task","input":{...}}]}}
type claudeCodeLine struct {
	Type    string `json:"type"`
	Message struct {
		Content json.RawMessage `json:"content"`
	} `json:"message"`
}

// contentBlock represents an element in Claude Code's content array.
type contentBlock struct {
	Type  string          `json:"type"`
	Name  string          `json:"name"`
	Input json.RawMessage `json:"input"`
}

// agentInput extracts the subagent_type from task/delegate_task tool_input.
type agentInput struct {
	SubagentType string `json:"subagent_type"`
	Description  string `json:"description"`
}

// collectAgents tallies agent invocations from transcripts.
func collectAgents(homeDir, projectFilter string) (map[string]int, error) {
	return collectFromTranscripts(homeDir, projectFilter, extractAgent)
}

func extractAgent(line []byte) (name string, ok bool) {
	if !bytes.Contains(line, []byte(`subagent_type`)) {
		return "", false
	}

	// opencode format
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

	// Claude Code format
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

// toolAsAgent lists built-in infrastructure subagents excluded from agent rankings.
// These overlap with tool names or are too generic to distinguish from custom agents.
var toolAsAgent = map[string]bool{
	"Bash": true, "bash": true,
	"Explore": true, "explore": true,
}

// resolveAgentName converts a subagent type into a meaningful name.
// Returns an empty string for agents that overlap with tool names (excluded).
// Extracts custom agent names from the description of general-purpose agents.
func resolveAgentName(input agentInput) string {
	if toolAsAgent[input.SubagentType] {
		return ""
	}

	if input.SubagentType != "general-purpose" || input.Description == "" {
		return input.SubagentType
	}

	// "horner: evaluate code" → "horner"
	idx := strings.Index(input.Description, ":")
	if idx <= 0 || idx > 30 {
		return input.SubagentType
	}

	candidate := strings.TrimSpace(input.Description[:idx])
	// Names with spaces are not agent names (exclude one-off labels like "Agent #1: xxx")
	if strings.Contains(candidate, " ") {
		return input.SubagentType
	}
	return candidate
}

// extractFromClaudeCode iterates over tool_use blocks in a Claude Code assistant message.
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

// transcriptDirs returns the list of transcript directories to scan.
func transcriptDirs(homeDir, projectFilter string) []string {
	if projectFilter != "" {
		encoded := encodeProjectPath(projectFilter)
		return []string{filepath.Join(homeDir, ".claude", "projects", encoded)}
	}

	// All scope: transcripts/ + all subdirectories under projects/
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

// encodeProjectPath converts a project path to Claude's directory encoding format.
// Example: /Users/jeremy/code/ccfg → -Users-jeremy-code-ccfg
func encodeProjectPath(path string) string {
	return strings.ReplaceAll(path, "/", "-")
}

type extractFunc func(line []byte) (name string, ok bool)

// collectFromTranscripts collects data from transcript files using the given extract function.
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

// scanTranscripts scans JSONL files in a directory and extracts data using the extract function.
func scanTranscripts(dir string, counts map[string]int, extract extractFunc) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read transcript directory: %w", err)
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
