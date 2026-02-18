package usage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCollectAgents_Opencode(t *testing.T) {
	home := t.TempDir()
	dir := filepath.Join(home, ".claude", "transcripts")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}

	lines := []string{
		`{"type":"tool_use","tool_name":"task","tool_input":{"subagent_type":"explore","prompt":"find files"}}`,
		`{"type":"tool_use","tool_name":"task","tool_input":{"subagent_type":"explore","prompt":"search"}}`,
		`{"type":"tool_use","tool_name":"delegate_task","tool_input":{"subagent_type":"librarian","prompt":"research"}}`,
		`{"type":"tool_use","tool_name":"Read","tool_input":{"path":"foo.go"}}`,
	}
	writeJSONL(t, filepath.Join(dir, "test.jsonl"), lines)

	counts, err := collectAgents(home, "", time.Time{})
	if err != nil {
		t.Fatal(err)
	}
	// explore is a built-in infrastructure agent and should be excluded
	if counts["explore"] != 0 {
		t.Errorf("explore: got %d, want 0 (excluded)", counts["explore"])
	}
	if counts["librarian"] != 1 {
		t.Errorf("librarian: got %d, want 1", counts["librarian"])
	}
	if counts["Read"] != 0 {
		t.Errorf("Read should not be counted as agent")
	}
}

func TestCollectAgents_ClaudeCode(t *testing.T) {
	home := t.TempDir()
	dir := filepath.Join(home, ".claude", "projects", "-project-bar")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Claude Code format: tool_use blocks inside an assistant message
	lines := []string{
		`{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Task","input":{"subagent_type":"explore","prompt":"search"}}]}}`,
		`{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Task","input":{"subagent_type":"code-reviewer","prompt":"review"}},{"type":"text","text":"reviewing..."}]}}`,
		`{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Read","input":{"file_path":"foo.go"}}]}}`,
		`{"type":"user","message":"hello"}`,
	}
	writeJSONL(t, filepath.Join(dir, "session.jsonl"), lines)

	// Verify that projects/ subdirectories are scanned in All scope
	counts, err := collectAgents(home, "", time.Time{})
	if err != nil {
		t.Fatal(err)
	}
	// explore is a built-in infrastructure agent and should be excluded
	if counts["explore"] != 0 {
		t.Errorf("explore: got %d, want 0 (excluded)", counts["explore"])
	}
	if counts["code-reviewer"] != 1 {
		t.Errorf("code-reviewer: got %d, want 1", counts["code-reviewer"])
	}
}

func TestCollectAgents_GeneralPurposeResolve(t *testing.T) {
	home := t.TempDir()
	dir := filepath.Join(home, ".claude", "projects", "-project-x")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}

	lines := []string{
		// Custom agent: "name:" pattern in description
		`{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Task","input":{"subagent_type":"general-purpose","description":"horner: evaluate code","prompt":"..."}}]}}`,
		`{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Task","input":{"subagent_type":"general-purpose","description":"horner: review PR","prompt":"..."}}]}}`,
		`{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Task","input":{"subagent_type":"general-purpose","description":"newey: design spec","prompt":"..."}}]}}`,
		// No name pattern: stays as general-purpose
		`{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Task","input":{"subagent_type":"general-purpose","description":"Generate rules file","prompt":"..."}}]}}`,
		// Prefix with spaces: stays as general-purpose
		`{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Task","input":{"subagent_type":"general-purpose","description":"Agent #1: do something","prompt":"..."}}]}}`,
	}
	writeJSONL(t, filepath.Join(dir, "session.jsonl"), lines)

	counts, err := collectAgents(home, "", time.Time{})
	if err != nil {
		t.Fatal(err)
	}
	if counts["horner"] != 2 {
		t.Errorf("horner: got %d, want 2", counts["horner"])
	}
	if counts["newey"] != 1 {
		t.Errorf("newey: got %d, want 1", counts["newey"])
	}
	if counts["general-purpose"] != 2 {
		t.Errorf("general-purpose: got %d, want 2", counts["general-purpose"])
	}
}

func TestCollectAgents_ProjectScope(t *testing.T) {
	home := t.TempDir()
	projDir := filepath.Join(home, ".claude", "projects", "-project-foo")
	if err := os.MkdirAll(projDir, 0o755); err != nil {
		t.Fatal(err)
	}

	lines := []string{
		`{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Task","input":{"subagent_type":"code-reviewer","prompt":"review"}}]}}`,
	}
	writeJSONL(t, filepath.Join(projDir, "session.jsonl"), lines)

	counts, err := collectAgents(home, "/project/foo", time.Time{})
	if err != nil {
		t.Fatal(err)
	}
	if counts["code-reviewer"] != 1 {
		t.Errorf("code-reviewer: got %d, want 1", counts["code-reviewer"])
	}
}

