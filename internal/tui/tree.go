package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jeremy-kr/ccfg/internal/model"
)

// scopeStyle maps each scope to its emoji and color.
var scopeStyle = map[model.Scope]struct {
	emoji string
	color lipgloss.Color
}{
	model.ScopeManaged: {"🔒", colorRed},
	model.ScopeUser:    {"👤", colorGreen},
	model.ScopeProject: {"📁", colorCyan},
}

// categoryEmoji maps each config category to its emoji.
var categoryEmoji = map[model.ConfigCategory]string{
	model.CategorySettings:    "⚙️ ",
	model.CategoryInstructions: "📝",
	model.CategoryMCP:          "🔧",
	model.CategoryPolicy:       "🔑",
	model.CategoryCommands:     "⌨️ ",
	model.CategorySkills:       "🧠",
	model.CategoryAgents:       "🤖",
	model.CategoryKeybindings:  "🎮",
	model.CategoryHooks:        "🪝",
}

// TreeNode represents a single item in the tree.
type TreeNode struct {
	Label    string            // Display text.
	Scope    model.Scope       // Owning scope.
	File     *model.ConfigFile // nil for scope header nodes.
	Expanded bool              // Whether children are expanded.
	Children []TreeNode        // Child nodes.
}

// TreeModel manages the state of the left tree panel.
type TreeModel struct {
	roots  []TreeNode // Top-level nodes (one per scope).
	cursor int        // Currently selected visible index.
	offset int        // Scroll offset.
	height int        // Number of visible rows.
	filter string     // Search filter (empty string means no filter).
}

// NewTreeModel builds a tree from a ScanResult.
func NewTreeModel(result *model.ScanResult) TreeModel {
	var roots []TreeNode

	if len(result.Managed) > 0 {
		roots = append(roots, makeScopeNode("Managed", model.ScopeManaged, result.Managed))
	}
	if len(result.User) > 0 {
		roots = append(roots, makeScopeNode("User", model.ScopeUser, result.User))
	}
	if len(result.Project) > 0 {
		roots = append(roots, makeScopeNode("Project", model.ScopeProject, result.Project))
	}

	// Start with the first scope expanded.
	if len(roots) > 0 {
		roots[0].Expanded = true
	}

	return TreeModel{roots: roots}
}

func makeScopeNode(label string, scope model.Scope, files []model.ConfigFile) TreeNode {
	var children []TreeNode
	for _, f := range files {
		f := f
		children = append(children, makeFileNode(f, scope))
	}
	return TreeNode{
		Label:    label,
		Scope:    scope,
		Children: children,
	}
}

// makeFileNode recursively creates a TreeNode from a ConfigFile.
func makeFileNode(f model.ConfigFile, scope model.Scope) TreeNode {
	node := TreeNode{
		Label: f.Description,
		Scope: scope,
		File:  &f,
	}
	for _, child := range f.Children {
		child := child
		node.Children = append(node.Children, makeFileNode(child, scope))
	}
	return node
}

// visibleNodes returns the currently expanded nodes as a flat list.
func (t *TreeModel) visibleNodes() []TreeNode {
	var nodes []TreeNode
	filter := strings.ToLower(t.filter)

	for _, root := range t.roots {
		if filter != "" {
			// Filter mode: only show scopes with matching children.
			var matched []TreeNode
			for _, child := range root.Children {
				if strings.Contains(strings.ToLower(child.Label), filter) ||
					(child.File != nil && strings.Contains(strings.ToLower(child.File.Path), filter)) {
					matched = append(matched, child)
				}
			}
			if len(matched) > 0 {
				nodes = append(nodes, root)
				nodes = append(nodes, matched...)
			}
		} else {
			nodes = append(nodes, root)
			if root.Expanded {
				nodes = append(nodes, flattenExpanded(root.Children)...)
			}
		}
	}
	return nodes
}

// flattenExpanded recursively flattens expanded nodes into a flat list.
func flattenExpanded(nodes []TreeNode) []TreeNode {
	var flat []TreeNode
	for _, node := range nodes {
		flat = append(flat, node)
		if node.Expanded && len(node.Children) > 0 {
			flat = append(flat, flattenExpanded(node.Children)...)
		}
	}
	return flat
}

// Filter filters the tree by the given search text.
func (t *TreeModel) Filter(text string) {
	t.filter = text
	t.cursor = 0
	t.offset = 0
}

// ClearFilter clears the active filter.
func (t *TreeModel) ClearFilter() {
	t.filter = ""
}

