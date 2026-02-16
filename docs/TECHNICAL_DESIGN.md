# ccfg — Technical Design

## Core Type Definitions

### Scope

An enum representing the scope of a config file.

```go
type Scope int

const (
    ScopeManaged Scope = iota  // System administrator settings
    ScopeUser                   // User global settings
    ScopeProject                // Project-specific settings
)
```

### FileType

The format of a config file.

```go
type FileType int

const (
    FileTypeJSON     FileType = iota  // .json
    FileTypeJSONC                      // .json with comments
    FileTypeMarkdown                   // .md
)
```

### ConfigFile

A single scanned config file.

```go
type ConfigFile struct {
    Path        string    // Absolute path
    Scope       Scope     // Owning Scope
    FileType    FileType  // File format
    Exists      bool      // Whether the file exists
    Size        int64     // Size in bytes (if exists)
    ModTime     time.Time // Last modified time (if exists)
    Description string    // User-facing description
}
```

### ScanResult

The complete scan result.

```go
type ScanResult struct {
    Managed  []ConfigFile  // Managed Scope files
    User     []ConfigFile  // User Scope files
    Project  []ConfigFile  // Project Scope files
    RootDir  string        // Detected project root
}
```

## Component Interfaces

### Scanner

Responsible for discovering config files.

```go
// Scanner discovers Claude Code config files.
type Scanner interface {
    Scan() (*ScanResult, error)
}
```

Implementation: `scanner.DefaultScanner` -- queries OS-specific paths and collects file metadata.

### Parser

Responsible for parsing file contents.

```go
// Parser parses config file contents.
type Parser interface {
    Parse(file ConfigFile) (any, error)
}
```

- `JSONParser`: Parses JSON/JSONC files into `map[string]any`
- `MarkdownParser`: Reads Markdown files as strings

### Merger

Merges JSON settings from multiple Scopes.

```go
// Merger merges settings from multiple Scopes according to priority.
type Merger interface {
    Merge(result *ScanResult) (*MergedConfig, error)
}
```

Merge priority: Project > User > Managed (higher priority overrides lower).

## TUI Architecture

Based on the Elm Architecture (TEA) -- the foundational pattern of Bubbletea.

```
Model -> Update(Msg) -> (Model, Cmd) -> View() -> string
```

### Layout

```
+---------------------------------------------+
|  ccfg v0.1.0 -- Claude Code Config Viewer   |  <- Header
+--------------+------------------------------+
|              |                              |
|  [Tree View] |  [Content Preview]           |  <- Main (30/70)
|  30% width   |  70% width                   |
|              |                              |
|  > Managed   |  {                           |
|    settings  |    "permissions": { ... }    |
|  v User      |    "hooks": { ... }          |
|    settings  |  }                           |
|    CLAUDE.md |                              |
|  > Project   |                              |
|              |                              |
+--------------+------------------------------+
|  [Status Bar]  Tab: switch pane  q: quit    |  <- Footer
+---------------------------------------------+
```

### Key Bindings

| Key | Action |
|---|---|
| `j/k` or `Up/Down` | Move between tree items |
| `Enter` | Expand/collapse node or select file |
| `Tab` | Switch between left/right panels |
| `h/l` | Switch between left/right panels (vim style) |
| `/` | Enter search mode |
| `m` | Toggle merged view |
| `q` / `Ctrl+C` | Quit |

### TUI Model Structure

```go
type Model struct {
    scan       *ScanResult
    tree       TreeModel      // Left tree state
    preview    PreviewModel   // Right preview state
    focus      Pane           // Currently focused panel
    width      int            // Terminal width
    height     int            // Terminal height
    searchMode bool           // Whether search mode is active
    searchText string         // Search query
}
```

## Data Flow

```
main()
  +- flag.Parse()              // Process CLI arguments
  +- scanner.New()             // Create Scanner
  +- scanner.Scan()            // Discover config files
  |    +- ManagedPaths()       // Query Managed paths
  |    +- UserPaths()          // Query User paths
  |    +- ProjectPaths(root)   // Query Project paths
  |    +- stat each file       // Collect file metadata
  +- tui.NewModel(scanResult)  // Initialize TUI model
  +- tea.NewProgram(model)     // Start Bubbletea
       +- Run()
```

## Error Handling Strategy

- File not found: Marked as `ConfigFile.Exists = false`, not treated as an error
- Parse failure: Error message displayed in preview, falls back to raw text
- Permission denied: File marked as "access denied"
- Project root not detected: Project Scope shown as empty
