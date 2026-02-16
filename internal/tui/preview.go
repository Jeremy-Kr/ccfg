package tui

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jeremy-kr/ccfg/internal/model"
	"github.com/jeremy-kr/ccfg/internal/parser"
)

// PreviewModel은 우측 미리보기 패널의 상태를 관리한다.
type PreviewModel struct {
	file       *model.ConfigFile // 현재 표시 중인 파일
	content    string            // 파일 내용
	lines      []string          // 줄 단위 분할
	offset     int               // 스크롤 오프셋
	height     int               // 표시 가능한 행 수
	isCardMode bool              // 카드 모드 (agents/skills 디렉토리)
	lastWidth  int               // 카드 모드에서 사용한 마지막 폭
}

// SetFile은 미리보기에 표시할 파일을 설정한다.
func (p *PreviewModel) SetFile(file *model.ConfigFile) {
	if file == nil {
		p.file = nil
		p.content = ""
		p.lines = nil
		p.offset = 0
		p.isCardMode = false
		return
	}

	// 이미 같은 파일이면 스킵
	if p.file != nil && p.file.Path == file.Path {
		return
	}

	p.file = file
	p.offset = 0
	p.isCardMode = false

	// 가상 노드 — JSON 내부 섹션 미리보기
	if file.IsVirtual {
		p.content = p.renderVirtualNode(file)
		p.lines = strings.Split(p.content, "\n")
		return
	}

	if !file.Exists {
		p.content = "(파일이 존재하지 않습니다)"
		p.lines = []string{p.content}
		return
	}

	// 디렉토리인 경우
	if file.IsDir {
		// agents/skills 디렉토리는 카드 모드
		if file.Category == model.CategoryAgents || file.Category == model.CategorySkills {
			p.isCardMode = true
			// lastWidth가 있으면 즉시 카드 생성, 없으면 PrepareCardContent 호출 대기
			if p.lastWidth > 0 {
				p.generateCardLines(p.lastWidth)
			} else {
				p.content = ""
				p.lines = nil
			}
			return
		}
		p.content = p.renderDir(file)
		p.lines = strings.Split(p.content, "\n")
		return
	}

	data, err := os.ReadFile(file.Path)
	if err != nil {
		p.content = fmt.Sprintf("(읽기 실패: %v)", err)
		p.lines = []string{p.content}
		return
	}

	raw := string(data)
	switch file.FileType {
	case model.FileTypeJSON, model.FileTypeJSONC:
		p.content = parser.FormatJSON(raw)
	case model.FileTypeMarkdown:
		p.content = parser.FormatMarkdown(raw)
	default:
		p.content = raw
	}
	p.lines = strings.Split(p.content, "\n")
}

func (p *PreviewModel) renderDir(file *model.ConfigFile) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("디렉토리: %s\n\n", file.Path))

	if len(file.Children) == 0 {
		entries, err := os.ReadDir(file.Path)
		if err != nil {
			b.WriteString(fmt.Sprintf("(읽기 실패: %v)", err))
			return b.String()
		}
		for _, entry := range entries {
			icon := dirIcon(entry.IsDir())
			b.WriteString(fmt.Sprintf("  %s %s\n", icon, entry.Name()))
		}
	} else {
		for _, child := range file.Children {
			icon := dirIcon(child.IsDir)
			detail := ""
			if child.Exists {
				detail = fmt.Sprintf("  (%d bytes)", child.Size)
			}
			b.WriteString(fmt.Sprintf("  %s %s%s\n", icon, child.Description, detail))
		}
	}

	return b.String()
}

// dirIcon은 디렉토리 여부에 따라 아이콘을 반환한다.
func dirIcon(isDir bool) string {
	if isDir {
		return "📁"
	}
	return "📄"
}

// InvalidateCache는 현재 캐시된 파일을 무효화하여 다음 SetFile 호출 시 강제 갱신한다.
func (p *PreviewModel) InvalidateCache() { p.file = nil }

