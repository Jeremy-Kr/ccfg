package usage

import (
	"fmt"
	"strings"
)

// Collector는 Claude Code 사용 데이터를 수집한다.
type Collector struct {
	HomeDir     string // 사용자 홈 디렉토리
	ProjectPath string // 현재 프로젝트 경로 (빈 문자열이면 프로젝트 필터링 안 함)
}

// Collect는 지정된 범위의 사용 데이터를 수집하고 등급을 산정한다.
func (c *Collector) Collect(scope DataScope) (*UsageData, error) {
	projectFilter := ""
	if scope == ScopeProject && c.ProjectPath != "" {
		projectFilter = c.ProjectPath
	}

	toolCounts, err := collectTools(c.HomeDir, projectFilter)
	if err != nil {
		return nil, fmt.Errorf("도구 수집 실패: %w", err)
	}

	agentCounts, err := collectAgents(c.HomeDir, projectFilter)
	if err != nil {
		return nil, fmt.Errorf("에이전트 수집 실패: %w", err)
	}

	skillCounts, err := collectSkills(c.HomeDir, projectFilter)
	if err != nil {
		return nil, fmt.Errorf("스킬 수집 실패: %w", err)
	}

	return &UsageData{
		Tools:  Rank(normalizeCounts(toolCounts)),
		Agents: Rank(normalizeCounts(agentCounts)),
		Skills: Rank(normalizeCounts(skillCounts)),
	}, nil
}

// normalizeMap은 opencode 소문자 이름을 Claude Code PascalCase로 매핑한다.
var normalizeMap = map[string]string{
	// 도구
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
	// 에이전트
	"explore": "Explore",
	"plan":    "Plan",
}

// normalizeCounts는 동일한 도구/에이전트의 대소문자 변형을 통합한다.
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
