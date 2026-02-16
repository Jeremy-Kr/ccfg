package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jeremy-kr/ccfg/internal/model"
)

// TreeNode는 트리의 한 항목을 나타낸다.
type TreeNode struct {
	Label    string            // 표시 텍스트
	Scope    model.Scope       // 소속 Scope
	File     *model.ConfigFile // nil이면 Scope 헤더 노드
	Expanded bool              // 자식이 펼쳐져 있는지
	Children []TreeNode        // 하위 노드
}

// TreeModel은 좌측 트리 패널의 상태를 관리한다.
type TreeModel struct {
	roots  []TreeNode // 최상위 노드 (Scope별)
	cursor int        // 현재 선택된 visible 인덱스
	offset int        // 스크롤 오프셋
	height int        // 표시 가능한 행 수
	filter string     // 검색 필터 (빈 문자열이면 필터 없음)
}

// NewTreeModel은 ScanResult로부터 트리를 구성한다.
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

	// 첫 번째 Scope를 펼친 상태로 시작
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

// makeFileNode은 ConfigFile로부터 TreeNode를 재귀적으로 생성한다.
func makeFileNode(f model.ConfigFile, scope model.Scope) TreeNode {
	node := TreeNode{
		Label: f.Description,
		Scope: scope,
		File:  &f,
	}
	if f.IsDir && len(f.Children) > 0 {
		for _, child := range f.Children {
			child := child
			node.Children = append(node.Children, makeFileNode(child, scope))
		}
	}
	return node
}

// visibleNodes는 현재 펼쳐진 노드들을 플랫 리스트로 반환한다.
func (t *TreeModel) visibleNodes() []TreeNode {
	var nodes []TreeNode
	filter := strings.ToLower(t.filter)

	for _, root := range t.roots {
		if filter != "" {
			// 필터 모드: 매칭되는 자식이 있는 Scope만 표시
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

// flattenExpanded는 펼쳐진 노드를 재귀적으로 플랫 리스트로 변환한다.
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

// Filter는 트리를 검색어로 필터링한다.
func (t *TreeModel) Filter(text string) {
	t.filter = text
	t.cursor = 0
	t.offset = 0
}

// ClearFilter는 필터를 해제한다.
func (t *TreeModel) ClearFilter() {
	t.filter = ""
}

// SelectedFile은 현재 커서가 가리키는 파일을 반환한다. Scope 노드면 nil.
func (t *TreeModel) SelectedFile() *model.ConfigFile {
	visible := t.visibleNodes()
	if t.cursor >= 0 && t.cursor < len(visible) {
		return visible[t.cursor].File
	}
	return nil
}

// MoveUp은 커서를 위로 이동한다.
func (t *TreeModel) MoveUp() {
	if t.cursor > 0 {
		t.cursor--
		t.adjustScroll()
	}
}

// MoveDown은 커서를 아래로 이동한다.
func (t *TreeModel) MoveDown() {
	visible := t.visibleNodes()
	if t.cursor < len(visible)-1 {
		t.cursor++
		t.adjustScroll()
	}
}

// Toggle은 펼칠 수 있는 노드(Scope 헤더, 디렉토리)를 펼치거나 접는다.
func (t *TreeModel) Toggle() {
	visible := t.visibleNodes()
	if t.cursor < 0 || t.cursor >= len(visible) {
		return
	}
	node := visible[t.cursor]

	// Scope 헤더 노드 토글
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

	// 디렉토리 파일 노드 토글 (Children이 있는 경우)
	if node.File.IsDir && len(node.Children) > 0 {
		toggleByPath(t.roots, node.File.Path)
		t.clampCursor()
	}
}

// toggleByPath는 트리에서 경로가 일치하는 노드의 Expanded를 토글한다.
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

// SetHeight는 표시 가능한 행 수를 설정한다.
func (t *TreeModel) SetHeight(h int) {
	t.height = h
	t.adjustScroll()
}

// View는 트리를 문자열로 렌더링한다.
func (t *TreeModel) View(width int, focused bool) string {
	visible := t.visibleNodes()
	var b strings.Builder

	end := t.offset + t.height
	if end > len(visible) {
		end = len(visible)
	}

	for i := t.offset; i < end; i++ {
		node := visible[i]
		line := t.renderNode(node, i == t.cursor, focused)
		// 너비에 맞게 자르기
		if lipgloss.Width(line) > width {
			line = line[:width]
		}
		b.WriteString(line)
		if i < end-1 {
			b.WriteString("\n")
		}
	}

	// 남는 행은 빈 줄로 채우기
	rendered := end - t.offset
	for i := rendered; i < t.height; i++ {
		if i > 0 {
			b.WriteString("\n")
		}
	}

	style := panelStyle.Width(width)
	if focused {
		style = panelFocusedStyle.Width(width)
	}

	return style.Render(b.String())
}

func (t *TreeModel) renderNode(node TreeNode, selected, focused bool) string {
	if node.File == nil {
		// Scope 헤더
		arrow := "▸"
		for _, r := range t.roots {
			if r.Label == node.Label && r.Expanded {
				arrow = "▾"
				break
			}
		}
		text := fmt.Sprintf("%s %s", arrow, node.Label)
		if selected && focused {
			return treeSelectedStyle.Render(text)
		}
		return scopeHeaderStyle.Render(text)
	}

	depth := nodeDepth(t.roots, node.File)
	indent := strings.Repeat("  ", depth)

	// 디렉토리 노드 (펼침 가능)
	if node.File.IsDir && len(node.Children) > 0 {
		arrow := "▸"
		if node.Expanded {
			arrow = "▾"
		}
		count := fmt.Sprintf("(%d)", len(node.Children))
		text := fmt.Sprintf("%s%s %s %s", indent, arrow, node.Label, count)
		if selected && focused {
			return treeSelectedStyle.Render(text)
		}
		if node.File.Exists {
			return fileExistsStyle.Render(text)
		}
		return fileMissingStyle.Render(text)
	}

	// 파일 노드
	status := fileMissingStyle.Render("✗")
	if node.File.Exists {
		status = fileExistsStyle.Render("✓")
	}
	text := fmt.Sprintf("%s%s %s", indent, status, node.Label)
	if selected && focused {
		return treeSelectedStyle.Render(fmt.Sprintf("%s%s %s", indent, "›", node.Label))
	}
	return treeItemStyle.Render(text)
}

// nodeDepth는 트리에서 해당 파일 노드의 깊이를 반환한다.
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
