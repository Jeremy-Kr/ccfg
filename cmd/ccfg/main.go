package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jeremy-kr/ccfg/internal/scanner"
	"github.com/jeremy-kr/ccfg/internal/tui"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		short := commit
		if len(short) > 7 {
			short = short[:7]
		}
		fmt.Printf("ccfg %s (%s, %s)\n", version, short, date)
		return
	}

	s := scanner.New("")

	start := time.Now()
	result, err := s.Scan()
	scanDuration := time.Since(start)

	if err != nil {
		fmt.Fprintf(os.Stderr, "scan failed: %v\n", err)
		os.Exit(1)
	}

	m := tui.NewModel(result, scanDuration, s)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to start TUI: %v\n", err)
		os.Exit(1)
	}
}