func TestCollectSkills_Opencode(t *testing.T) {
	home := t.TempDir()
	dir := filepath.Join(home, ".claude", "transcripts")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}

	lines := []string{
		`{"type":"tool_use","tool_name":"skill","tool_input":{"name":"git-master"}}`,
		`{"type":"tool_use","tool_name":"skill","tool_input":{"name":"git-master"}}`,
		`{"type":"tool_use","tool_name":"skill","tool_input":{"name":"commit"}}`,
		`{"type":"tool_use","tool_name":"task","tool_input":{"subagent_type":"explore"}}`,
	}
	writeJSONL(t, filepath.Join(dir, "test.jsonl"), lines)

	counts, err := collectSkills(home, "", time.Time{})
	if err != nil {
		t.Fatal(err)
	}
	if counts["git-master"] != 2 {
		t.Errorf("git-master: got %d, want 2", counts["git-master"])
	}
	if counts["commit"] != 1 {
		t.Errorf("commit: got %d, want 1", counts["commit"])
	}
	if counts["explore"] != 0 {
		t.Errorf("explore should not be counted as skill")
	}
}

func TestCollectSkills_ClaudeCode(t *testing.T) {
	home := t.TempDir()
	dir := filepath.Join(home, ".claude", "projects", "-project-bar")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}

	lines := []string{
		`{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Skill","input":{"skill":"commit"}}]}}`,
		`{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Skill","input":{"skill":"commit"}}]}}`,
		`{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Task","input":{"subagent_type":"explore"}}]}}`,
	}
	writeJSONL(t, filepath.Join(dir, "session.jsonl"), lines)

	counts, err := collectSkills(home, "", time.Time{})
	if err != nil {
		t.Fatal(err)
	}
	if counts["commit"] != 2 {
		t.Errorf("commit: got %d, want 2", counts["commit"])
	}
	if counts["explore"] != 0 {
		t.Errorf("explore should not be counted as skill")
	}
}

func TestCollectTools_ClaudeCode(t *testing.T) {
	home := t.TempDir()
	dir := filepath.Join(home, ".claude", "projects", "-project-baz")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Multiple tool_use blocks in a single assistant line
	lines := []string{
		`{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Read","input":{"file_path":"a.go"}},{"type":"tool_use","name":"Read","input":{"file_path":"b.go"}}]}}`,
		`{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Bash","input":{"command":"go test"}}]}}`,
		`{"type":"assistant","message":{"content":[{"type":"text","text":"done"}]}}`,
	}
	writeJSONL(t, filepath.Join(dir, "session.jsonl"), lines)

	counts, err := collectTools(home, "", time.Time{})
	if err != nil {
		t.Fatal(err)
	}
	if counts["Read"] != 2 {
		t.Errorf("Read: got %d, want 2", counts["Read"])
	}
	if counts["Bash"] != 1 {
		t.Errorf("Bash: got %d, want 1", counts["Bash"])
	}
}

func TestEncodeProjectPath(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"/Users/jeremy/code/ccfg", "-Users-jeremy-code-ccfg"},
		{"/project/foo", "-project-foo"},
	}
	for _, tt := range tests {
		got := encodeProjectPath(tt.input)
		if got != tt.want {
			t.Errorf("encodeProjectPath(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestTranscriptDirs_AllScope(t *testing.T) {
	home := t.TempDir()

	// Create both transcripts/ and projects/
	os.MkdirAll(filepath.Join(home, ".claude", "transcripts"), 0o755)
	os.MkdirAll(filepath.Join(home, ".claude", "projects", "-proj-a"), 0o755)
	os.MkdirAll(filepath.Join(home, ".claude", "projects", "-proj-b"), 0o755)

	dirs := transcriptDirs(home, "")
	if len(dirs) != 3 {
		t.Errorf("expected 3 dirs (transcripts + 2 projects), got %d: %v", len(dirs), dirs)
	}
}

func TestCollectTools_SessionMeta(t *testing.T) {
	home := t.TempDir()
	dir := filepath.Join(home, ".claude", "usage-data", "session-meta")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}

	writeMeta(t, dir, "s1.json", sessionMeta{
		ProjectPath: "/project/a",
		ToolCounts:  map[string]int{"Read": 10, "Bash": 5},
	})

	counts, err := collectTools(home, "/project/a", time.Time{})
	if err != nil {
		t.Fatal(err)
	}
	if counts["Read"] != 10 {
		t.Errorf("Read: got %d, want 10", counts["Read"])
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

func writeJSONL(t *testing.T, path string, lines []string) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	for _, line := range lines {
		f.WriteString(line + "\n")
	}
}
