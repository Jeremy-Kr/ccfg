package usage

import (
	"testing"
)

func TestRank_Empty(t *testing.T) {
	entries := Rank(nil)
	if entries != nil {
		t.Errorf("expected nil, got %v", entries)
	}

	entries = Rank(map[string]int{})
	if entries != nil {
		t.Errorf("expected nil for empty map, got %v", entries)
	}
}

func TestRank_SingleItem(t *testing.T) {
	entries := Rank(map[string]int{"explore": 100})
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Grade != GradeSSS {
		t.Errorf("single item should be SSS, got %s", entries[0].Grade)
	}
	if entries[0].Name != "explore" {
		t.Errorf("expected name 'explore', got %s", entries[0].Name)
	}
}

func TestRank_SortedDescending(t *testing.T) {
	counts := map[string]int{
		"alpha":   10,
		"bravo":   100,
		"charlie": 50,
	}
	entries := Rank(counts)

	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[0].Name != "bravo" || entries[1].Name != "charlie" || entries[2].Name != "alpha" {
		t.Errorf("unexpected order: %s, %s, %s", entries[0].Name, entries[1].Name, entries[2].Name)
	}
}

func TestRank_GradeThresholds(t *testing.T) {
	// max=1000, 각 항목의 logScore가 임계값 경계에 위치하도록 설계
	counts := map[string]int{
		"max":  1000, // logScore = 1.0 → SSS
		"low":  1,    // logScore ≈ 0.1 → D or F
		"zero": 0,    // logScore = 0.0 → F (log(1)/log(1001) = 0)
	}
	entries := Rank(counts)

	gradeMap := make(map[string]Grade)
	for _, e := range entries {
		gradeMap[e.Name] = e.Grade
	}

	if gradeMap["max"] != GradeSSS {
		t.Errorf("max should be SSS, got %s", gradeMap["max"])
	}
	if gradeMap["zero"] != GradeF {
		t.Errorf("zero count should be F, got %s", gradeMap["zero"])
	}
}

func TestGradeFromScore(t *testing.T) {
	tests := []struct {
		score float64
		want  Grade
	}{
		{1.0, GradeSSS},
		{0.95, GradeSSS},
		{0.94, GradeSS},
		{0.80, GradeSS},
		{0.79, GradeS},
		{0.65, GradeS},
		{0.64, GradeA},
		{0.50, GradeA},
		{0.49, GradeB},
		{0.35, GradeB},
		{0.34, GradeC},
		{0.20, GradeC},
		{0.19, GradeD},
		{0.10, GradeD},
		{0.09, GradeF},
		{0.0, GradeF},
	}

	for _, tt := range tests {
		got := gradeFromScore(tt.score)
		if got != tt.want {
			t.Errorf("gradeFromScore(%.2f) = %s, want %s", tt.score, got, tt.want)
		}
	}
}

func TestGrade_String(t *testing.T) {
	if GradeSSS.String() != "SSS" {
		t.Errorf("GradeSSS.String() = %s", GradeSSS.String())
	}
	if GradeF.String() != "F" {
		t.Errorf("GradeF.String() = %s", GradeF.String())
	}
}
