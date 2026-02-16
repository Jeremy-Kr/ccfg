# ccfg — Roadmap

## Phase 1: Project Scaffolding ✅

- [x] Go environment setup
- [x] Project documentation (PRD, technical design, roadmap)
- [x] go.mod initialization
- [x] Minimal entry point (`cmd/ccfg/main.go`)
- [x] Path constants definition (`internal/scanner/paths.go`)
- [x] Git init + .gitignore

## Phase 2: Scanner Implementation ✅

- [x] `internal/model/types.go` -- Shared type definitions (Scope, FileType, ConfigFile, ScanResult)
- [x] `internal/scanner/scanner.go` -- DefaultScanner implementation
- [x] `internal/scanner/root.go` -- Project root detection (.git traversal)
- [x] `internal/scanner/scanner_test.go` -- Tests (temp directory based)
- [x] Add Bubbletea, Lipgloss, and other dependencies

## Phase 3: Basic TUI ✅

- [x] `internal/tui/model.go` -- Main TUI model (Init, Update, View)
- [x] `internal/tui/tree.go` -- Tree view component
- [x] `internal/tui/preview.go` -- Preview panel (raw text)
- [x] `internal/tui/keys.go` -- Key binding definitions
- [x] `internal/tui/styles.go` -- Lipgloss style definitions
- [x] 30/70 split layout
- [x] Keyboard navigation (j/k/Enter/Tab/q)

## Phase 4: Parser + Syntax Highlighting ✅

- [x] `internal/parser/json.go` -- JSON/JSONC parser (strips comments/trailing commas)
- [x] `internal/parser/markdown.go` -- Markdown renderer (Glamour)
- [x] Syntax highlighting (Chroma monokai theme) applied
- [x] Raw text fallback on parse errors

## Phase 5: Extended Scan Targets ✅

- [x] `~/.claude/commands/`, `.claude/commands/` -- Custom slash commands
- [x] `~/.claude/skills/`, `.claude/skills/` -- Installed agent skills
- [x] `~/.claude/keybindings.json` -- Keybinding settings
- [x] Directory scanning support (list files under commands/, skills/ in addition to individual files)
- [x] ConfigCategory extension (Commands, Skills, Keybindings added)
- [x] Parse `"hooks"` key inside settings.json -> display hook list in tree
- [x] Parse `"mcpServers"` key inside settings.json -> display MCP server list

## Phase 6: Advanced Features ✅

- [x] Merged View -- merge all Scope settings + show source
- [x] Search functionality (enter with `/` key)
- [x] Linux path support
- [x] File change detection (fsnotify) and auto-refresh

## Phase 7: Distribution + Quality ✅

- [x] golangci-lint configuration (`.golangci.yml`)
- [x] goreleaser configuration (`.goreleaser.yml`)
- [x] Makefile (build/test/lint/run/clean)
- [ ] README.md (installation instructions, screenshots)
- [x] Homebrew formula
- [x] GitHub Actions release automation
- [x] MIT license

## Phase 8: Game-Style Theme ✅

- [x] Colorful borders (unique color per Scope)
- [x] Emoji icons (Scope, file status, category)
- [x] Section title styling (labels on top of borders)
- [x] Bottom key hint bar (game HUD style)
- [x] Full color palette redesign (dark background + neon accents)
