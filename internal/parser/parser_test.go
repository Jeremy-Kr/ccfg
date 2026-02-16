package parser

import (
	"strings"
	"testing"
)

func TestStripJSONC(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string // 정리된 JSON의 핵심 내용
	}{
		{
			name:  "한 줄 주석 제거",
			input: `{"key": "value" // 주석}`,
			want:  `"key"`,
		},
		{
			name:  "블록 주석 제거",
			input: `{"key": /* 주석 */ "value"}`,
			want:  `"value"`,
		},
		{
			name: "trailing comma 제거",
			input: `{
				"a": 1,
				"b": 2,
			}`,
			want: `"b": 2`,
		},
		{
			name:  "문자열 내 슬래시 보존",
			input: `{"url": "https://example.com"}`,
			want:  `https://example.com`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StripJSONC(tt.input)
			if !strings.Contains(got, tt.want) {
				t.Errorf("StripJSONC 결과에 %q가 없음\ngot: %s", tt.want, got)
			}
		})
	}
}

func TestFormatJSON(t *testing.T) {
	input := `{"name":"ccfg","version":"0.1.0"}`
	result := FormatJSON(input)

	// pretty-print 되었으면 개행이 있어야 함
	if !strings.Contains(result, "\n") {
		t.Error("FormatJSON이 pretty-print하지 않음")
	}
	// 핵심 값이 포함되어야 함
	if !strings.Contains(result, "ccfg") {
		t.Error("FormatJSON 결과에 'ccfg'가 없음")
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
		t.Error("JSONC 파싱 후 'permissions'가 없음")
	}
}

func TestFormatMarkdown(t *testing.T) {
	input := "# Hello\n\nThis is **bold** text."
	result := FormatMarkdown(input)

	// 렌더링 되었으면 원본과 달라야 함
	if result == input {
		t.Error("FormatMarkdown이 렌더링하지 않음")
	}
	// 핵심 텍스트가 포함되어야 함
	if !strings.Contains(result, "Hello") {
		t.Error("FormatMarkdown 결과에 'Hello'가 없음")
	}
}
