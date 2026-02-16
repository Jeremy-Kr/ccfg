package usage

import "fmt"

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
		Tools:  Rank(toolCounts),
		Agents: Rank(agentCounts),
		Skills: Rank(skillCounts),
	}, nil
}
