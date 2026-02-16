package usage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCollectAgents(t *testing.T) {
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
		`{"type":"tool_use","tool_name":"skill","tool_input":{"name":"git-master"}}`,
	}
	writeJSONL(t, filepath.Join(dir, "test.jsonl"), lines)

	counts, err := collectAgents(home, "")
	if err != nil {
		t.Fatal(err)
	}
	if counts["explore"] != 2 {
		t.Errorf("explore: got %d, want 2", counts["explore"])
	}
	if counts["librarian"] != 1 {
		t.Errorf("librarian: got %d, want 1", counts["librarian"])
	}
	// Read와 skill은 에이전트가 아님
	if counts["Read"] != 0 {
		t.Errorf("Read should not be counted as agent")
	}
}

func TestCollectAgents_ProjectScope(t *testing.T) {
	home := t.TempDir()
	// 프로젝트별 transcript 디렉토리
	projDir := filepath.Join(home, ".claude", "projects", "-project-foo")
	if err := os.MkdirAll(projDir, 0o755); err != nil {
		t.Fatal(err)
	}

	lines := []string{
		`{"type":"tool_use","tool_name":"task","tool_input":{"subagent_type":"code-reviewer","prompt":"review"}}`,
	}
	writeJSONL(t, filepath.Join(projDir, "session.jsonl"), lines)

	counts, err := collectAgents(home, "/project/foo")
	if err != nil {
		t.Fatal(err)
	}
	if counts["code-reviewer"] != 1 {
		t.Errorf("code-reviewer: got %d, want 1", counts["code-reviewer"])
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

func TestCollectSkills(t *testing.T) {
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

	counts, err := collectSkills(home, "")
	if err != nil {
		t.Fatal(err)
	}
	if counts["git-master"] != 2 {
		t.Errorf("git-master: got %d, want 2", counts["git-master"])
	}
	if counts["commit"] != 1 {
		t.Errorf("commit: got %d, want 1", counts["commit"])
	}
	// task는 스킬이 아님
	if counts["explore"] != 0 {
		t.Errorf("explore should not be counted as skill")
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
