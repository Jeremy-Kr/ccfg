package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

type keyMap struct {
	Up       key.Binding
	Down     key.Binding
	Toggle   key.Binding
	Tab      key.Binding
	Left     key.Binding
	Right    key.Binding
	PageUp   key.Binding
	PageDown key.Binding
	Search   key.Binding
	Merge    key.Binding
	Ranking  key.Binding
	Quit     key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("↓/j", "down"),
	),
	Toggle: key.NewBinding(
		key.WithKeys("enter", " "),
		key.WithHelp("enter", "toggle"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "switch panel"),
	),
	Left: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "left panel"),
	),
	Right: key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l", "right panel"),
	),
	PageUp: key.NewBinding(
		key.WithKeys("pgup", "ctrl+u"),
		key.WithHelp("pgup", "page up"),
	),
	PageDown: key.NewBinding(
		key.WithKeys("pgdown", "ctrl+d"),
		key.WithHelp("pgdn", "page down"),
	),
	Search: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "search"),
	),
	Merge: key.NewBinding(
		key.WithKeys("m"),
		key.WithHelp("m", "merge view"),
	),
	Ranking: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "ranking"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

// renderHUD renders the HUD footer.
func renderHUD(existCount, totalCount int, scopeName string, scanSec float64, watching bool) string {
	sep := hudSep.Render(" │ ")

	nav := hudLabelNav.Render("[NAV]") + " " +
		hudKey.Render("↑↓") + hudDesc.Render(" move  ") +
		hudKey.Render("⏎") + hudDesc.Render(" toggle  ") +
		hudKey.Render("⇥") + hudDesc.Render(" panel")

	cmd := hudLabelCmd.Render("[CMD]") + " " +
		hudKey.Render("/") + hudDesc.Render(" search  ") +
		hudKey.Render("m") + hudDesc.Render(" merge  ") +
		hudKey.Render("r") + hudDesc.Render(" ranking  ") +
		hudKey.Render("q") + hudDesc.Render(" quit")

	stats := fmt.Sprintf("📊 %s/%s",
		fileExistsStyle.Render(fmt.Sprintf("%d", existCount)),
		fileMissingStyle.Render(fmt.Sprintf("%d", totalCount)),
	)

	scope := hudDesc.Render(scopeName)
	scan := hudDesc.Render(fmt.Sprintf("⏱ %.1fs", scanSec))

	hud := nav + sep + cmd + sep + stats + sep + scope + sep + scan
	if watching {
		hud += sep + lipgloss.NewStyle().Foreground(colorGreen).Render("👁 watching")
	}
	return hud
}
