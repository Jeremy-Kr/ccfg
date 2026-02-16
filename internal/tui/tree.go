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
	children := make([]TreeNode, len(files))
	for i, f := range files {
		f := f
		children[i] = TreeNode{
			Label: f.Description,
			Scope: scope,
			File:  &f,
		}
	}
	return TreeNode{
		Label:    label,
		Scope:    scope,
		Children: children,
	}
}

// visibleNodes는 현재 펼쳐진 노드들을 플랫 리스트로 반환한다.
func (t *TreeModel) visibleNodes() []TreeNode {
	var nodes []TreeNode
	for _, root := range t.roots {
		nodes = append(nodes, root)
		if root.Expanded {
			nodes = append(nodes, root.Children...)
		}
	}
	return nodes
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

// Toggle은 현재 Scope 노드를 펼치거나 접는다.
func (t *TreeModel) Toggle() {
	visible := t.visibleNodes()
	if t.cursor < 0 || t.cursor >= len(visible) {
		return
	}
	node := visible[t.cursor]

	// Scope 헤더 노드만 토글 가능
	if node.File != nil {
		return
	}

	for i := range t.roots {
		if t.roots[i].Label == node.Label {
			t.roots[i].Expanded = !t.roots[i].Expanded
			// 접을 때 커서가 자식에 있었으면 부모로 이동
			if !t.roots[i].Expanded {
				t.clampCursor()
			}
			break
		}
	}
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

	// 파일 노드
	status := fileMissingStyle.Render("✗")
	if node.File.Exists {
		status = fileExistsStyle.Render("✓")
	}
	text := fmt.Sprintf("  %s %s", status, node.Label)
	if selected && focused {
		return treeSelectedStyle.Render(fmt.Sprintf("  %s %s", "›", node.Label))
	}
	return treeItemStyle.Render(text)
}
