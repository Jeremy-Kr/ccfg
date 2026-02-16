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
			t.Fatal("base is empty on a supported OS")
		}
		if len(entries) == 0 {
			t.Fatal("entries is empty on a supported OS")
		}
	}

	for _, e := range entries {
		if e.RelPath == "" {
			t.Error("found entry with empty RelPath")
		}
		if e.Description == "" {
			t.Error("found entry with empty Description")
		}
	}
}

func TestUserPaths(t *testing.T) {
	base, entries := UserPaths()

	if base == "" {
		t.Fatal("UserPaths base is empty")
	}
	if len(entries) != 9 {
		t.Errorf("UserPaths entries count: got %d, want 9", len(entries))
	}

	if entries[0].RelPath != ".claude/settings.json" {
		t.Errorf("first entry: got %q, want .claude/settings.json", entries[0].RelPath)
	}
}

func TestProjectPaths(t *testing.T) {
	base, entries := ProjectPaths("")
	if base != "" || entries != nil {
		t.Error("non-nil result returned for empty root")
	}

	base, entries = ProjectPaths("/tmp/myproject")
	if base != "/tmp/myproject" {
		t.Errorf("base: got %q, want /tmp/myproject", base)
	}
	if len(entries) != 8 {
		t.Errorf("ProjectPaths entries count: got %d, want 8", len(entries))
	}
}

func TestCategoryAssignment(t *testing.T) {
	_, userEntries := UserPaths()

	categoryMap := make(map[string]model.ConfigCategory)
	for _, e := range userEntries {
		categoryMap[e.RelPath] = e.Category
	}

	if categoryMap[".claude/settings.json"] != model.CategorySettings {
		t.Error("settings.json Category is not CategorySettings")
	}
	if categoryMap[".claude/CLAUDE.md"] != model.CategoryInstructions {
		t.Error("CLAUDE.md Category is not CategoryInstructions")
	}
	if categoryMap[".mcp.json"] != model.CategoryMCP {
		t.Error(".mcp.json Category is not CategoryMCP")
	}
	if categoryMap[".claude/commands"] != model.CategoryCommands {
		t.Error("commands Category is not CategoryCommands")
	}
	if categoryMap[".claude/skills"] != model.CategorySkills {
		t.Error("skills Category is not CategorySkills")
	}
	if categoryMap[".claude/agents"] != model.CategoryAgents {
		t.Error("agents Category is not CategoryAgents")
	}
	if categoryMap[".claude/keybindings.json"] != model.CategoryKeybindings {
		t.Error("keybindings.json Category is not CategoryKeybindings")
	}
}
