package parser

import (
	"sort"
	"testing"
)

func TestParseSettingsHooks(t *testing.T) {
	raw := `{
		// JSONC 주석 테스트
		"hooks": {
			"SessionStart": [
				{"command": "echo hello"},
				{"command": "echo world"}
			],
			"Stop": [
				{"command": "cleanup.sh"}
			]
		}
	}`

	entries := ParseSettingsHooks(raw)
	if len(entries) != 2 {
		t.Fatalf("expected 2 hook events, got %d", len(entries))
	}

	// Sort for stable test assertions
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Event < entries[j].Event
	})

	if entries[0].Event != "SessionStart" {
		t.Errorf("expected SessionStart, got %s", entries[0].Event)
	}
	if entries[0].Count != 2 {
		t.Errorf("expected count 2, got %d", entries[0].Count)
	}
	if len(entries[0].Commands) != 2 {
		t.Errorf("expected 2 commands, got %d", len(entries[0].Commands))
	}

	if entries[1].Event != "Stop" {
		t.Errorf("expected Stop, got %s", entries[1].Event)
	}
	if entries[1].Count != 1 {
		t.Errorf("expected count 1, got %d", entries[1].Count)
	}
}

func TestParseSettingsHooks_NoHooks(t *testing.T) {
	raw := `{"permissions": {}}`
	entries := ParseSettingsHooks(raw)
	if entries != nil {
		t.Errorf("expected nil for no hooks, got %v", entries)
	}
}

func TestParseSettingsHooks_InvalidJSON(t *testing.T) {
	entries := ParseSettingsHooks("not json")
	if entries != nil {
		t.Errorf("expected nil for invalid json, got %v", entries)
	}
}

func TestParseMCPServers(t *testing.T) {
	raw := `{
		"mcpServers": {
			"serena": {
				"type": "stdio",
				"command": "uvx serena"
			},
			"pencil": {
				"command": "npx pencil-mcp"
			}
		}
	}`

	entries := ParseMCPServers(raw)
	if len(entries) != 2 {
		t.Fatalf("expected 2 servers, got %d", len(entries))
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name < entries[j].Name
	})

	if entries[0].Name != "pencil" {
		t.Errorf("expected pencil, got %s", entries[0].Name)
	}
	if entries[0].Command != "npx pencil-mcp" {
		t.Errorf("expected 'npx pencil-mcp', got %s", entries[0].Command)
	}

	if entries[1].Name != "serena" {
		t.Errorf("expected serena, got %s", entries[1].Name)
	}
	if entries[1].Type != "stdio" {
		t.Errorf("expected stdio, got %s", entries[1].Type)
	}
}

func TestParseMCPServers_NoServers(t *testing.T) {
	raw := `{"hooks": {}}`
	entries := ParseMCPServers(raw)
	if entries != nil {
		t.Errorf("expected nil for no mcpServers, got %v", entries)
	}
}
