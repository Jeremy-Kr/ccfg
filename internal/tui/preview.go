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

// PreviewModel manages the state of the right preview panel.
type PreviewModel struct {
	file       *model.ConfigFile // Currently displayed file.
	content    string            // File content.
	lines      []string          // Content split by line.
	offset     int               // Scroll offset.
	height     int               // Number of visible rows.
	isCardMode bool              // Card mode (agents/skills directory).
	lastWidth  int               // Last width used in card mode.
}

// SetFile sets the file to display in the preview.
func (p *PreviewModel) SetFile(file *model.ConfigFile) {
	if file == nil {
		p.file = nil
		p.content = ""
		p.lines = nil
		p.offset = 0
		p.isCardMode = false
		return
	}

	// Skip if the same file is already displayed.
	if p.file != nil && p.file.Path == file.Path {
		return
	}

	p.file = file
	p.offset = 0
	p.isCardMode = false

	// Virtual node — JSON internal section preview.
	if file.IsVirtual {
		p.content = p.renderVirtualNode(file)
		p.lines = strings.Split(p.content, "\n")
		return
	}

	if !file.Exists {
		p.content = "(file not found)"
		p.lines = []string{p.content}
		return
	}

	// Directory case.
	if file.IsDir {
		// agents/skills directories use card mode.
		if file.Category == model.CategoryAgents || file.Category == model.CategorySkills {
			p.isCardMode = true
			// Generate cards immediately if lastWidth is available, otherwise wait for PrepareCardContent.
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
		p.content = fmt.Sprintf("(failed to read: %v)", err)
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
	b.WriteString(fmt.Sprintf("Directory: %s\n\n", file.Path))

	if len(file.Children) == 0 {
		entries, err := os.ReadDir(file.Path)
		if err != nil {
			b.WriteString(fmt.Sprintf("(failed to read: %v)", err))
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

// dirIcon returns a directory or file icon based on whether it is a directory.
func dirIcon(isDir bool) string {
	if isDir {
		return "📁"
	}
	return "📄"
}

// InvalidateCache invalidates the cached file so the next SetFile call forces a refresh.
func (p *PreviewModel) InvalidateCache() { p.file = nil }

// ScrollUp scrolls the preview up by n lines.
func (p *PreviewModel) ScrollUp(n int) {
	p.offset -= n
	if p.offset < 0 {
		p.offset = 0
	}
}

// ScrollDown scrolls the preview down by n lines.
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

// SetHeight sets the number of visible rows.
func (p *PreviewModel) SetHeight(h int) {
	p.height = h
}

// PrepareCardContent pre-generates card lines at the given width for card mode.
// Must be called from the Update() flow (View() uses a value receiver so state is not persisted).
func (p *PreviewModel) PrepareCardContent(width int) {
	if !p.isCardMode || p.file == nil || width <= 0 {
		return
	}
	if p.lastWidth == width && p.lines != nil {
		return // Already generated at the same width.
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

// View renders the preview as a string.
func (p *PreviewModel) View(width int, focused bool) string {
	var b strings.Builder
	availW := width - panelStyle.GetHorizontalFrameSize()

	if p.file == nil {
		b.WriteString("Select a file")
	} else if p.isCardMode {
		// Card mode: use p.lines pre-generated by PrepareCardContent().
		renderScrollableLines(&b, p.lines, p.height, p.offset, availW)
	} else {
		// File info header (decorated line).
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

		// Display content (minus 1 row for header).
		renderScrollableLines(&b, p.lines, p.height-1, p.offset, availW)
	}

	// Fix panel height and prevent line wrapping.
	base := panelStyleFor(focused)
	style := base.Width(width - base.GetHorizontalBorderSize()).Height(p.height)
	availWidth := width - style.GetHorizontalFrameSize()
	content := lipgloss.NewStyle().MaxWidth(availWidth).Render(b.String())

	return style.Render(content)
}

// renderScrollableLines renders lines with a scrollbar and writes to b.
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

// renderAgentCards renders .md files in the agents directory as character cards.
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
			return fmt.Sprintf("(failed to read: %v)", err)
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
		return "(no agent files)"
	}
	return strings.Join(cards, "\n")
}

// renderSkillCards renders SKILL.md files in the skills directory as ability cards.
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
			return fmt.Sprintf("(failed to read: %v)", err)
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
		return "(no skill files)"
	}
	return strings.Join(cards, "\n")
}

// renderAgentCard renders a single agent character card.
// width is the total card box width (including border).
func renderAgentCard(meta *parser.AgentMeta, width int) string {
	var lines []string

	// lipgloss Width(w) word-wraps at w - padding.
	// So contentW = width - border - padding is the actual content width.
	borderW := agentCardStyle.GetHorizontalBorderSize()
	paddingW := agentCardStyle.GetHorizontalFrameSize() - borderW
	contentW := width - borderW - paddingW

	// Title line: 🤖 name.
	title := agentCardTitleStyle.Render("🤖 " + meta.Name)
	lines = append(lines, title)

	// Role separator line.
	roleLine := "━━"
	if meta.Role != "" {
		roleLine += " " + meta.Role + " "
	}
	if pad := contentW - lipgloss.Width(roleLine); pad > 0 {
		roleLine += strings.Repeat("━", pad)
	}
	lines = append(lines, agentCardRoleStyle.Render(roleLine))

	// Description.
	if meta.Desc != "" {
		lines = append(lines, "")
		descLines := wrapText(meta.Desc, contentW)
		lines = append(lines, descLines...)
		lines = append(lines, "")
	}

	// Meta info (model, color).
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
	// Width = width - borderW sets the inner (padding+content) width. Total render width = width.
	return agentCardStyle.Width(width - borderW).Render(content)
}

// renderSkillCard renders a single skill ability card.
// width is the total card box width (including border).
func renderSkillCard(meta *parser.SkillMeta, width int) string {
	var lines []string

	borderW := skillCardStyle.GetHorizontalBorderSize()
	paddingW := skillCardStyle.GetHorizontalFrameSize() - borderW
	contentW := width - borderW - paddingW

	// Title line: ⚡ name      [category].
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

	// Separator line.
	sep := strings.Repeat("━", contentW)
	lines = append(lines, lipgloss.NewStyle().Foreground(colorCyan).Render(sep))

	// Description.
	if meta.Desc != "" {
		lines = append(lines, "")
		descLines := wrapText(meta.Desc, contentW)
		lines = append(lines, descLines...)
		lines = append(lines, "")
	}

	// Tags.
	if meta.Tags != "" {
		tagLine := lipgloss.NewStyle().Foreground(colorGreen).Render("🎯 " + meta.Tags)
		lines = append(lines, tagLine)
	}

	content := strings.Join(lines, "\n")
	return skillCardStyle.Width(width - borderW).Render(content)
}

// wrapText wraps text to fit within the given width.
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

// renderVirtualNode renders the preview of a virtual node (JSON internal section).
func (p *PreviewModel) renderVirtualNode(file *model.ConfigFile) string {
	parts := strings.SplitN(file.Path, "#", 2)
	if len(parts) != 2 {
		return "(failed to parse virtual node path)"
	}
	realPath := parts[0]
	dotPath := parts[1]

	data, err := os.ReadFile(realPath)
	if err != nil {
		return fmt.Sprintf("(failed to read: %v)", err)
	}

	cleaned := parser.StripJSONC(string(data))
	var obj any
	if err := json.Unmarshal([]byte(cleaned), &obj); err != nil {
		return fmt.Sprintf("(failed to parse JSON: %v)", err)
	}

	section := navigateJSON(obj, dotPath)
	if section == nil {
		return fmt.Sprintf("(section not found: %s)", dotPath)
	}

	// FormatJSON internally handles pretty-print and syntax highlighting.
	sectionBytes, err := json.Marshal(section)
	if err != nil {
		return fmt.Sprintf("%v", section)
	}
	return parser.FormatJSON(string(sectionBytes))
}

// navigateJSON navigates a JSON object using dot notation (dotPath).
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
