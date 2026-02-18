package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jeremy-kr/ccfg/internal/usage"
)

// rankingHeaderRows is the number of rows consumed by the ranking header
// (tab bar + scope bar + period bar + separator).
const rankingHeaderRows = 4

// RankingModel manages the state of the ranking view.
type RankingModel struct {
	data      *usage.UsageData
	tab       usage.RankCategory
	scope     usage.DataScope
	period    usage.TimePeriod
	cursor    int
	offset    int
	height    int
	collector *usage.Collector
	err       error
}

// NewRankingModel creates a RankingModel with the given Collector.
func NewRankingModel(collector *usage.Collector) RankingModel {
	return RankingModel{
		tab:       usage.RankAgents,
		scope:     usage.ScopeAll,
		collector: collector,
	}
}

// Load collects usage data.
func (r *RankingModel) Load() {
	r.collector.Period = r.period
	data, err := r.collector.Collect(r.scope)
	r.data = data
	r.err = err
	r.cursor = 0
	r.offset = 0
}

// SetHeight sets the number of visible rows.
func (r *RankingModel) SetHeight(h int) {
	r.height = h
}

// entries returns the list of entries for the current tab.
func (r *RankingModel) entries() []usage.RankEntry {
	if r.data == nil {
		return nil
	}
	switch r.tab {
	case usage.RankAgents:
		return r.data.Agents
	case usage.RankTools:
		return r.data.Tools
	case usage.RankSkills:
		return r.data.Skills
	default:
		return nil
	}
}

// NextTab moves to the next tab.
func (r *RankingModel) NextTab() {
	r.tab = (r.tab + 1) % 3
	r.cursor = 0
	r.offset = 0
}

// SetTab sets the tab directly.
func (r *RankingModel) SetTab(tab usage.RankCategory) {
	r.tab = tab
	r.cursor = 0
	r.offset = 0
}

// ToggleScope toggles the scope and reloads data.
func (r *RankingModel) ToggleScope() {
	if r.scope == usage.ScopeAll {
		r.scope = usage.ScopeProject
	} else {
		r.scope = usage.ScopeAll
	}
	r.Load()
}

// TogglePeriod cycles through time periods: All → 30d → 7d → 24h → All.
func (r *RankingModel) TogglePeriod() {
	r.period = r.period.Next()
	r.Load()
}

// MoveUp moves the cursor up.
func (r *RankingModel) MoveUp() {
	if r.cursor > 0 {
		r.cursor--
		r.adjustScroll()
	}
}

// MoveDown moves the cursor down.
func (r *RankingModel) MoveDown() {
	entries := r.entries()
	if r.cursor < len(entries)-1 {
		r.cursor++
		r.adjustScroll()
	}
}

func (r *RankingModel) adjustScroll() {
	if r.height <= 0 {
		return
	}
	if r.cursor < r.offset {
		r.offset = r.cursor
	}
	if r.cursor >= r.offset+r.height {
		r.offset = r.cursor - r.height + 1
	}
}

