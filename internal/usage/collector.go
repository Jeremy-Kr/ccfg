package usage

import (
	"fmt"
	"strings"
)

// Collector collects Claude Code usage data.
type Collector struct {
	HomeDir     string // User home directory
	ProjectPath string // Current project path (empty string disables project filtering)
}

// Collect gathers usage data for the given scope and assigns grades.
func (c *Collector) Collect(scope DataScope) (*UsageData, error) {
	projectFilter := ""
	if scope == ScopeProject && c.ProjectPath != "" {
		projectFilter = c.ProjectPath
	}

	toolCounts, err := collectTools(c.HomeDir, projectFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to collect tools: %w", err)
	}

	agentCounts, err := collectAgents(c.HomeDir, projectFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to collect agents: %w", err)
	}

	skillCounts, err := collectSkills(c.HomeDir, projectFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to collect skills: %w", err)
	}

	return &UsageData{
		Tools:  Rank(normalizeCounts(toolCounts)),
		Agents: Rank(normalizeCounts(agentCounts)),
		Skills: Rank(normalizeCounts(skillCounts)),
	}, nil
}

// normalizeMap maps opencode lowercase names to Claude Code PascalCase names.
var normalizeMap = map[string]string{
	// Tools
	"bash":         "Bash",
	"read":         "Read",
	"edit":         "Edit",
	"write":        "Write",
	"glob":         "Glob",
	"grep":         "Grep",
	"task":         "Task",
	"skill":        "Skill",
	"websearch":    "WebSearch",
	"todowrite":    "TodoWrite",
	"todoread":     "TodoRead",
	"question":     "AskUserQuestion",
	"slashcommand": "Skill",
	// Agents
	"explore": "Explore",
	"plan":    "Plan",
}

// normalizeCounts merges case variants of the same tool or agent into a single canonical name.
func normalizeCounts(counts map[string]int) map[string]int {
	if len(counts) == 0 {
		return counts
	}

	result := make(map[string]int, len(counts))
	for name, count := range counts {
		canonical := name
		if mapped, ok := normalizeMap[name]; ok {
			canonical = mapped
		} else if mapped, ok := normalizeMap[strings.ToLower(name)]; ok {
			canonical = mapped
		}
		result[canonical] += count
	}
	return result
}
