package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jeremy-kr/ccfg/internal/model"
)

// Pane은 현재 포커스된 패널을 나타낸다.
type Pane int

const (
	PaneTree Pane = iota
	PanePreview
)

const version = "0.1.0"

// Model은 TUI 전체 상태를 관리한다.
type Model struct {
	scan    *model.ScanResult
	tree    TreeModel
	preview PreviewModel
	focus   Pane
	width   int
	height  int
	ready   bool
}

// NewModel은 ScanResult로부터 TUI 모델을 생성한다.
func NewModel(result *model.ScanResult) Model {
	tree := NewTreeModel(result)
	m := Model{
		scan:  result,
		tree:  tree,
		focus: PaneTree,
	}
	// 초기 파일 선택 반영
	if f := tree.SelectedFile(); f != nil {
		m.preview.SetFile(f)
	}
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		m.updateLayout()
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, keys.Tab):
			m.toggleFocus()
			return m, nil

		case key.Matches(msg, keys.Left):
			m.focus = PaneTree
			return m, nil

		case key.Matches(msg, keys.Right):
			m.focus = PanePreview
			return m, nil

		case key.Matches(msg, keys.Up):
			if m.focus == PaneTree {
				m.tree.MoveUp()
				m.syncPreview()
			} else {
				m.preview.ScrollUp(1)
			}
			return m, nil

		case key.Matches(msg, keys.Down):
			if m.focus == PaneTree {
				m.tree.MoveDown()
				m.syncPreview()
			} else {
				m.preview.ScrollDown(1)
			}
			return m, nil

		case key.Matches(msg, keys.Toggle):
			if m.focus == PaneTree {
				m.tree.Toggle()
				m.syncPreview()
			}
			return m, nil

		case key.Matches(msg, keys.PageUp):
			if m.focus == PanePreview {
				m.preview.ScrollUp(m.contentHeight() / 2)
			}
			return m, nil

		case key.Matches(msg, keys.PageDown):
			if m.focus == PanePreview {
				m.preview.ScrollDown(m.contentHeight() / 2)
			}
			return m, nil
		}
	}

	return m, nil
}

func (m Model) View() string {
	if !m.ready {
		return "로딩 중..."
	}

	// 헤더
	header := headerStyle.Render(fmt.Sprintf("ccfg v%s — Claude Code Config Viewer", version))

	// 풋터
	footer := footerStyle.Render(keys.helpLine())

	// 메인 영역 치수
	contentH := m.contentHeight()
	treeW := m.treeWidth()
	previewW := m.previewWidth()

	// 패널 렌더링
	m.tree.SetHeight(contentH)
	m.preview.SetHeight(contentH)
	treeView := m.tree.View(treeW, m.focus == PaneTree)
	previewView := m.preview.View(previewW, m.focus == PanePreview)

	// 좌우 배치
	main := lipgloss.JoinHorizontal(lipgloss.Top, treeView, previewView)

	return lipgloss.JoinVertical(lipgloss.Left, header, main, footer)
}

func (m *Model) toggleFocus() {
	if m.focus == PaneTree {
		m.focus = PanePreview
	} else {
		m.focus = PaneTree
	}
}

func (m *Model) syncPreview() {
	m.preview.SetFile(m.tree.SelectedFile())
}

func (m *Model) updateLayout() {
	h := m.contentHeight()
	m.tree.SetHeight(h)
	m.preview.SetHeight(h)
}

// contentHeight는 헤더/풋터를 제외한 메인 영역 높이.
func (m *Model) contentHeight() int {
	// 헤더 1줄 + 풋터 1줄 + 패널 테두리 위아래 2줄
	h := m.height - 4
	if h < 3 {
		h = 3
	}
	return h
}

// treeWidth는 좌측 패널 너비 (30%).
func (m *Model) treeWidth() int {
	w := m.width * 30 / 100
	if w < 20 {
		w = 20
	}
	return w
}

// previewWidth는 우측 패널 너비 (나머지).
func (m *Model) previewWidth() int {
	return m.width - m.treeWidth()
}
