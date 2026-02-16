package tui

import (
	"fmt"
	"os"
	"strings"

	"github.com/jeremy-kr/ccfg/internal/model"
	"github.com/jeremy-kr/ccfg/internal/parser"
)

// PreviewModel은 우측 미리보기 패널의 상태를 관리한다.
type PreviewModel struct {
	file    *model.ConfigFile // 현재 표시 중인 파일
	content string            // 파일 내용
	lines   []string          // 줄 단위 분할
	offset  int               // 스크롤 오프셋
	height  int               // 표시 가능한 행 수
}

// SetFile은 미리보기에 표시할 파일을 설정한다.
func (p *PreviewModel) SetFile(file *model.ConfigFile) {
	if file == nil {
		p.file = nil
		p.content = ""
		p.lines = nil
		p.offset = 0
		return
	}

	// 이미 같은 파일이면 스킵
	if p.file != nil && p.file.Path == file.Path {
		return
	}

	p.file = file
	p.offset = 0

	if !file.Exists {
		p.content = "(파일이 존재하지 않습니다)"
		p.lines = []string{p.content}
		return
	}

	// 디렉토리인 경우 내용 목록 표시
	if file.IsDir {
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
			icon := "  📄 "
			if entry.IsDir() {
				icon = "  📁 "
			}
			b.WriteString(fmt.Sprintf("%s%s\n", icon, entry.Name()))
		}
	} else {
		for _, child := range file.Children {
			icon := "  📄 "
			if child.IsDir {
				icon = "  📁 "
			}
			detail := ""
			if child.Exists {
				detail = fmt.Sprintf("  (%d bytes)", child.Size)
			}
			b.WriteString(fmt.Sprintf("%s%s%s\n", icon, child.Description, detail))
		}
	}

	return b.String()
}

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

// View는 미리보기를 문자열로 렌더링한다.
func (p *PreviewModel) View(width int, focused bool) string {
	var b strings.Builder

	if p.file == nil {
		b.WriteString("파일을 선택하세요")
	} else {
		// 파일 정보 헤더
		header := p.file.Path
		if p.file.Exists {
			header = fmt.Sprintf("%s  (%d bytes)", p.file.Path, p.file.Size)
		}
		b.WriteString(headerStyle.Render(header))
		b.WriteString("\n")

		// 내용 표시
		end := p.offset + p.height - 1 // 헤더 1줄 제외
		if end > len(p.lines) {
			end = len(p.lines)
		}
		for i := p.offset; i < end; i++ {
			b.WriteString(p.lines[i])
			if i < end-1 {
				b.WriteString("\n")
			}
		}
	}

	// 패널 높이 고정 — 내용이 길어도 패널이 커지지 않도록
	style := panelStyle.Width(width).Height(p.height)
	if focused {
		style = panelFocusedStyle.Width(width).Height(p.height)
	}

	return style.Render(b.String())
}
