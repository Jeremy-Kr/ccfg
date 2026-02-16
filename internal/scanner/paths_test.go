package scanner

import (
	"runtime"
	"testing"

	"github.com/jeremy-kr/ccfg/internal/model"
)

func TestManagedPaths(t *testing.T) {
	base, entries := ManagedPaths()

	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		if base == "" {
			t.Fatal("지원되는 OS에서 base가 빈 문자열")
		}
		if len(entries) == 0 {
			t.Fatal("지원되는 OS에서 entries가 비어있음")
		}
	}

	for _, e := range entries {
		if e.RelPath == "" {
			t.Error("RelPath가 빈 문자열인 항목이 있음")
		}
		if e.Description == "" {
			t.Error("Description이 빈 문자열인 항목이 있음")
		}
	}
}

func TestUserPaths(t *testing.T) {
	base, entries := UserPaths()

	if base == "" {
		t.Fatal("UserPaths base가 빈 문자열")
	}
	if len(entries) != 9 {
		t.Errorf("UserPaths entries 개수: got %d, want 9", len(entries))
	}

	if entries[0].RelPath != ".claude/settings.json" {
		t.Errorf("첫 번째 항목: got %q, want .claude/settings.json", entries[0].RelPath)
	}
}

func TestProjectPaths(t *testing.T) {
	base, entries := ProjectPaths("")
	if base != "" || entries != nil {
		t.Error("빈 root에서 nil이 아닌 결과 반환")
	}

	base, entries = ProjectPaths("/tmp/myproject")
	if base != "/tmp/myproject" {
		t.Errorf("base: got %q, want /tmp/myproject", base)
	}
	if len(entries) != 7 {
		t.Errorf("ProjectPaths entries 개수: got %d, want 7", len(entries))
	}
}

func TestCategoryAssignment(t *testing.T) {
	_, userEntries := UserPaths()

	categoryMap := make(map[string]model.ConfigCategory)
	for _, e := range userEntries {
		categoryMap[e.RelPath] = e.Category
	}

	if categoryMap[".claude/settings.json"] != model.CategorySettings {
		t.Error("settings.json의 Category가 CategorySettings가 아님")
	}
	if categoryMap[".claude/CLAUDE.md"] != model.CategoryInstructions {
		t.Error("CLAUDE.md의 Category가 CategoryInstructions가 아님")
	}
	if categoryMap[".mcp.json"] != model.CategoryMCP {
		t.Error(".mcp.json의 Category가 CategoryMCP가 아님")
	}
	if categoryMap[".claude/commands"] != model.CategoryCommands {
		t.Error("commands의 Category가 CategoryCommands가 아님")
	}
	if categoryMap[".claude/skills"] != model.CategorySkills {
		t.Error("skills의 Category가 CategorySkills가 아님")
	}
	if categoryMap[".claude/agents"] != model.CategoryAgents {
		t.Error("agents의 Category가 CategoryAgents가 아님")
	}
	if categoryMap[".claude/keybindings.json"] != model.CategoryKeybindings {
		t.Error("keybindings.json의 Category가 CategoryKeybindings가 아님")
	}
}