// SelectedFile returns the file at the current cursor position. Returns nil for scope nodes.
func (t *TreeModel) SelectedFile() *model.ConfigFile {
	visible := t.visibleNodes()
	if t.cursor >= 0 && t.cursor < len(visible) {
		return visible[t.cursor].File
	}
	return nil
}

// SelectedScope returns the scope of the node at the current cursor position.
func (t *TreeModel) SelectedScope() model.Scope {
	visible := t.visibleNodes()
	if t.cursor >= 0 && t.cursor < len(visible) {
		return visible[t.cursor].Scope
	}
	return model.ScopeUser
}

// MoveUp moves the cursor up.
func (t *TreeModel) MoveUp() {
	if t.cursor > 0 {
		t.cursor--
		t.adjustScroll()
	}
}

// MoveDown moves the cursor down.
func (t *TreeModel) MoveDown() {
	visible := t.visibleNodes()
	if t.cursor < len(visible)-1 {
		t.cursor++
		t.adjustScroll()
	}
}

// Toggle expands or collapses a toggleable node (scope header or directory).
func (t *TreeModel) Toggle() {
	visible := t.visibleNodes()
	if t.cursor < 0 || t.cursor >= len(visible) {
		return
	}
	node := visible[t.cursor]

	// Toggle scope header node.
	if node.File == nil {
		for i := range t.roots {
			if t.roots[i].Label == node.Label {
				t.roots[i].Expanded = !t.roots[i].Expanded
				if !t.roots[i].Expanded {
					t.clampCursor()
				}
				return
			}
		}
		return
	}

	// Toggle file node (when it has children — directory or virtual group).
	if len(node.Children) > 0 {
		toggleByPath(t.roots, node.File.Path)
		t.clampCursor()
	}
}

// toggleByPath toggles the Expanded field of the node matching the given path.
func toggleByPath(nodes []TreeNode, path string) bool {
	for i := range nodes {
		if nodes[i].File != nil && nodes[i].File.Path == path {
			nodes[i].Expanded = !nodes[i].Expanded
			return true
		}
		if toggleByPath(nodes[i].Children, path) {
			return true
		}
	}
	return false
}

func (t *TreeModel) clampCursor() {
	visible := t.visibleNodes()
	if t.cursor >= len(visible) {
		t.cursor = len(visible) - 1
	}
	if t.cursor < 0 {
		t.cursor = 0
	}
}

func (t *TreeModel) adjustScroll() {
	if t.height <= 0 {
		return
	}
	if t.cursor < t.offset {
		t.offset = t.cursor
	}
	if t.cursor >= t.offset+t.height {
		t.offset = t.cursor - t.height + 1
	}
}

// SetHeight sets the number of visible rows.
func (t *TreeModel) SetHeight(h int) {
	t.height = h
	t.adjustScroll()
}

// View renders the tree as a string.
func (t *TreeModel) View(width int, focused bool) string {
	visible := t.visibleNodes()
	var b strings.Builder

	end := t.offset + t.height
	if end > len(visible) {
		end = len(visible)
	}

	base := panelStyleFor(focused)
	style := base.Width(width - base.GetHorizontalBorderSize()).Height(t.height)
	availWidth := width - style.GetHorizontalFrameSize()

	scrollBars := renderScrollbar(len(visible), t.height, t.offset)
	contentW := availWidth
	if scrollBars != nil {
		contentW = availWidth - 1
	}

	for i := t.offset; i < end; i++ {
		node := visible[i]
		line := t.renderNode(node, i == t.cursor, focused)
		if scrollBars != nil {
			line = lipgloss.NewStyle().MaxWidth(contentW).Render(line)
			if gap := contentW - lipgloss.Width(line); gap > 0 {
				line += strings.Repeat(" ", gap)
			}
			line += scrollBars[i-t.offset]
		}
		b.WriteString(line)
		if i < end-1 {
			b.WriteString("\n")
		}
	}

	// Fill remaining rows with empty lines.
	rendered := end - t.offset
	for i := rendered; i < t.height; i++ {
		if i > 0 {
			b.WriteString("\n")
		}
	}

	content := lipgloss.NewStyle().MaxWidth(availWidth).Render(b.String())
	return style.Render(content)
}

