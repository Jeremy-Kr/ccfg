package tui

import "github.com/charmbracelet/lipgloss"

var (
	scrollThumbStyle = lipgloss.NewStyle().Foreground(colorCyan)
	scrollTrackStyle = lipgloss.NewStyle().Foreground(colorDimGray)
)

// renderScrollbar는 세로 스크롤바 문자열 슬라이스를 반환한다.
// 각 요소는 한 행의 스크롤바 문자(1열).
// total: 전체 항목 수, visible: 보이는 항목 수, offset: 현재 오프셋.
// 스크롤 불필요 시 빈 슬라이스를 반환한다.
func renderScrollbar(total, visible, offset int) []string {
	if total <= visible || visible <= 0 {
		return nil
	}

	// thumb 크기: 최소 1행
	thumbSize := visible * visible / total
	if thumbSize < 1 {
		thumbSize = 1
	}

	// thumb 위치
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
