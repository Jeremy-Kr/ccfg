package tui

import (
	"strings"
	"testing"
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
