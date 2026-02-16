package main

import (
	"fmt"
	"os"

	"github.com/jeremy-kr/ccfg/internal/model"
	"github.com/jeremy-kr/ccfg/internal/scanner"
)

const version = "0.1.0"

func main() {
	fmt.Printf("ccfg v%s — Claude Code Config Viewer\n\n", version)

	s := scanner.New("")
	result, err := s.Scan()
	if err != nil {
		fmt.Fprintf(os.Stderr, "스캔 실패: %v\n", err)
		os.Exit(1)
	}

	printScope("Managed", result.Managed)
	printScope("User", result.User)

	if result.RootDir != "" {
		fmt.Printf("Project root: %s\n", result.RootDir)
		printScope("Project", result.Project)
	} else {
		fmt.Println("Project: (프로젝트 루트 미감지)")
	}
}

func printScope(name string, files []model.ConfigFile) {
	fmt.Printf("── %s ──\n", name)
	for _, f := range files {
		status := "✗"
		detail := ""
		if f.Exists {
			status = "✓"
			detail = fmt.Sprintf("  (%d bytes, %s)", f.Size, f.ModTime.Format("2006-01-02 15:04"))
		}
		fmt.Printf("  %s %s — %s%s\n", status, f.Description, f.Path, detail)
	}
	fmt.Println()
}
