package tui

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jeremy-kr/ccfg/internal/merger"
	"github.com/jeremy-kr/ccfg/internal/model"
	"github.com/jeremy-kr/ccfg/internal/scanner"
	"github.com/jeremy-kr/ccfg/internal/usage"
	"github.com/jeremy-kr/ccfg/internal/watcher"
)

// Pane represents the currently focused panel.
type Pane int

const (
	PaneTree Pane = iota
	PanePreview
)

const version = "0.1.0"

// Model manages the overall TUI state.
type Model struct {
	scan         *model.ScanResult
	tree         TreeModel
	preview      PreviewModel
	focus        Pane
	width        int
	height       int
	ready        bool
	searchMode   bool
	searchText   string
	mergeMode    bool
	merged       *merger.MergedConfig
	rankingMode  bool
	ranking      RankingModel
	scanDuration time.Duration
	watcher      *watcher.Watcher  // File watcher (nil if inactive).
	sc           *scanner.Scanner  // For rescanning.
}

// NewModel creates a TUI model from a ScanResult.
func NewModel(result *model.ScanResult, scanDuration time.Duration, s *scanner.Scanner) Model {
	tree := NewTreeModel(result)
	homeDir, _ := os.UserHomeDir()
	m := Model{
		scan:         result,
		tree:         tree,
		focus:        PaneTree,
		merged:       merger.Merge(result),
		ranking:      NewRankingModel(&usage.Collector{HomeDir: homeDir, ProjectPath: result.RootDir}),
		scanDuration: scanDuration,
		sc:           s,
	}
	if f := tree.SelectedFile(); f != nil {
		m.preview.SetFile(f)
	}

	// Create file watcher (nil on failure — operates without watching).
	paths := scanner.WatchPaths(result.RootDir)
	if w, err := watcher.New(paths); err == nil {
		m.watcher = w
	}

	return m
}

// fileStats returns the count of existing files and the total file count.
// Virtual nodes (IsVirtual) are excluded since they are not real files.
func (m *Model) fileStats() (exist, total int) {
	for _, f := range m.scan.All() {
		total++
		if f.Exists {
			exist++
		}
		for _, c := range f.Children {
			if c.IsVirtual {
				continue
			}
			total++
			if c.Exists {
				exist++
			}
		}
	}
	return
}

func (m Model) Init() tea.Cmd {
	return m.waitCmd()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case watcher.FileChangedMsg:
		return m.handleFileChanged()

	case watcher.ErrorMsg:
		return m, m.waitCmd()

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		m.updateLayout()
		return m, nil

	case tea.KeyMsg:
		// Search mode.
		if m.searchMode {
			return m.updateSearch(msg)
		}

		// Ranking mode.
		if m.rankingMode {
			return m.updateRanking(msg)
		}

		switch {
		case key.Matches(msg, keys.Quit):
			if m.watcher != nil {
				m.watcher.Close()
			}
			return m, tea.Quit

		case key.Matches(msg, keys.Search):
			m.searchMode = true
			m.searchText = ""
			return m, nil

		case key.Matches(msg, keys.Merge):
			m.mergeMode = !m.mergeMode
			return m, nil

		case key.Matches(msg, keys.Ranking):
			m.rankingMode = true
			m.mergeMode = false
			m.ranking.Load()
			m.ranking.SetHeight(m.contentHeight() - rankingHeaderRows)
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
		// Keep filter active.
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
		return "Loading..."
	}

	// Ranking mode — fullscreen.
	if m.rankingMode {
		return m.renderRankingView()
	}

	// Header — decorated line.
	header := m.renderHeader()

	// Footer — HUD or search bar.
	var footer string
	if m.searchMode {
		searchBar := lipgloss.NewStyle().Foreground(colorMagenta).Render(
			fmt.Sprintf("🔍 /%s█  (Enter: confirm, Esc: cancel)", m.searchText),
		)
		footer = footerStyle.Render(searchBar)
	} else {
		existCount, totalCount := m.fileStats()
		scopeName := m.tree.SelectedScope().String()
		scanSec := m.scanDuration.Seconds()
		footer = footerStyle.Render(renderHUD(existCount, totalCount, scopeName, scanSec, m.watcher != nil))
	}

	// Main area dimensions.
	contentH := m.contentHeight()
	treeW := m.treeWidth()
	previewW := m.previewWidth()

	// Render panels.
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

func (m *Model) renderHeader() string {
	subtitle := "Claude Code Config Viewer ⚡"
	if m.rankingMode {
		subtitle = lipgloss.NewStyle().Bold(true).Foreground(colorYellow).Render("🏆 RANKING VIEW 🏆")
	} else if m.mergeMode {
		subtitle = lipgloss.NewStyle().Bold(true).Foreground(colorMagenta).Render("⚡ MERGE VIEW ⚡")
	}
	title := fmt.Sprintf("⚡ CCFG v%s — %s", version, subtitle)

	label := fmt.Sprintf("[ %s ]", title)
	pad := m.width - lipgloss.Width(label)
	if pad < 2 {
		pad = 2
	}
	left := pad / 2
	right := pad - left
	line := strings.Repeat("═", left) + label + strings.Repeat("═", right)
	return headerStyle.Render(line)
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

	base := panelStyleFor(m.focus == PanePreview)
	style := base.Width(width - base.GetHorizontalBorderSize()).Height(height)
	availWidth := width - style.GetHorizontalFrameSize()
	truncated := lipgloss.NewStyle().MaxWidth(availWidth).Render(b.String())

	return style.Render(truncated)
}

func (m Model) updateRanking(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Quit):
		return m, tea.Quit
	case msg.Type == tea.KeyEscape, key.Matches(msg, keys.Ranking):
		m.rankingMode = false
		return m, nil
	case key.Matches(msg, keys.Up):
		m.ranking.MoveUp()
		return m, nil
	case key.Matches(msg, keys.Down):
		m.ranking.MoveDown()
		return m, nil
	case key.Matches(msg, keys.Tab):
		m.ranking.NextTab()
		return m, nil
	case msg.Type == tea.KeyRunes:
		switch string(msg.Runes) {
		case "1":
			m.ranking.SetTab(usage.RankAgents)
		case "2":
			m.ranking.SetTab(usage.RankTools)
		case "3":
			m.ranking.SetTab(usage.RankSkills)
		case "s":
			m.ranking.ToggleScope()
		case "p":
			m.ranking.TogglePeriod()
		}
		return m, nil
	}
	return m, nil
}

