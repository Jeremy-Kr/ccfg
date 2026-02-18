<p align="center">
  <h1 align="center">ccfg</h1>
</p>

<p align="center">
  A TUI dashboard for viewing Claude Code configuration files at a glance.
</p>

<p align="center">
  <a href="https://go.dev/"><img src="https://img.shields.io/github/go-mod/go-version/jeremy-kr/ccfg" alt="Go Version"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="License: MIT"></a>
  <a href="https://github.com/jeremy-kr/ccfg/releases"><img src="https://img.shields.io/github/v/release/jeremy-kr/ccfg" alt="Release"></a>
  <a href="https://github.com/jeremy-kr/tap"><img src="https://img.shields.io/badge/homebrew-available-orange" alt="Homebrew"></a>
</p>

## Demo

<p align="center">
  <img src="./docs/demo.gif" alt="ccfg demo" width="800">
</p>

> [!TIP]
> Record your own demo with [VHS](https://github.com/charmbracelet/vhs): `vhs docs/demo.tape`

## Why?

Claude Code stores configuration across **8+ files** in 3 different scopes:

```
~/.claude/settings.json              # user preferences
~/.claude/CLAUDE.md                  # global instructions
~/.mcp.json                          # global MCP servers
project/.claude/settings.json        # project overrides
project/CLAUDE.md                    # project instructions
project/.mcp.json                    # project MCP servers
managed_settings.json                # organization policies
policies.json                        # managed policies
...and more (hooks, keybindings, agents, skills)
```

When settings conflict or something isn't working, you need to check each file manually. **ccfg** gives you a single, navigable view of everything.

## Features

- **Unified view** — See managed, user, and project config files side by side
- **Tree navigation** — Browse config files organized by scope in a collapsible tree
- **Syntax highlighting** — JSON/JSONC highlighted with Chroma, Markdown rendered with Glamour
- **Merged view** — View the final merged configuration with source annotations
- **Search** — Find settings by key or value across all files
- **Auto-refresh** — Detects file changes via fsnotify and updates in real time
- **Extended scanning** — Custom commands, agent skills, hooks, MCP servers, and keybindings
- **Usage rankings** — Gamified tool/agent/skill statistics with SSS~F tier grades and time period filters (24h/7d/30d/All)
- **Character cards** — Custom agents and skills displayed as game-style cards
- **Read-only** — Never modifies any configuration file

## Installation

### Homebrew (recommended)

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

| Key                | Action                                        |
| ------------------ | --------------------------------------------- |
| `j/k` or `Up/Down` | Move between tree items                       |
| `Enter`            | Expand/collapse node or select file           |
| `Tab` or `h/l`     | Switch between left/right panels              |
| `/`                | Enter search mode                             |
| `Esc`              | Exit search / back                            |
| `m`                | Toggle merged view                            |
| `1/2/3`            | Switch ranking tabs (tools / agents / skills) |
| `s`                | Toggle ranking scope (all / project)          |
| `p`                | Cycle ranking period (All / 30d / 7d / 24h)   |
| `q` / `Ctrl+C`     | Quit                                          |

### Flags

```bash
ccfg --version    # Print version
```

## Scanned Files

ccfg discovers config files from three scopes:

| Scope       | Examples                                                        |
| ----------- | --------------------------------------------------------------- |
| **Managed** | `managed_settings.json`, `policies.json`                        |
| **User**    | `~/.claude/settings.json`, `~/.claude/CLAUDE.md`, `~/.mcp.json` |
| **Project** | `.claude/settings.json`, `CLAUDE.md`, `.mcp.json`               |

See [docs/PRD.md](docs/PRD.md) for the complete list.

## Tech Stack

- **Language:** [Go](https://go.dev/) 1.26
- **TUI Framework:** [Bubbletea](https://github.com/charmbracelet/bubbletea) v1.3.x (Elm Architecture)
- **Styling:** [Lipgloss](https://github.com/charmbracelet/lipgloss)
- **Syntax Highlighting:** [Glamour](https://github.com/charmbracelet/glamour) + [Chroma](https://github.com/alecthomas/chroma)
- **File Watching:** [fsnotify](https://github.com/fsnotify/fsnotify)

## Support

If you find ccfg useful, consider giving it a **star** on GitHub — it helps others discover the project!

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup, coding conventions, and the PR process.

## License

[MIT](LICENSE)
