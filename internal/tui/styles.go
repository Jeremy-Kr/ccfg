package tui

import "github.com/charmbracelet/lipgloss"

var (
	// 색상
	colorPrimary   = lipgloss.Color("#7C3AED") // 보라색
	colorSecondary = lipgloss.Color("#6B7280") // 회색
	colorAccent    = lipgloss.Color("#10B981") // 초록 (존재하는 파일)
	colorMuted     = lipgloss.Color("#4B5563") // 어두운 회색 (미존재 파일)
	colorBorder    = lipgloss.Color("#374151")
	colorFocused   = lipgloss.Color("#7C3AED")

	// 헤더
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPrimary).
			Padding(0, 1)

	// 풋터 (키 힌트)
	footerStyle = lipgloss.NewStyle().
			Foreground(colorSecondary).
			Padding(0, 1)

	// 패널 테두리 (포커스 없음)
	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder).
			Padding(0, 1)

	// 패널 테두리 (포커스)
	panelFocusedStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colorFocused).
				Padding(0, 1)

	// 트리 항목
	treeItemStyle = lipgloss.NewStyle()

	// 트리 선택된 항목
	treeSelectedStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorPrimary)

	// Scope 헤더 (트리 내)
	scopeHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#F9FAFB"))

	// 파일 존재
	fileExistsStyle = lipgloss.NewStyle().
			Foreground(colorAccent)

	// 파일 미존재
	fileMissingStyle = lipgloss.NewStyle().
				Foreground(colorMuted)
)
