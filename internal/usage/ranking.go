package usage

import (
	"math"
	"sort"
)

// Rank assigns grades on a log scale to invocation counts and sorts entries in descending order.
func Rank(counts map[string]int) []RankEntry {
	if len(counts) == 0 {
		return nil
	}

	maxCount := 0
	for _, c := range counts {
		if c > maxCount {
			maxCount = c
		}
	}

	entries := make([]RankEntry, 0, len(counts))
	for name, count := range counts {
		score := logScore(count, maxCount)
		entries = append(entries, RankEntry{
			Name:     name,
			Count:    count,
			Grade:    gradeFromScore(score),
			LogScore: score,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Count != entries[j].Count {
			return entries[i].Count > entries[j].Count
		}
		return entries[i].Name < entries[j].Name
	})

	return entries
}

func logScore(count, maxCount int) float64 {
	if maxCount <= 0 {
		return 0
	}
	return math.Log(float64(count)+1) / math.Log(float64(maxCount)+1)
}

func gradeFromScore(score float64) Grade {
	switch {
	case score >= 0.95:
		return GradeSSS
	case score >= 0.80:
		return GradeSS
	case score >= 0.65:
		return GradeS
	case score >= 0.50:
		return GradeA
	case score >= 0.35:
		return GradeB
	case score >= 0.20:
		return GradeC
	case score >= 0.10:
		return GradeD
	default:
		return GradeF
	}
}
