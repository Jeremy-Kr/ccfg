package usage

import "time"

// Grade represents a usage frequency tier from SSS to F.
type Grade int

const (
	GradeSSS Grade = iota
	GradeSS
	GradeS
	GradeA
	GradeB
	GradeC
	GradeD
	GradeF
)

func (g Grade) String() string {
	switch g {
	case GradeSSS:
		return "SSS"
	case GradeSS:
		return "SS"
	case GradeS:
		return "S"
	case GradeA:
		return "A"
	case GradeB:
		return "B"
	case GradeC:
		return "C"
	case GradeD:
		return "D"
	case GradeF:
		return "F"
	default:
		return "?"
	}
}

// RankCategory represents the type of ranking.
type RankCategory int

const (
	RankAgents RankCategory = iota
	RankTools
	RankSkills
)

func (r RankCategory) String() string {
	switch r {
	case RankAgents:
		return "Agents"
	case RankTools:
		return "Tools"
	case RankSkills:
		return "Skills"
	default:
		return "Unknown"
	}
}

// DataScope represents the data collection scope.
type DataScope int

const (
	ScopeAll     DataScope = iota // All projects
	ScopeProject                  // Current project only
)

func (d DataScope) String() string {
	if d == ScopeProject {
		return "Project"
	}
	return "All"
}

// TimePeriod represents a time range filter for usage data.
type TimePeriod int

const (
	PeriodAll   TimePeriod = iota // 전체 기간
	PeriodMonth                   // 최근 30일
	PeriodWeek                    // 최근 7일
	PeriodDay                     // 최근 24시간

	timePeriodCount = iota // number of TimePeriod values (must remain last)
)

func (p TimePeriod) String() string {
	switch p {
	case PeriodMonth:
		return "30d"
	case PeriodWeek:
		return "7d"
	case PeriodDay:
		return "24h"
	default:
		return "All"
	}
}

// Next returns the next period in the cycle: All → 30d → 7d → 24h → All.
func (p TimePeriod) Next() TimePeriod {
	return (p + 1) % timePeriodCount
}

// Cutoff returns the cutoff time for this period. Zero value means no filtering.
func (p TimePeriod) Cutoff() time.Time {
	now := time.Now()
	switch p {
	case PeriodMonth:
		return now.Add(-30 * 24 * time.Hour)
	case PeriodWeek:
		return now.Add(-7 * 24 * time.Hour)
	case PeriodDay:
		return now.Add(-24 * time.Hour)
	default:
		return time.Time{}
	}
}

// RankEntry represents a single item in a ranking list.
type RankEntry struct {
	Name     string
	Count    int
	Grade    Grade
	LogScore float64
}

// UsageData holds the collected usage data for all categories.
type UsageData struct {
	Agents []RankEntry
	Tools  []RankEntry
	Skills []RankEntry
}