// ScrollUp은 미리보기를 위로 스크롤한다.
func (p *PreviewModel) ScrollUp(n int) {
	p.offset -= n
	if p.offset < 0 {
		p.offset = 0
	}
}

// ScrollDown은 미리보기를 아래로 스크롤한다.
func (p *PreviewModel) ScrollDown(n int) {
	maxOffset := len(p.lines) - p.height
	if maxOffset < 0 {
		maxOffset = 0
	}
	p.offset += n
	if p.offset > maxOffset {
		p.offset = maxOffset
	}
}

// SetHeight는 표시 가능한 행 수를 설정한다.
func (p *PreviewModel) SetHeight(h int) {
	p.height = h
}

// PrepareCardContent는 카드 모드일 때 주어진 폭으로 카드 lines를 미리 생성한다.
// Update() 흐름에서 호출해야 한다 (View()는 value receiver라 상태가 유지되지 않음).
func (p *PreviewModel) PrepareCardContent(width int) {
	if !p.isCardMode || p.file == nil || width <= 0 {
		return
	}
	if p.lastWidth == width && p.lines != nil {
		return // 이미 같은 폭으로 생성됨
	}
	p.generateCardLines(width)
}

func (p *PreviewModel) generateCardLines(width int) {
	p.lastWidth = width
	availW := width - panelStyle.GetHorizontalFrameSize()
	cardW := max(availW-2, 20)

	var cardContent string
	switch p.file.Category {
	case model.CategoryAgents:
		cardContent = p.renderAgentCards(p.file, cardW)
	case model.CategorySkills:
		cardContent = p.renderSkillCards(p.file, cardW)
	}

	p.content = cardContent
	p.lines = strings.Split(cardContent, "\n")
}

// View는 미리보기를 문자열로 렌더링한다.
func (p *PreviewModel) View(width int, focused bool) string {
	var b strings.Builder
	availW := width - panelStyle.GetHorizontalFrameSize()

	if p.file == nil {
		b.WriteString("파일을 선택하세요")
	} else if p.isCardMode {
		// 카드 모드: PrepareCardContent()에서 미리 생성된 p.lines 사용
		renderScrollableLines(&b, p.lines, p.height, p.offset, availW)
	} else {
		// 파일 정보 헤더 (장식 라인)
		icon := dirIcon(p.file.IsDir)
		info := p.file.Path
		if p.file.Exists && !p.file.IsDir {
			info = fmt.Sprintf("%s (%d bytes)", p.file.Path, p.file.Size)
		}
		label := fmt.Sprintf("[ %s %s ]", icon, info)
		pad := max(availW-lipgloss.Width(label), 2)
		left := pad / 2
		right := pad - left
		decoratedHeader := strings.Repeat("━", left) + label + strings.Repeat("━", right)
		b.WriteString(lipgloss.NewStyle().Foreground(colorCyan).Render(decoratedHeader))
		b.WriteString("\n")

		// 내용 표시 (헤더 1줄 제외)
		renderScrollableLines(&b, p.lines, p.height-1, p.offset, availW)
	}

	// 패널 높이 고정 + 줄바꿈 방지
	base := panelStyleFor(focused)
	style := base.Width(width - base.GetHorizontalBorderSize()).Height(p.height)
	availWidth := width - style.GetHorizontalFrameSize()
	content := lipgloss.NewStyle().MaxWidth(availWidth).Render(b.String())

	return style.Render(content)
}

// renderScrollableLines는 lines를 스크롤바와 함께 렌더링하여 b에 기록한다.
func renderScrollableLines(b *strings.Builder, lines []string, visibleRows, offset, availW int) {
	end := offset + visibleRows
	if end > len(lines) {
		end = len(lines)
	}

	scrollBars := renderScrollbar(len(lines), visibleRows, offset)

	if scrollBars != nil {
		contentW := availW - 1
		for i := offset; i < end; i++ {
			line := lipgloss.NewStyle().MaxWidth(contentW).Render(lines[i])
			if gap := contentW - lipgloss.Width(line); gap > 0 {
				line += strings.Repeat(" ", gap)
			}
			line += scrollBars[i-offset]
			b.WriteString(line)
			if i < end-1 {
				b.WriteString("\n")
			}
		}
	} else {
		for i := offset; i < end; i++ {
			b.WriteString(lines[i])
			if i < end-1 {
				b.WriteString("\n")
			}
		}
	}
}

