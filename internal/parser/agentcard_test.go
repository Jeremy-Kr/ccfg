package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.md")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestParseAgentMeta_WithFrontmatter(t *testing.T) {
	content := `---
name: brawn
description: "체계적으로 요구사항을 분석합니다"
model: opus
color: yellow
---

You are a PM/PO agent.
`
	meta := ParseAgentMeta(writeTempFile(t, content))
	if meta == nil {
		t.Fatal("expected non-nil AgentMeta")
	}
	if meta.Name != "brawn" {
		t.Errorf("Name = %q, want %q", meta.Name, "brawn")
	}
	if meta.Desc != "체계적으로 요구사항을 분석합니다" {
		t.Errorf("Desc = %q", meta.Desc)
	}
	if meta.Model != "opus" {
		t.Errorf("Model = %q, want %q", meta.Model, "opus")
	}
	if meta.Color != "yellow" {
		t.Errorf("Color = %q, want %q", meta.Color, "yellow")
	}
}

func TestParseAgentMeta_WithoutFrontmatter(t *testing.T) {
	content := `# brawn — PM/PO

체계적으로 요구사항을 분석하고 정밀한 스펙을 작성합니다.
`
	meta := ParseAgentMeta(writeTempFile(t, content))
	if meta == nil {
		t.Fatal("expected non-nil AgentMeta")
	}
	if meta.Name != "brawn" {
		t.Errorf("Name = %q, want %q", meta.Name, "brawn")
	}
	if meta.Role != "PM/PO" {
		t.Errorf("Role = %q, want %q", meta.Role, "PM/PO")
	}
	if !strings.Contains(meta.Desc, "체계적으로") {
		t.Errorf("Desc = %q, expected to contain 체계적으로", meta.Desc)
	}
}

func TestParseAgentMeta_EmptyFile(t *testing.T) {
	meta := ParseAgentMeta(writeTempFile(t, ""))
	if meta != nil {
		t.Error("expected nil for empty file")
	}
}

func TestParseAgentMeta_NonExistentFile(t *testing.T) {
	meta := ParseAgentMeta("/nonexistent/path.md")
	if meta != nil {
		t.Error("expected nil for nonexistent file")
	}
}

func TestParseAgentMeta_LongDescription(t *testing.T) {
	long := strings.Repeat("가", 200) // 200자 한글
	content := `---
name: test-agent
description: "` + long + `"
---
`
	meta := ParseAgentMeta(writeTempFile(t, content))
	if meta == nil {
		t.Fatal("expected non-nil AgentMeta")
	}
	runes := []rune(meta.Desc)
	if len(runes) > maxDescLen+1 { // +1 for "…"
		t.Errorf("Desc length = %d runes, expected <= %d", len(runes), maxDescLen+1)
	}
	if !strings.HasSuffix(meta.Desc, "…") {
		t.Errorf("Desc should end with …, got %q", meta.Desc[len(meta.Desc)-3:])
	}
}

func TestParseSkillMeta_WithFrontmatter(t *testing.T) {
	content := `---
name: peon-ping-config
description: Update peon-ping configuration settings.
category: toolchain
---

# peon-ping-config

Update peon-ping configuration settings.
`
	meta := ParseSkillMeta(writeTempFile(t, content))
	if meta == nil {
		t.Fatal("expected non-nil SkillMeta")
	}
	if meta.Name != "peon-ping-config" {
		t.Errorf("Name = %q, want %q", meta.Name, "peon-ping-config")
	}
	if !strings.Contains(meta.Desc, "peon-ping") {
		t.Errorf("Desc = %q", meta.Desc)
	}
	if meta.Category != "toolchain" {
		t.Errorf("Category = %q, want %q", meta.Category, "toolchain")
	}
}

func TestParseSkillMeta_WithoutFrontmatter(t *testing.T) {
	content := `# biome

Fast all-in-one linter and formatter.

## Tags

- linting
- formatting
- code-quality
`
	meta := ParseSkillMeta(writeTempFile(t, content))
	if meta == nil {
		t.Fatal("expected non-nil SkillMeta")
	}
	if meta.Name != "biome" {
		t.Errorf("Name = %q, want %q", meta.Name, "biome")
	}
	if !strings.Contains(meta.Desc, "linter") {
		t.Errorf("Desc = %q", meta.Desc)
	}
	if meta.Tags != "linting, formatting, code-quality" {
		t.Errorf("Tags = %q", meta.Tags)
	}
}

func TestParseSkillMeta_EmptyFile(t *testing.T) {
	meta := ParseSkillMeta(writeTempFile(t, ""))
	if meta != nil {
		t.Error("expected nil for empty file")
	}
}

func TestParseFrontmatter_NoFrontmatter(t *testing.T) {
	fm, body := parseFrontmatter("Hello world")
	if len(fm) != 0 {
		t.Errorf("expected empty frontmatter, got %v", fm)
	}
	if body != "Hello world" {
		t.Errorf("body = %q", body)
	}
}

func TestParseFrontmatter_QuotedValues(t *testing.T) {
	content := `---
name: "test"
desc: 'hello'
---
body`
	fm, _ := parseFrontmatter(content)
	if fm["name"] != "test" {
		t.Errorf("name = %q, want %q", fm["name"], "test")
	}
	if fm["desc"] != "hello" {
		t.Errorf("desc = %q, want %q", fm["desc"], "hello")
	}
}

func TestCleanDesc_EscapedNewlines(t *testing.T) {
	desc := `First sentence.\nSecond sentence.\nThird.`
	got := cleanDesc(desc)
	if strings.Contains(got, "\\n") {
		t.Errorf("desc still contains \\n: %q", got)
	}
}
