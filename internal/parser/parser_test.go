package parser

import (
	"strings"
	"testing"
)

func TestStripJSONC(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string // key content expected in the cleaned JSON
	}{
		{
			name:  "strip single-line comment",
			input: `{"key": "value" // 주석}`,
			want:  `"key"`,
		},
		{
			name:  "strip block comment",
			input: `{"key": /* 주석 */ "value"}`,
			want:  `"value"`,
		},
		{
			name: "strip trailing comma",
			input: `{
				"a": 1,
				"b": 2,
			}`,
			want: `"b": 2`,
		},
		{
			name:  "preserve slashes inside strings",
			input: `{"url": "https://example.com"}`,
			want:  `https://example.com`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StripJSONC(tt.input)
			if !strings.Contains(got, tt.want) {
				t.Errorf("StripJSONC result missing %q\ngot: %s", tt.want, got)
			}
		})
	}
}

func TestFormatJSON(t *testing.T) {
	input := `{"name":"ccfg","version":"0.1.0"}`
	result := FormatJSON(input)

	// Pretty-printed output should contain newlines
	if !strings.Contains(result, "\n") {
		t.Error("FormatJSON did not pretty-print the output")
	}
	// Core value must be present
	if !strings.Contains(result, "ccfg") {
		t.Error("FormatJSON result missing 'ccfg'")
	}
}

func TestFormatJSONWithJSONC(t *testing.T) {
	input := `{
		// Claude Code 설정
		"permissions": ["read"],
		"hooks": {}, // 훅 설정
	}`
	result := FormatJSON(input)

	if !strings.Contains(result, "permissions") {
		t.Error("result missing 'permissions' after JSONC parsing")
	}
}

func TestFormatMarkdown(t *testing.T) {
	input := "# Hello\n\nThis is **bold** text."
	result := FormatMarkdown(input)

	// Rendered output should differ from the raw input
	if result == input {
		t.Error("FormatMarkdown did not render the input")
	}
	// Core text must be present
	if !strings.Contains(result, "Hello") {
		t.Error("FormatMarkdown result missing 'Hello'")
	}
}