// renderAgentCards는 에이전트 디렉토리의 .md 파일들을 캐릭터 카드로 렌더링한다.
func (p *PreviewModel) renderAgentCards(file *model.ConfigFile, width int) string {
	var cards []string

	if len(file.Children) > 0 {
		for _, child := range file.Children {
			if child.IsDir || !child.Exists {
				continue
			}
			meta := parser.ParseAgentMeta(child.Path)
			if meta != nil {
				cards = append(cards, renderAgentCard(meta, width))
			}
		}
	} else {
		entries, err := os.ReadDir(file.Path)
		if err != nil {
			return fmt.Sprintf("(읽기 실패: %v)", err)
		}
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
				continue
			}
			meta := parser.ParseAgentMeta(filepath.Join(file.Path, entry.Name()))
			if meta != nil {
				cards = append(cards, renderAgentCard(meta, width))
			}
		}
	}

	if len(cards) == 0 {
		return "(에이전트 파일 없음)"
	}
	return strings.Join(cards, "\n")
}

// renderSkillCards는 스킬 디렉토리의 SKILL.md 파일들을 어빌리티 카드로 렌더링한다.
func (p *PreviewModel) renderSkillCards(file *model.ConfigFile, width int) string {
	var cards []string

	if len(file.Children) > 0 {
		for _, child := range file.Children {
			if !child.IsDir || !child.Exists {
				continue
			}
			skillPath := filepath.Join(child.Path, "SKILL.md")
			meta := parser.ParseSkillMeta(skillPath)
			if meta != nil {
				cards = append(cards, renderSkillCard(meta, width))
			}
		}
	} else {
		entries, err := os.ReadDir(file.Path)
		if err != nil {
			return fmt.Sprintf("(읽기 실패: %v)", err)
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			skillPath := filepath.Join(file.Path, entry.Name(), "SKILL.md")
			meta := parser.ParseSkillMeta(skillPath)
			if meta != nil {
				cards = append(cards, renderSkillCard(meta, width))
			}
		}
	}

	if len(cards) == 0 {
		return "(스킬 파일 없음)"
	}
	return strings.Join(cards, "\n")
}

// renderAgentCard는 개별 에이전트 캐릭터 카드를 렌더링한다.
// width는 카드 박스의 총 폭 (border 포함).
func renderAgentCard(meta *parser.AgentMeta, width int) string {
	var lines []string

	// lipgloss Width(w)는 w - padding 에서 word wrap.
	// 따라서 contentW = width - border - padding 이 실제 콘텐츠 폭.
	borderW := agentCardStyle.GetHorizontalBorderSize()
	paddingW := agentCardStyle.GetHorizontalFrameSize() - borderW
	contentW := width - borderW - paddingW

	// 타이틀 라인: 🤖 name
	title := agentCardTitleStyle.Render("🤖 " + meta.Name)
	lines = append(lines, title)

	// 역할 구분선
	roleLine := "━━"
	if meta.Role != "" {
		roleLine += " " + meta.Role + " "
	}
	if pad := contentW - lipgloss.Width(roleLine); pad > 0 {
		roleLine += strings.Repeat("━", pad)
	}
	lines = append(lines, agentCardRoleStyle.Render(roleLine))

	// 설명
	if meta.Desc != "" {
		lines = append(lines, "")
		descLines := wrapText(meta.Desc, contentW)
		lines = append(lines, descLines...)
		lines = append(lines, "")
	}

	// 메타 정보 (model, color)
	var metaParts []string
	if meta.Model != "" {
		metaParts = append(metaParts, "🧠 "+meta.Model)
	}
	if meta.Color != "" {
		metaParts = append(metaParts, "🎨 "+meta.Color)
	}
	if len(metaParts) > 0 {
		metaLine := lipgloss.NewStyle().Foreground(colorGreen).Render(strings.Join(metaParts, "   "))
		lines = append(lines, metaLine)
	}

	content := strings.Join(lines, "\n")
	// Width = width - borderW → 내부(padding+content) 폭 설정. 총 렌더 폭 = width.
	return agentCardStyle.Width(width - borderW).Render(content)
}

