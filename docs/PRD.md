# ccfg — Product Requirements Document

## Product Overview

**ccfg** is a read-only tool that displays Claude Code's scattered config files in a single TUI dashboard.

Claude Code distributes its configuration across 8+ files and directories:
- Managed settings controlled by system administrators
- Global settings in the user's home directory
- Project-specific local settings

It is difficult to determine which settings come from where and what the currently active values are. ccfg solves this problem.

## Target Users

- Developers who use Claude Code on a daily basis
- Users who maintain different Claude Code settings across multiple projects
- Power users who frequently modify MCP servers, hooks, and permission settings

## Scanned Files (3 Scopes)

### Scope 1: Managed (System Administrator Settings)
| File | Description |
|---|---|
| `/Library/Application Support/ClaudeCode/managed_settings.json` | Managed settings (macOS) |
| `/Library/Application Support/ClaudeCode/policies.json` | Policy file (macOS) |

### Scope 2: User (User Global Settings)
| File | Description |
|---|---|
| `~/.claude/settings.json` | User global settings |
| `~/.claude/settings.local.json` | User local settings (gitignored) |
| `~/.claude.json` | Legacy global settings |
| `~/.claude/CLAUDE.md` | User global instructions |
| `~/.mcp.json` | MCP server global settings |

### Scope 3: Project (Project-Specific Settings)
| File | Description |
|---|---|
| `<root>/.claude/settings.json` | Project settings |
| `<root>/.claude/settings.local.json` | Project local settings |
| `<root>/CLAUDE.md` | Project instructions |
| `<root>/.claude/CLAUDE.md` | Project instructions (alternate location) |
| `<root>/.mcp.json` | MCP server project settings |

## Functional Requirements

### FR-1: Config File Discovery
- Automatically discover all config files across the 3 Scopes
- Display file existence, size, and modification time
- Auto-detect project root from CWD (by traversing up to find `.git` directory)

### FR-2: Tree View Navigation
- Display a Scope > File hierarchy tree in the left panel
- Navigate the tree with keyboard (j/k/Enter/Esc)
- Visually distinguish files based on existence status

### FR-3: Content Preview
- Display the selected file's content in the right panel
- JSON/JSONC files: syntax highlighted
- Markdown files: rendered display

### FR-4: Merged View
- Display the final merged result of all Scope JSON settings
- Show the source Scope for each value
- Priority: Project > User > Managed

### FR-5: Search
- Search by config key/value
- Highlight search results

## Non-Functional Requirements

- **Performance:** Under 500ms from launch to TUI display
- **Read-only:** Never modifies config files
- **Single binary:** No external runtime dependencies
- **OS Support:** macOS first, Linux support planned
- **Accessibility:** Full keyboard-only navigation
