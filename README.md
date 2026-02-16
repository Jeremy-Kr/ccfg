# ccfg

A TUI dashboard for viewing Claude Code configuration files at a glance.

Claude Code stores its configuration across 8+ files spread across managed, user, and project scopes. **ccfg** collects them all into a single, navigable terminal interface so you can instantly see what is configured and where each setting comes from.

## Features

- **Unified view** -- See managed, user, and project config files side by side
- **Tree navigation** -- Browse config files organized by scope in a collapsible tree
- **Syntax highlighting** -- JSON/JSONC files highlighted with Chroma, Markdown rendered with Glamour
- **Merged view** -- View the final merged configuration with source annotations
- **Search** -- Find settings by key or value across all files
- **Auto-refresh** -- Detects file changes via fsnotify and updates in real time
- **Extended scanning** -- Custom commands, agent skills, hooks, MCP servers, and keybindings
- **Read-only** -- Never modifies any configuration file

## Installation

### Homebrew

```bash
brew install jeremy-kr/tap/ccfg
```

### Go Install

```bash
go install github.com/jeremy-kr/ccfg/cmd/ccfg@latest
```

### From Source

```bash
git clone https://github.com/jeremy-kr/ccfg.git
cd ccfg
go build -o ccfg ./cmd/ccfg
```

## Usage

Run `ccfg` from within any project directory:

```bash
ccfg
```

### Key Bindings

| Key | Action |
|---|---|
| `j/k` or `Up/Down` | Move between tree items |
| `Enter` | Expand/collapse node or select file |
| `Tab` | Switch between left/right panels |
| `h/l` | Switch panels (vim style) |
| `/` | Enter search mode |
| `m` | Toggle merged view |
| `q` / `Ctrl+C` | Quit |

### Flags

```bash
ccfg --version    # Print version
```

## Screenshots

<!-- TODO: Add screenshots -->

## Tech Stack

- **Language:** [Go](https://go.dev/) 1.26
- **TUI Framework:** [Bubbletea](https://github.com/charmbracelet/bubbletea) v1.3.x
- **Styling:** [Lipgloss](https://github.com/charmbracelet/lipgloss)
- **Syntax Highlighting:** [Glamour](https://github.com/charmbracelet/glamour) + [Chroma](https://github.com/alecthomas/chroma)
- **File Watching:** [fsnotify](https://github.com/fsnotify/fsnotify)

## Scanned Files

ccfg discovers config files from three scopes:

| Scope | Examples |
|---|---|
| **Managed** | `/Library/Application Support/ClaudeCode/managed_settings.json`, `policies.json` |
| **User** | `~/.claude/settings.json`, `~/.claude/CLAUDE.md`, `~/.mcp.json` |
| **Project** | `.claude/settings.json`, `CLAUDE.md`, `.mcp.json` |

See [docs/PRD.md](docs/PRD.md) for the complete list.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup, coding conventions, and the PR process.

## License

[MIT](LICENSE)