func (m *Model) renderRankingView() string {
	header := m.renderHeader()
	contentH := m.contentHeight()
	panelFrameW := panelFocusedStyle.GetHorizontalFrameSize()
	rankingContent := m.ranking.View(m.width-2-panelFrameW, contentH)

	// Ranking HUD.
	footer := footerStyle.Render(renderRankingHUD())

	style := panelFocusedStyle.Width(m.width - 2).Height(contentH)
	body := style.Render(rankingContent)

	return lipgloss.JoinVertical(lipgloss.Left, header, body, footer)
}

func renderRankingHUD() string {
	sep := hudSep.Render(" │ ")

	nav := hudLabelNav.Render("[NAV]") + " " +
		hudKey.Render("↑↓") + hudDesc.Render(" move  ") +
		hudKey.Render("1/2/3") + hudDesc.Render(" tab  ") +
		hudKey.Render("⇥") + hudDesc.Render(" next tab")

	cmd := hudLabelCmd.Render("[CMD]") + " " +
		hudKey.Render("s") + hudDesc.Render(" scope  ") +
		hudKey.Render("p") + hudDesc.Render(" period  ") +
		hudKey.Render("r/Esc") + hudDesc.Render(" close  ") +
		hudKey.Render("q") + hudDesc.Render(" quit")

	return nav + sep + cmd
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
	m.preview.PrepareCardContent(m.previewWidth())
}

func (m *Model) updateLayout() {
	h := m.contentHeight()
	m.tree.SetHeight(h)
	m.preview.SetHeight(h)
	m.preview.PrepareCardContent(m.previewWidth())
	m.ranking.SetHeight(h - rankingHeaderRows)
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

// waitCmd returns a command to wait for the next file change if the watcher is active.
func (m *Model) waitCmd() tea.Cmd {
	if m.watcher != nil {
		return m.watcher.WaitForChange()
	}
	return nil
}

// handleFileChanged performs a rescan and rebuilds the tree when a file change is detected.
func (m Model) handleFileChanged() (tea.Model, tea.Cmd) {
	if m.sc == nil {
		return m, m.waitCmd()
	}

	// Capture tree state.
	state := m.tree.CaptureState()

	// Rescan.
	start := time.Now()
	result, err := m.sc.Scan()
	scanDuration := time.Since(start)
	if err != nil {
		return m, m.waitCmd()
	}

	// Rebuild tree and restore state.
	m.scan = result
	m.scanDuration = scanDuration
	m.tree = NewTreeModel(result)
	m.tree.RestoreState(state)
	m.tree.SetHeight(m.contentHeight())

	// Update merge.
	m.merged = merger.Merge(result)

	// Update preview.
	m.preview.InvalidateCache()
	m.syncPreview()

	return m, m.waitCmd()
}

func (m *Model) previewWidth() int {
	return m.width - m.treeWidth()
}
