package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jeremy-kr/ccfg/internal/scanner"
	"github.com/jeremy-kr/ccfg/internal/tui"
)

func main() {
	s := scanner.New("")
	result, err := s.Scan()
	if err != nil {
		fmt.Fprintf(os.Stderr, "스캔 실패: %v\n", err)
		os.Exit(1)
	}

	m := tui.NewModel(result)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "TUI 실행 실패: %v\n", err)
		os.Exit(1)
	}
}
