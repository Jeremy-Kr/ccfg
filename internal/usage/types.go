package usage

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
