package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jeremy-kr/ccfg/internal/usage"
)

// RankingModel은 랭킹 뷰의 상태를 관리한다.
type RankingModel struct {
	data      *usage.UsageData
	tab       usage.RankCategory
	scope     usage.DataScope
	cursor    int
	offset    int
	height    int
	collector *usage.Collector
	err       error
}

// NewRankingModel은 Collector로 RankingModel을 생성한다.
func NewRankingModel(collector *usage.Collector) RankingModel {
	return RankingModel{
		tab:       usage.RankAgents,
		scope:     usage.ScopeAll,
		collector: collector,
	}
}

// Load는 사용 데이터를 수집한다.
func (r *RankingModel) Load() {
	data, err := r.collector.Collect(r.scope)
	r.data = data
	r.err = err
	r.cursor = 0
	r.offset = 0
}

// SetHeight는 표시 가능한 행 수를 설정한다.
func (r *RankingModel) SetHeight(h int) {
	r.height = h
}

// entries는 현재 탭에 해당하는 항목 리스트를 반환한다.
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

// NextTab은 다음 탭으로 이동한다.
func (r *RankingModel) NextTab() {
	r.tab = (r.tab + 1) % 3
	r.cursor = 0
	r.offset = 0
}

// SetTab은 탭을 직접 설정한다.
func (r *RankingModel) SetTab(tab usage.RankCategory) {
	r.tab = tab
	r.cursor = 0
	r.offset = 0
}

// ToggleScope는 범위를 전환하고 데이터를 다시 로드한다.
func (r *RankingModel) ToggleScope() {
	if r.scope == usage.ScopeAll {
		r.scope = usage.ScopeProject
	} else {
		r.scope = usage.ScopeAll
	}
	r.Load()
}

// MoveUp은 커서를 위로 이동한다.
func (r *RankingModel) MoveUp() {
	if r.cursor > 0 {
		r.cursor--
		r.adjustScroll()
	}
}

// MoveDown은 커서를 아래로 이동한다.
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

// View는 랭킹 뷰를 렌더링한다.
func (r *RankingModel) View(width, height int) string {
	r.height = height - 3 // 탭바 + 범위바 + 구분선

	var b strings.Builder

	// 탭바
	b.WriteString(r.renderTabs(width))
	b.WriteString("\n")

	// 범위바
	b.WriteString(r.renderScopeBar(width))
	b.WriteString("\n")

	// 구분선
	sep := lipgloss.NewStyle().Foreground(colorDimGray).Render(strings.Repeat("─", width-4))
	b.WriteString(sep)
	b.WriteString("\n")

	// 에러 표시
	if r.err != nil {
		errMsg := lipgloss.NewStyle().Foreground(colorRed).Render(fmt.Sprintf("오류: %v", r.err))
		b.WriteString(errMsg)
		return b.String()
	}

	// 데이터 없음
	entries := r.entries()
	if len(entries) == 0 {
		empty := lipgloss.NewStyle().Foreground(colorDimGray).Render("  데이터 없음")
		b.WriteString(empty)
		return b.String()
	}

	// 랭킹 리스트
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
	barWidth := width - 35 // 번호(4) + 등급(6) + 이름(15) + 카운트(6) + 여백(4)
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

	// 키 힌트를 오른쪽에 배치
	hint := hudDesc.Render("1/2/3: 탭  Tab: 다음")
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

	scopeBar := hudDesc.Render("범위: ") + allStyle.Render(" All ") + hudDesc.Render(" / ") + projStyle.Render(" Project ")

	hint := hudDesc.Render("s: 전환")
	pad := width - lipgloss.Width(scopeBar) - lipgloss.Width(hint) - 4
	if pad < 1 {
		pad = 1
	}
	return scopeBar + strings.Repeat(" ", pad) + hint
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
		// 선택된 항목: 배경색으로 확실하게 표시
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
