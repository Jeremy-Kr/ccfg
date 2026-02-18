package usage

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
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
	// max=1000, each item's logScore is designed to sit near grade threshold boundaries
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

func TestTimePeriod_String(t *testing.T) {
	tests := []struct {
		p    TimePeriod
		want string
	}{
		{PeriodAll, "All"},
		{PeriodMonth, "30d"},
		{PeriodWeek, "7d"},
		{PeriodDay, "24h"},
	}
	for _, tt := range tests {
		if got := tt.p.String(); got != tt.want {
			t.Errorf("TimePeriod(%d).String() = %q, want %q", tt.p, got, tt.want)
		}
	}
}

func TestTimePeriod_Next(t *testing.T) {
	// All → Month → Week → Day → All (cycle)
	p := PeriodAll
	p = p.Next()
	if p != PeriodMonth {
		t.Errorf("expected PeriodMonth, got %s", p)
	}
	p = p.Next()
	if p != PeriodWeek {
		t.Errorf("expected PeriodWeek, got %s", p)
	}
	p = p.Next()
	if p != PeriodDay {
		t.Errorf("expected PeriodDay, got %s", p)
	}
	p = p.Next()
	if p != PeriodAll {
		t.Errorf("expected PeriodAll, got %s", p)
	}
}

func TestTimePeriod_Cutoff(t *testing.T) {
	// PeriodAll returns zero time.
	if c := PeriodAll.Cutoff(); !c.IsZero() {
		t.Errorf("PeriodAll.Cutoff() should be zero, got %v", c)
	}

	// Other periods return a time in the past.
	now := time.Now()
	for _, p := range []TimePeriod{PeriodMonth, PeriodWeek, PeriodDay} {
		c := p.Cutoff()
		if c.IsZero() || c.After(now) {
			t.Errorf("%s.Cutoff() should be in the past, got %v", p, c)
		}
	}
}

func TestExtractTimestamp(t *testing.T) {
	ts := "2025-01-15T10:30:00.123Z"
	line := []byte(fmt.Sprintf(`{"type":"assistant","timestamp":"%s","message":{}}`, ts))

	got, ok := extractTimestamp(line)
	if !ok {
		t.Fatal("extractTimestamp returned false")
	}
	want, _ := time.Parse(time.RFC3339Nano, ts)
	if !got.Equal(want) {
		t.Errorf("got %v, want %v", got, want)
	}

	// No timestamp field.
	_, ok = extractTimestamp([]byte(`{"type":"user","message":"hello"}`))
	if ok {
		t.Error("expected false for line without timestamp")
	}
}

func TestCollectWithCutoff_FiltersOldLines(t *testing.T) {
	home := t.TempDir()
	dir := filepath.Join(home, ".claude", "transcripts")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}

	now := time.Now()
	recent := now.Add(-1 * time.Hour).Format(time.RFC3339Nano)
	old := now.Add(-48 * time.Hour).Format(time.RFC3339Nano)

	lines := []string{
		fmt.Sprintf(`{"type":"assistant","timestamp":"%s","message":{"content":[{"type":"tool_use","name":"Task","input":{"subagent_type":"code-reviewer","prompt":"review"}}]}}`, recent),
		fmt.Sprintf(`{"type":"assistant","timestamp":"%s","message":{"content":[{"type":"tool_use","name":"Task","input":{"subagent_type":"librarian","prompt":"research"}}]}}`, old),
	}
	writeJSONL(t, filepath.Join(dir, "session.jsonl"), lines)

	// No cutoff: both counted.
	counts, err := collectAgents(home, "", time.Time{})
	if err != nil {
		t.Fatal(err)
	}
	if counts["code-reviewer"] != 1 || counts["librarian"] != 1 {
		t.Errorf("no cutoff: code-reviewer=%d, librarian=%d", counts["code-reviewer"], counts["librarian"])
	}

	// Cutoff at 24h ago: only recent line counted.
	cutoff := now.Add(-24 * time.Hour)
	counts, err = collectAgents(home, "", cutoff)
	if err != nil {
		t.Fatal(err)
	}
	if counts["code-reviewer"] != 1 {
		t.Errorf("with cutoff: code-reviewer=%d, want 1", counts["code-reviewer"])
	}
	if counts["librarian"] != 0 {
		t.Errorf("with cutoff: librarian=%d, want 0 (filtered out)", counts["librarian"])
	}
}

func TestCollectTools_WithCutoff(t *testing.T) {
	home := t.TempDir()
	dir := filepath.Join(home, ".claude", "transcripts")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}

	now := time.Now()
	recent := now.Add(-2 * time.Hour).Format(time.RFC3339Nano)
	old := now.Add(-72 * time.Hour).Format(time.RFC3339Nano)

	lines := []string{
		fmt.Sprintf(`{"type":"assistant","timestamp":"%s","message":{"content":[{"type":"tool_use","name":"Read","input":{"file_path":"a.go"}}]}}`, recent),
		fmt.Sprintf(`{"type":"assistant","timestamp":"%s","message":{"content":[{"type":"tool_use","name":"Bash","input":{"command":"ls"}}]}}`, old),
	}
	writeJSONL(t, filepath.Join(dir, "session.jsonl"), lines)

	// With 24h cutoff: only Read counted.
	cutoff := now.Add(-24 * time.Hour)
	counts, err := collectTools(home, "", cutoff)
	if err != nil {
		t.Fatal(err)
	}
	if counts["Read"] != 1 {
		t.Errorf("Read: got %d, want 1", counts["Read"])
	}
	if counts["Bash"] != 0 {
		t.Errorf("Bash: got %d, want 0 (filtered out)", counts["Bash"])
	}
}
