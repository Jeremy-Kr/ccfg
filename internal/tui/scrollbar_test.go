package tui

import (
	"strings"
	"testing"

	"github.com/jeremy-kr/ccfg/internal/model"
)

func TestRenderScrollbar_NoScrollNeeded(t *testing.T) {
	// total <= visible이면 nil 반환
	if got := renderScrollbar(5, 10, 0); got != nil {
		t.Errorf("expected nil, got %v", got)
	}
	if got := renderScrollbar(10, 10, 0); got != nil {
		t.Errorf("expected nil for equal, got %v", got)
	}
}

func TestRenderScrollbar_ZeroVisible(t *testing.T) {
	if got := renderScrollbar(10, 0, 0); got != nil {
		t.Errorf("expected nil for zero visible, got %v", got)
	}
}

func TestRenderScrollbar_Length(t *testing.T) {
	bars := renderScrollbar(100, 20, 0)
	if len(bars) != 20 {
		t.Errorf("expected 20 bars, got %d", len(bars))
	}
}

func TestRenderScrollbar_MinThumbSize(t *testing.T) {
	// total이 매우 크면 thumbSize가 최소 1이어야 함
	bars := renderScrollbar(10000, 5, 0)
	if len(bars) != 5 {
		t.Errorf("expected 5 bars, got %d", len(bars))
	}
	// thumb이 최소 1개 존재해야 함
	thumbCount := countThumb(bars)
	if thumbCount < 1 {
		t.Errorf("expected at least 1 thumb, got %d", thumbCount)
	}
}

func TestRenderScrollbar_ThumbAtTop(t *testing.T) {
	bars := renderScrollbar(40, 10, 0)
	// offset=0이면 thumb이 상단에 위치
	if len(bars) == 0 {
		t.Fatal("expected non-nil bars")
	}
	thumbCount := countThumb(bars)
	if thumbCount < 1 {
		t.Error("expected at least 1 thumb char")
	}
}

func TestRenderScrollbar_ThumbAtBottom(t *testing.T) {
	total, visible := 40, 10
	maxOffset := total - visible
	bars := renderScrollbar(total, visible, maxOffset)
	if len(bars) != visible {
		t.Fatalf("expected %d bars, got %d", visible, len(bars))
	}
	// offset=maxOffset이면 마지막 행이 thumb이어야 함
	thumbCount := countThumb(bars)
	if thumbCount < 1 {
		t.Error("expected at least 1 thumb char at bottom")
	}
}

func TestPreviewView_ScrollbarInOutput(t *testing.T) {
	p := &PreviewModel{height: 12, offset: 0}
	lines := make([]string, 50)
	for i := range lines {
		lines[i] = "test content"
	}
	p.lines = lines
	p.file = &model.ConfigFile{Path: "/test.txt", Exists: true, Size: 100}

	output := p.View(80, false)
	if !strings.Contains(output, "┃") && !strings.Contains(output, "│") {
		t.Error("scrollbar chars (┃ or │) not found in preview output")
	}
}

func TestTreeView_ScrollbarInOutput(t *testing.T) {
	// visibleNodes가 height보다 많으면 스크롤바 표시
	roots := []TreeNode{
		{Label: "Test", Scope: model.ScopeUser, Expanded: true},
	}
	for i := 0; i < 30; i++ {
		roots[0].Children = append(roots[0].Children, TreeNode{
			Label: "item",
			Scope: model.ScopeUser,
			File:  &model.ConfigFile{Path: "/test", Description: "item"},
		})
	}
	tm := TreeModel{roots: roots, height: 10}
	output := tm.View(40, false)
	if !strings.Contains(output, "┃") && !strings.Contains(output, "│") {
		t.Error("scrollbar chars not found in tree output")
	}
}

// countThumb은 ANSI 스타일 적용된 thumb 문자("┃")를 포함하는 행 수를 센다.
func countThumb(bars []string) int {
	count := 0
	for _, b := range bars {
		if strings.Contains(b, "┃") {
			count++
		}
	}
	return count
}
