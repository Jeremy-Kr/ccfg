package tui

import "github.com/charmbracelet/bubbles/key"

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
	Quit     key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("↑/k", "위로"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("↓/j", "아래로"),
	),
	Toggle: key.NewBinding(
		key.WithKeys("enter", " "),
		key.WithHelp("enter", "펼치기/접기"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "패널 전환"),
	),
	Left: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "왼쪽 패널"),
	),
	Right: key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l", "오른쪽 패널"),
	),
	PageUp: key.NewBinding(
		key.WithKeys("pgup", "ctrl+u"),
		key.WithHelp("pgup", "페이지 위"),
	),
	PageDown: key.NewBinding(
		key.WithKeys("pgdown", "ctrl+d"),
		key.WithHelp("pgdn", "페이지 아래"),
	),
	Search: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "검색"),
	),
	Merge: key.NewBinding(
		key.WithKeys("m"),
		key.WithHelp("m", "병합 뷰"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "종료"),
	),
}

func (k keyMap) helpLine() string {
	return "jk: 이동  enter: 펼치기  tab/hl: 패널  /: 검색  m: 병합  q: 종료"
}
