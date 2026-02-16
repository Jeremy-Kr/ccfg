package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/jeremy-kr/ccfg/internal/usage"
)

var (
	// 레트로 아케이드 색상 팔레트
	colorYellow  = lipgloss.Color("#FFD700") // 골드 옐로우 — 제목, 포커스
	colorOrange  = lipgloss.Color("#FF8C00") // 오렌지 — 테두리, 강조
	colorRed     = lipgloss.Color("#FF4444") // 레드 — 미존재 파일, 경고
	colorGreen   = lipgloss.Color("#39FF14") // 네온 그린 — 존재 파일, 성공
	colorCyan    = lipgloss.Color("#00FFFF") // 시안 — 정보, Project scope
	colorMagenta = lipgloss.Color("#FF00FF") // 마젠타 — 검색, 특수 강조
	colorDimGray = lipgloss.Color("#555555") // 어두운 회색 — 비활성 테두리

	// 헤더
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorYellow)

	// 풋터 (HUD)
	footerStyle = lipgloss.NewStyle().
			Padding(0, 1)

	// 패널 테두리 (포커스 없음) — Double Line
	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(colorDimGray).
			Padding(0, 1)

	// 패널 테두리 (포커스) — Double Line
	panelFocusedStyle = lipgloss.NewStyle().
				Border(lipgloss.DoubleBorder()).
				BorderForeground(colorOrange).
				Padding(0, 1)

	// 트리 항목
	treeItemStyle = lipgloss.NewStyle()

	// 트리 선택된 항목
	treeSelectedStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorYellow)

	// Scope 헤더 스타일 (기본 — renderNode에서 Scope별로 오버라이드)
	scopeHeaderStyle = lipgloss.NewStyle().
				Bold(true)

	// 파일 존재
	fileExistsStyle = lipgloss.NewStyle().
			Foreground(colorGreen)

	// 파일 미존재
	fileMissingStyle = lipgloss.NewStyle().
				Foreground(colorRed)

	// 디렉토리 노드
	dirStyle = lipgloss.NewStyle().
			Foreground(colorOrange)

	// HUD 요소별 스타일
	hudLabelNav = lipgloss.NewStyle().Bold(true).Foreground(colorGreen)
	hudLabelCmd = lipgloss.NewStyle().Bold(true).Foreground(colorCyan)
	hudKey      = lipgloss.NewStyle().Bold(true).Foreground(colorYellow)
	hudDesc     = lipgloss.NewStyle().Foreground(colorDimGray)
	hudSep      = lipgloss.NewStyle().Foreground(colorOrange)
)

// 등급별 색상
var gradeColors = map[usage.Grade]lipgloss.Color{
	usage.GradeSSS: colorMagenta,
	usage.GradeSS:  colorYellow,
	usage.GradeS:   colorOrange,
	usage.GradeA:   colorRed,
	usage.GradeB:   colorCyan,
	usage.GradeC:   colorGreen,
	usage.GradeD:   lipgloss.Color("#888888"),
	usage.GradeF:   lipgloss.Color("#444444"),
}

// gradeStyle은 등급에 해당하는 스타일을 반환한다.
func gradeStyle(g usage.Grade) lipgloss.Style {
	color, ok := gradeColors[g]
	if !ok {
		color = colorDimGray
	}
	return lipgloss.NewStyle().Bold(true).Foreground(color)
}

// panelStyleFor는 포커스 상태에 따라 패널 스타일을 반환한다.
func panelStyleFor(focused bool) lipgloss.Style {
	if focused {
		return panelFocusedStyle
	}
	return panelStyle
}
