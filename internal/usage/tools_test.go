package usage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestCollectTools(t *testing.T) {
	home := t.TempDir()
	dir := filepath.Join(home, ".claude", "usage-data", "session-meta")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}

	// 세션 1: 프로젝트 A
	writeMeta(t, dir, "s1.json", sessionMeta{
		ProjectPath: "/project/a",
		ToolCounts:  map[string]int{"Read": 10, "Bash": 5},
	})
	// 세션 2: 프로젝트 B
	writeMeta(t, dir, "s2.json", sessionMeta{
		ProjectPath: "/project/b",
		ToolCounts:  map[string]int{"Read": 3, "Write": 7},
	})

	// 전체 범위
	counts, err := collectTools(home, "")
	if err != nil {
		t.Fatal(err)
	}
	if counts["Read"] != 13 {
		t.Errorf("Read: got %d, want 13", counts["Read"])
	}
	if counts["Bash"] != 5 {
		t.Errorf("Bash: got %d, want 5", counts["Bash"])
	}

	// 프로젝트 필터
	counts, err = collectTools(home, "/project/a")
	if err != nil {
		t.Fatal(err)
	}
	if counts["Read"] != 10 {
		t.Errorf("filtered Read: got %d, want 10", counts["Read"])
	}
	if counts["Write"] != 0 {
		t.Errorf("filtered Write: got %d, want 0", counts["Write"])
	}
}

func TestCollectTools_MissingDir(t *testing.T) {
	counts, err := collectTools(t.TempDir(), "")
	if err != nil {
		t.Errorf("missing dir should not error: %v", err)
	}
	if counts != nil {
		t.Errorf("expected nil counts, got %v", counts)
	}
}

func writeMeta(t *testing.T, dir, name string, meta sessionMeta) {
	t.Helper()
	data, err := json.Marshal(meta)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, name), data, 0o644); err != nil {
		t.Fatal(err)
	}
}
