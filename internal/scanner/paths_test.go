package scanner

import (
	"runtime"
	"testing"
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
	if len(entries) != 5 {
		t.Errorf("UserPaths entries 개수: got %d, want 5", len(entries))
	}

	// settings.json이 첫 번째 항목인지 확인
	if entries[0].RelPath != ".claude/settings.json" {
		t.Errorf("첫 번째 항목: got %q, want .claude/settings.json", entries[0].RelPath)
	}
}

func TestProjectPaths(t *testing.T) {
	// 빈 root는 빈 결과
	base, entries := ProjectPaths("")
	if base != "" || entries != nil {
		t.Error("빈 root에서 nil이 아닌 결과 반환")
	}

	// 유효한 root
	base, entries = ProjectPaths("/tmp/myproject")
	if base != "/tmp/myproject" {
		t.Errorf("base: got %q, want /tmp/myproject", base)
	}
	if len(entries) != 5 {
		t.Errorf("ProjectPaths entries 개수: got %d, want 5", len(entries))
	}
}

func TestCategoryAssignment(t *testing.T) {
	_, userEntries := UserPaths()

	categoryMap := make(map[string]ConfigCategory)
	for _, e := range userEntries {
		categoryMap[e.RelPath] = e.Category
	}

	// settings 계열은 CategorySettings
	if categoryMap[".claude/settings.json"] != CategorySettings {
		t.Error("settings.json의 Category가 CategorySettings가 아님")
	}

	// CLAUDE.md는 CategoryInstructions
	if categoryMap[".claude/CLAUDE.md"] != CategoryInstructions {
		t.Error("CLAUDE.md의 Category가 CategoryInstructions가 아님")
	}

	// .mcp.json은 CategoryMCP
	if categoryMap[".mcp.json"] != CategoryMCP {
		t.Error(".mcp.json의 Category가 CategoryMCP가 아님")
	}
}