func (t *TreeModel) renderNode(node TreeNode, selected, focused bool) string {
	if node.File == nil {
		// Scope header.
		arrow := "▶"
		for _, r := range t.roots {
			if r.Label == node.Label && r.Expanded {
				arrow = "▼"
				break
			}
		}

		ss, ok := scopeStyle[node.Scope]
		if !ok {
			ss = scopeStyle[model.ScopeUser]
		}
		text := fmt.Sprintf("%s %s %s", ss.emoji, arrow, strings.ToUpper(node.Label))
		style := scopeHeaderStyle.Foreground(ss.color)
		if selected && focused {
			return treeSelectedStyle.Render(text)
		}
		return style.Render(text)
	}

	depth := nodeDepth(t.roots, node.File)
	indent := strings.Repeat("  ", depth)

	// Category emoji.
	emoji := ""
	if e, ok := categoryEmoji[node.File.Category]; ok {
		emoji = e + " "
	}

	// Expandable node (directory or virtual group with children).
	if len(node.Children) > 0 {
		arrow := "▶"
		if node.Expanded {
			arrow = "▼"
		}
		count := fmt.Sprintf("(%d)", len(node.Children))
		text := fmt.Sprintf("%s%s %s%s %s", indent, arrow, emoji, node.Label, count)
		if selected && focused {
			return treeSelectedStyle.Render(text)
		}
		return dirStyle.Render(text)
	}

	// Virtual leaf node (individual item from a JSON internal section).
	if node.File.IsVirtual {
		text := fmt.Sprintf("%s%s%s", indent, emoji, node.Label)
		if selected && focused {
			return treeSelectedStyle.Render(text)
		}
		return treeItemStyle.Render(text)
	}

	// File node.
	if selected && focused {
		text := fmt.Sprintf("%s▸ %s%s", indent, emoji, node.Label)
		return treeSelectedStyle.Render(text)
	}

	status := fileMissingStyle.Render("○")
	if node.File.Exists {
		status = fileExistsStyle.Render("●")
	}
	text := fmt.Sprintf("%s%s %s%s", indent, status, emoji, node.Label)
	return treeItemStyle.Render(text)
}

// nodeDepth returns the depth of the given file node in the tree.
func nodeDepth(roots []TreeNode, target *model.ConfigFile) int {
	if target == nil {
		return 0
	}
	for _, root := range roots {
		if d := findDepth(root.Children, target.Path, 1); d > 0 {
			return d
		}
	}
	return 1
}

func findDepth(nodes []TreeNode, path string, depth int) int {
	for _, node := range nodes {
		if node.File != nil && node.File.Path == path {
			return depth
		}
		if d := findDepth(node.Children, path, depth+1); d > 0 {
			return d
		}
	}
	return 0
}

// TreeState is a snapshot of the tree's expansion state and cursor position.
type TreeState struct {
	Expanded     map[string]bool // key -> expanded ("__scope__"+label for scope headers, path for files).
	SelectedPath string          // Currently selected file path.
}

// CaptureState captures the current expansion state and cursor position of the tree.
func (t *TreeModel) CaptureState() TreeState {
	expanded := make(map[string]bool)
	for _, root := range t.roots {
		expanded["__scope__"+root.Label] = root.Expanded
		captureExpanded(root.Children, expanded)
	}

	var selectedPath string
	if f := t.SelectedFile(); f != nil {
		selectedPath = f.Path
	}

	return TreeState{
		Expanded:     expanded,
		SelectedPath: selectedPath,
	}
}

func captureExpanded(nodes []TreeNode, out map[string]bool) {
	for _, node := range nodes {
		if node.File != nil && len(node.Children) > 0 {
			out[node.File.Path] = node.Expanded
		}
		captureExpanded(node.Children, out)
	}
}

// RestoreState restores a previously captured state to the current tree.
func (t *TreeModel) RestoreState(state TreeState) {
	for i := range t.roots {
		if v, ok := state.Expanded["__scope__"+t.roots[i].Label]; ok {
			t.roots[i].Expanded = v
		}
		restoreExpanded(t.roots[i].Children, state.Expanded)
	}

	// Restore cursor position.
	if state.SelectedPath != "" {
		visible := t.visibleNodes()
		for i, node := range visible {
			if node.File != nil && node.File.Path == state.SelectedPath {
				t.cursor = i
				t.adjustScroll()
				return
			}
		}
	}
	t.clampCursor()
}

func restoreExpanded(nodes []TreeNode, expanded map[string]bool) {
	for i := range nodes {
		if nodes[i].File != nil && len(nodes[i].Children) > 0 {
			if v, ok := expanded[nodes[i].File.Path]; ok {
				nodes[i].Expanded = v
			}
		}
		restoreExpanded(nodes[i].Children, expanded)
	}
}
