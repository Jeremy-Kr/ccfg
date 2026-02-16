package parser

import (
	"github.com/charmbracelet/glamour"
)

// renderer is the singleton Glamour terminal renderer.
var renderer *glamour.TermRenderer

func init() {
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(0), // Wrapping is handled by the TUI panel
	)
	if err != nil {
		// If renderer creation fails, leave nil; FormatMarkdown falls back to raw text
		return
	}
	renderer = r
}

// FormatMarkdown renders Markdown for the terminal.
// Returns the original text on failure.
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
