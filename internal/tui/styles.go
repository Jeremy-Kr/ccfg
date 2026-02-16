package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/jeremy-kr/ccfg/internal/usage"
)

var (
	// Retro arcade color palette.
	colorYellow  = lipgloss.Color("#FFD700") // Gold yellow — title, focus.
	colorOrange  = lipgloss.Color("#FF8C00") // Orange — border, accent.
	colorRed     = lipgloss.Color("#FF4444") // Red — missing file, warning.
	colorGreen   = lipgloss.Color("#39FF14") // Neon green — existing file, success.
	colorCyan    = lipgloss.Color("#00FFFF") // Cyan — info, Project scope.
	colorMagenta = lipgloss.Color("#FF00FF") // Magenta — search, special accent.
	colorDimGray = lipgloss.Color("#555555") // Dim gray — inactive border.

	// Header style.
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorYellow)

	// Footer (HUD) style.
	footerStyle = lipgloss.NewStyle().
			Padding(0, 1)

	// Panel border (unfocused) — double line.
	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(colorDimGray).
			Padding(0, 1)

	// Panel border (focused) — double line.
	panelFocusedStyle = lipgloss.NewStyle().
				Border(lipgloss.DoubleBorder()).
				BorderForeground(colorOrange).
				Padding(0, 1)

	// Tree item style.
	treeItemStyle = lipgloss.NewStyle()

	// Tree selected item style.
	treeSelectedStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorYellow)

	// Scope header style (default — overridden per scope in renderNode).
	scopeHeaderStyle = lipgloss.NewStyle().
				Bold(true)

	// File exists style.
	fileExistsStyle = lipgloss.NewStyle().
			Foreground(colorGreen)

	// File missing style.
	fileMissingStyle = lipgloss.NewStyle().
				Foreground(colorRed)

	// Directory node style.
	dirStyle = lipgloss.NewStyle().
			Foreground(colorOrange)

	// HUD element styles.
	hudLabelNav = lipgloss.NewStyle().Bold(true).Foreground(colorGreen)
	hudLabelCmd = lipgloss.NewStyle().Bold(true).Foreground(colorCyan)
	hudKey      = lipgloss.NewStyle().Bold(true).Foreground(colorYellow)
	hudDesc     = lipgloss.NewStyle().Foreground(colorDimGray)
	hudSep      = lipgloss.NewStyle().Foreground(colorOrange)

	// Agent character card styles.
	agentCardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorOrange).
			Padding(0, 1)
	agentCardTitleStyle = lipgloss.NewStyle().Bold(true).Foreground(colorYellow)
	agentCardRoleStyle  = lipgloss.NewStyle().Foreground(colorOrange)

	// Skill ability card styles.
	skillCardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorCyan).
			Padding(0, 1)
	skillCardTitleStyle = lipgloss.NewStyle().Bold(true).Foreground(colorCyan)
	skillCardTagStyle   = lipgloss.NewStyle().Foreground(colorMagenta)
)

// gradeColors maps each grade to its display color.
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

// gradeStyle returns the style for the given grade.
func gradeStyle(g usage.Grade) lipgloss.Style {
	color, ok := gradeColors[g]
	if !ok {
		color = colorDimGray
	}
	return lipgloss.NewStyle().Bold(true).Foreground(color)
}

// panelStyleFor returns the panel style based on focus state.
func panelStyleFor(focused bool) lipgloss.Style {
	if focused {
		return panelFocusedStyle
	}
	return panelStyle
}
