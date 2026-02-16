package parser

import (
	"github.com/charmbracelet/glamour"
)

// renderer는 Glamour 렌더러의 싱글턴.
var renderer *glamour.TermRenderer

func init() {
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(0), // 래핑은 TUI 패널이 처리
	)
	if err != nil {
		// 렌더러 생성 실패 시 nil — FormatMarkdown에서 폴백
		return
	}
	renderer = r
}

// FormatMarkdown은 Markdown을 터미널에 렌더링한다.
// 실패하면 원본 텍스트를 반환한다.
func FormatMarkdown(raw string) string {
	if renderer == nil {
		return raw
	}
	out, err := renderer.Render(raw)
	if err != nil {
		return raw
	}
	return out
}
