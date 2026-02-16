package usage

// Grade는 사용 빈도에 따른 등급을 나타낸다 (SSS ~ F).
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

// Emoji는 등급에 해당하는 이모지를 반환한다.
func (g Grade) Emoji() string {
	switch g {
	case GradeSSS:
		return "🏆"
	case GradeSS:
		return "⭐"
	case GradeS:
		return "✨"
	case GradeA:
		return "🔥"
	case GradeB:
		return "💎"
	case GradeC:
		return "🌿"
	case GradeD:
		return "💤"
	case GradeF:
		return "🪨"
	default:
		return "❓"
	}
}

// RankCategory는 랭킹의 종류를 나타낸다.
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

// Emoji는 카테고리에 해당하는 이모지를 반환한다.
func (r RankCategory) Emoji() string {
	switch r {
	case RankAgents:
		return "🤖"
	case RankTools:
		return "🔧"
	case RankSkills:
		return "🧠"
	default:
		return "❓"
	}
}

// DataScope는 데이터 범위를 나타낸다.
type DataScope int

const (
	ScopeAll     DataScope = iota // 모든 프로젝트
	ScopeProject                  // 현재 프로젝트만
)

func (d DataScope) String() string {
	if d == ScopeProject {
		return "Project"
	}
	return "All"
}

// RankEntry는 랭킹 목록의 한 항목이다.
type RankEntry struct {
	Name     string
	Count    int
	Grade    Grade
	LogScore float64
}

// UsageData는 수집된 전체 사용 데이터이다.
type UsageData struct {
	Agents []RankEntry
	Tools  []RankEntry
	Skills []RankEntry
}
