package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jeremy-kr/ccfg/internal/merger"
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
	scan       *model.ScanResult
	tree       TreeModel
	preview    PreviewModel
	focus      Pane
	width      int
	height     int
	ready      bool
	searchMode bool
	searchText string
	mergeMode  bool
	merged     *merger.MergedConfig
}

// NewModel은 ScanResult로부터 TUI 모델을 생성한다.
func NewModel(result *model.ScanResult) Model {
	tree := NewTreeModel(result)
	m := Model{
		scan:   result,
		tree:   tree,
		focus:  PaneTree,
		merged: merger.Merge(result),
	}
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
		// 검색 모드
		if m.searchMode {
			return m.updateSearch(msg)
		}

		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, keys.Search):
			m.searchMode = true
			m.searchText = ""
			return m, nil

		case key.Matches(msg, keys.Merge):
			m.mergeMode = !m.mergeMode
			return m, nil

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

func (m Model) updateSearch(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEscape:
		m.searchMode = false
		m.searchText = ""
		m.tree.ClearFilter()
		return m, nil
	case tea.KeyEnter:
		m.searchMode = false
		// 필터 유지
		return m, nil
	case tea.KeyBackspace:
		if len(m.searchText) > 0 {
			m.searchText = m.searchText[:len(m.searchText)-1]
		}
		m.tree.Filter(m.searchText)
		return m, nil
	default:
		if msg.Type == tea.KeyRunes {
			m.searchText += string(msg.Runes)
			m.tree.Filter(m.searchText)
		}
		return m, nil
	}
}

func (m Model) View() string {
	if !m.ready {
		return "로딩 중..."
	}

	// 헤더
	title := fmt.Sprintf("ccfg v%s — Claude Code Config Viewer", version)
	if m.mergeMode {
		title += "  [MERGED]"
	}
	header := headerStyle.Render(title)

	// 풋터
	var footer string
	if m.searchMode {
		footer = footerStyle.Render(fmt.Sprintf("🔍 /%s█  (Enter: 확인, Esc: 취소)", m.searchText))
	} else {
		footer = footerStyle.Render(keys.helpLine())
	}

	// 메인 영역 치수
	contentH := m.contentHeight()
	treeW := m.treeWidth()
	previewW := m.previewWidth()

	// 패널 렌더링
	m.tree.SetHeight(contentH)
	m.preview.SetHeight(contentH)
	treeView := m.tree.View(treeW, m.focus == PaneTree)

	var previewView string
	if m.mergeMode {
		previewView = m.renderMergeView(previewW, contentH)
	} else {
		previewView = m.preview.View(previewW, m.focus == PanePreview)
	}

	main := lipgloss.JoinHorizontal(lipgloss.Top, treeView, previewView)

	return lipgloss.JoinVertical(lipgloss.Left, header, main, footer)
}

func (m *Model) renderMergeView(width, height int) string {
	content := m.merged.Render()
	lines := strings.Split(content, "\n")

	var b strings.Builder
	end := height
	if end > len(lines) {
		end = len(lines)
	}
	for i := 0; i < end; i++ {
		b.WriteString(lines[i])
		if i < end-1 {
			b.WriteString("\n")
		}
	}

	style := panelStyle.Width(width)
	if m.focus == PanePreview {
		style = panelFocusedStyle.Width(width)
	}
	return style.Render(b.String())
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

func (m *Model) contentHeight() int {
	h := m.height - 4
	if h < 3 {
		h = 3
	}
	return h
}

func (m *Model) treeWidth() int {
	w := m.width * 30 / 100
	if w < 20 {
		w = 20
	}
	return w
}

func (m *Model) previewWidth() int {
	return m.width - m.treeWidth()
}