// renderSkillCard는 개별 스킬 어빌리티 카드를 렌더링한다.
// width는 카드 박스의 총 폭 (border 포함).
func renderSkillCard(meta *parser.SkillMeta, width int) string {
	var lines []string

	borderW := skillCardStyle.GetHorizontalBorderSize()
	paddingW := skillCardStyle.GetHorizontalFrameSize() - borderW
	contentW := width - borderW - paddingW

	// 타이틀 라인: ⚡ name      [category]
	titlePart := skillCardTitleStyle.Render("⚡ " + meta.Name)
	if meta.Category != "" {
		tag := skillCardTagStyle.Render("[" + meta.Category + "]")
		gap := contentW - lipgloss.Width(titlePart) - lipgloss.Width(tag)
		if gap < 1 {
			gap = 1
		}
		titlePart += strings.Repeat(" ", gap) + tag
	}
	lines = append(lines, titlePart)

	// 구분선
	sep := strings.Repeat("━", contentW)
	lines = append(lines, lipgloss.NewStyle().Foreground(colorCyan).Render(sep))

	// 설명
	if meta.Desc != "" {
		lines = append(lines, "")
		descLines := wrapText(meta.Desc, contentW)
		lines = append(lines, descLines...)
		lines = append(lines, "")
	}

	// 태그
	if meta.Tags != "" {
		tagLine := lipgloss.NewStyle().Foreground(colorGreen).Render("🎯 " + meta.Tags)
		lines = append(lines, tagLine)
	}

	content := strings.Join(lines, "\n")
	return skillCardStyle.Width(width - borderW).Render(content)
}

// wrapText는 텍스트를 주어진 폭에 맞게 줄바꿈한다.
func wrapText(text string, width int) []string {
	if width <= 0 {
		return []string{text}
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return nil
	}

	var lines []string
	current := words[0]

	for _, word := range words[1:] {
		test := current + " " + word
		if lipgloss.Width(test) <= width {
			current = test
		} else {
			lines = append(lines, current)
			current = word
		}
	}
	lines = append(lines, current)
	return lines
}

// renderVirtualNode는 가상 노드(JSON 내부 섹션)의 미리보기를 렌더링한다.
func (p *PreviewModel) renderVirtualNode(file *model.ConfigFile) string {
	parts := strings.SplitN(file.Path, "#", 2)
	if len(parts) != 2 {
		return "(가상 노드 경로 파싱 실패)"
	}
	realPath := parts[0]
	dotPath := parts[1]

	data, err := os.ReadFile(realPath)
	if err != nil {
		return fmt.Sprintf("(읽기 실패: %v)", err)
	}

	cleaned := parser.StripJSONC(string(data))
	var obj any
	if err := json.Unmarshal([]byte(cleaned), &obj); err != nil {
		return fmt.Sprintf("(JSON 파싱 실패: %v)", err)
	}

	section := navigateJSON(obj, dotPath)
	if section == nil {
		return fmt.Sprintf("(섹션을 찾을 수 없습니다: %s)", dotPath)
	}

	// FormatJSON이 내부적으로 pretty-print + 구문 강조를 처리한다.
	sectionBytes, err := json.Marshal(section)
	if err != nil {
		return fmt.Sprintf("%v", section)
	}
	return parser.FormatJSON(string(sectionBytes))
}

// navigateJSON은 점 표기법(dotPath)으로 JSON 객체를 탐색한다.
func navigateJSON(obj any, dotPath string) any {
	keys := strings.Split(dotPath, ".")
	current := obj
	for _, k := range keys {
		m, ok := current.(map[string]any)
		if !ok {
			return nil
		}
		v, exists := m[k]
		if !exists {
			return nil
		}
		current = v
	}
	return current
}
