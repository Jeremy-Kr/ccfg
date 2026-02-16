package tui

import "github.com/charmbracelet/lipgloss"

var (
	scrollThumbStyle = lipgloss.NewStyle().Foreground(colorCyan)
	scrollTrackStyle = lipgloss.NewStyle().Foreground(colorDimGray)
)

// renderScrollbar returns a vertical scrollbar as a slice of strings.
// Each element is one row of the scrollbar (1 column wide).
// total is the total item count, visible is the visible item count, and offset is the current offset.
// Returns nil when scrolling is not needed.
func renderScrollbar(total, visible, offset int) []string {
	if total <= visible || visible <= 0 {
		return nil
	}

	// Thumb size: minimum 1 row.
	thumbSize := visible * visible / total
	if thumbSize < 1 {
		thumbSize = 1
	}

	// Thumb position.
	maxOffset := total - visible
	track := visible - thumbSize
	thumbPos := 0
	if maxOffset > 0 && track > 0 {
		thumbPos = offset * track / maxOffset
	}

	bars := make([]string, visible)
	for i := range bars {
		if i >= thumbPos && i < thumbPos+thumbSize {
			bars[i] = scrollThumbStyle.Render("┃")
		} else {
			bars[i] = scrollTrackStyle.Render("│")
		}
	}
	return bars
}
