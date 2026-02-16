# ccfg — Claude Code Config Viewer

A TUI dashboard built with Go + Bubbletea. View all of Claude Code's scattered config files in a single screen.

## Tech Stack

- **Language:** Go 1.26
- **TUI:** Bubbletea v1.3.x + Lipgloss + Bubbles
- **Syntax Highlighting:** Glamour + Chroma
- **CLI:** Standard `flag` package (no external framework)

## Build / Test / Lint

```bash
# Build
go build -o ccfg ./cmd/ccfg

# Run all tests
go test ./...

# Test a specific package
go test ./internal/scanner/...

# Lint (requires golangci-lint)
golangci-lint run

# Run
./ccfg
```

## Directory Structure

```
cmd/ccfg/          Entry point
internal/
  scanner/         Config file discovery (paths, scan logic)
  parser/          Config file parsing (JSON, JSONC, Markdown)
  model/           Shared type definitions
  tui/             Bubbletea model, view, components
docs/              PRD, technical design, roadmap
```

## Coding Conventions

- Go standard style (`gofmt`)
- Encapsulation via `internal/` packages
- Errors wrapped with `fmt.Errorf("context: %w", err)` pattern
- Interfaces defined at the call site
- Test files use `*_test.go` in the same package
- Naming follows Go conventions (camelCase for unexported, PascalCase for exported)

## Key Constraints

- Read-only tool: never modifies config files
- macOS first (Linux support added in Phase 5)
- Minimal dependencies: only add libraries when necessary