// View renders the ranking view.
func (r *RankingModel) View(width, height int) string {
	r.height = height - rankingHeaderRows

	var b strings.Builder

	// Tab bar.
	b.WriteString(r.renderTabs(width))
	b.WriteString("\n")

	// Scope bar.
	b.WriteString(r.renderScopeBar(width))
	b.WriteString("\n")

	// Period bar.
	b.WriteString(r.renderPeriodBar(width))
	b.WriteString("\n")

	// Separator.
	sep := lipgloss.NewStyle().Foreground(colorDimGray).Render(strings.Repeat("─", width-4))
	b.WriteString(sep)
	b.WriteString("\n")

	// Error display.
	if r.err != nil {
		errMsg := lipgloss.NewStyle().Foreground(colorRed).Render(fmt.Sprintf("Error: %v", r.err))
		b.WriteString(errMsg)
		return b.String()
	}

	// No data.
	entries := r.entries()
	if len(entries) == 0 {
		empty := lipgloss.NewStyle().Foreground(colorDimGray).Render("  No data")
		b.WriteString(empty)
		return b.String()
	}

	// Ranking list.
	visibleRows := r.height
	if visibleRows < 1 {
		visibleRows = 1
	}
	end := r.offset + visibleRows
	if end > len(entries) {
		end = len(entries)
	}

	scrollBars := renderScrollbar(len(entries), visibleRows, r.offset)

	contentW := width
	barWidth := width - 35 // rank(4) + grade(6) + name(15) + count(6) + padding(4).
	if scrollBars != nil {
		contentW = width - 1
		barWidth--
	}
	if barWidth < 5 {
		barWidth = 5
	}

	for i := r.offset; i < end; i++ {
		entry := entries[i]
		selected := i == r.cursor

		line := r.renderEntry(i+1, entry, barWidth, selected)
		if scrollBars != nil {
			if gap := contentW - lipgloss.Width(line); gap > 0 {
				line += strings.Repeat(" ", gap)
			}
			line += scrollBars[i-r.offset]
		}
		b.WriteString(line)
		if i < end-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

func (r *RankingModel) renderTabs(width int) string {
	tabs := []struct {
		cat   usage.RankCategory
		emoji string
		label string
	}{
		{usage.RankAgents, "🤖", "Agents"},
		{usage.RankTools, "🔧", "Tools"},
		{usage.RankSkills, "🧠", "Skills"},
	}

	var parts []string
	for _, t := range tabs {
		label := fmt.Sprintf(" %s %s ", t.emoji, t.label)
		if t.cat == r.tab {
			parts = append(parts, lipgloss.NewStyle().
				Bold(true).
				Foreground(colorYellow).
				Background(lipgloss.Color("#333333")).
				Render(label))
		} else {
			parts = append(parts, lipgloss.NewStyle().
				Foreground(colorDimGray).
				Render(label))
		}
	}

	tabBar := strings.Join(parts, lipgloss.NewStyle().Foreground(colorDimGray).Render(" │ "))

	// Place key hints on the right.
	hint := hudDesc.Render("1/2/3: tab  Tab: next")
	pad := width - lipgloss.Width(tabBar) - lipgloss.Width(hint) - 4
	if pad < 1 {
		pad = 1
	}
	return tabBar + strings.Repeat(" ", pad) + hint
}

func (r *RankingModel) renderScopeBar(width int) string {
	activeStyle := lipgloss.NewStyle().Bold(true).Background(lipgloss.Color("#333333"))
	inactiveStyle := lipgloss.NewStyle().Foreground(colorDimGray)

	allStyle, projStyle := inactiveStyle, activeStyle.Foreground(colorCyan)
	if r.scope == usage.ScopeAll {
		allStyle, projStyle = activeStyle.Foreground(colorYellow), inactiveStyle
	}

	scopeBar := hudDesc.Render("Scope: ") + allStyle.Render(" All ") + hudDesc.Render(" / ") + projStyle.Render(" Project ")

	hint := hudDesc.Render("s: scope  p: period")
	pad := width - lipgloss.Width(scopeBar) - lipgloss.Width(hint) - 4
	if pad < 1 {
		pad = 1
	}
	return scopeBar + strings.Repeat(" ", pad) + hint
}

func (r *RankingModel) renderPeriodBar(width int) string {
	activeStyle := lipgloss.NewStyle().Bold(true).Background(lipgloss.Color("#333333"))
	inactiveStyle := lipgloss.NewStyle().Foreground(colorDimGray)

	periods := []usage.TimePeriod{usage.PeriodAll, usage.PeriodMonth, usage.PeriodWeek, usage.PeriodDay}
	var parts []string
	for _, p := range periods {
		label := fmt.Sprintf(" %s ", p)
		if p == r.period {
			parts = append(parts, activeStyle.Foreground(colorYellow).Render(label))
		} else {
			parts = append(parts, inactiveStyle.Render(label))
		}
	}

	periodBar := hudDesc.Render("Period: ") + strings.Join(parts, hudDesc.Render(" / "))
	return periodBar
}

func (r *RankingModel) renderEntry(rank int, entry usage.RankEntry, barWidth int, selected bool) string {
	gs := gradeStyle(entry.Grade)

	rankStr := fmt.Sprintf("%2d.", rank)
	badge := fmt.Sprintf("[%-3s]", entry.Grade)
	name := entry.Name
	if len(name) > 15 {
		name = name[:14] + "…"
	}
	namePad := 15 - lipgloss.Width(name)
	if namePad < 0 {
		namePad = 0
	}

	filled := int(entry.LogScore * float64(barWidth))
	if filled < 1 && entry.Count > 0 {
		filled = 1
	}
	empty := barWidth - filled
	if empty < 0 {
		empty = 0
	}

	if selected {
		// Selected item: highlighted with background color.
		sel := lipgloss.NewStyle().Bold(true).Foreground(colorYellow).Background(lipgloss.Color("#333333"))
		bar := sel.Render(strings.Repeat("█", filled)) +
			lipgloss.NewStyle().Foreground(lipgloss.Color("#555555")).Background(lipgloss.Color("#333333")).Render(strings.Repeat("░", empty))
		return sel.Render(fmt.Sprintf(" %s %s %s%s ", rankStr, badge, name, strings.Repeat(" ", namePad))) +
			bar + sel.Render(fmt.Sprintf(" %d", entry.Count))
	}

	bar := gs.Render(strings.Repeat("█", filled)) +
		lipgloss.NewStyle().Foreground(lipgloss.Color("#333333")).Render(strings.Repeat("░", empty))
	return fmt.Sprintf(" %s %s %s%s %s %d",
		rankStr, gs.Render(badge), name, strings.Repeat(" ", namePad), bar, entry.Count)
}
