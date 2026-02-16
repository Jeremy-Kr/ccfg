package model

import "time"

// Scope represents the scope level of a config file.
type Scope int

const (
	ScopeManaged Scope = iota // Managed (admin) settings
	ScopeUser                 // User-level global settings
	ScopeProject              // Project-level settings
)

func (s Scope) String() string {
	switch s {
	case ScopeManaged:
		return "Managed"
	case ScopeUser:
		return "User"
	case ScopeProject:
		return "Project"
	default:
		return "Unknown"
	}
}

// FileType represents the format of a config file.
type FileType int

const (
	FileTypeJSON     FileType = iota // .json
	FileTypeJSONC                    // .json with comments
	FileTypeMarkdown                 // .md
)

func (f FileType) String() string {
	switch f {
	case FileTypeJSON:
		return "JSON"
	case FileTypeJSONC:
		return "JSONC"
	case FileTypeMarkdown:
		return "Markdown"
	default:
		return "Unknown"
	}
}

// ConfigCategory represents the functional category of a config file.
type ConfigCategory int

const (
	CategorySettings     ConfigCategory = iota // Behavior settings (permissions, hooks, etc.)
	CategoryInstructions                       // Instructions (CLAUDE.md family)
	CategoryMCP                                // MCP server settings
	CategoryPolicy                             // Admin policy
	CategoryCommands                           // Custom slash commands
	CategorySkills                             // Agent skills
	CategoryAgents                             // Custom agent definitions
	CategoryKeybindings                        // Keybinding settings
	CategoryHooks                              // Hooks settings
)

func (c ConfigCategory) String() string {
	switch c {
	case CategorySettings:
		return "Settings"
	case CategoryInstructions:
		return "Instructions"
	case CategoryMCP:
		return "MCP"
	case CategoryPolicy:
		return "Policy"
	case CategoryCommands:
		return "Commands"
	case CategorySkills:
		return "Skills"
	case CategoryAgents:
		return "Agents"
	case CategoryKeybindings:
		return "Keybindings"
	case CategoryHooks:
		return "Hooks"
	default:
		return "Unknown"
	}
}

// ConfigFile represents a single scanned config file.
type ConfigFile struct {
	Path        string         // Absolute path
	Scope       Scope          // Owning scope
	FileType    FileType       // File format
	Category    ConfigCategory // Functional category
	Exists      bool           // Whether the file exists
	IsDir       bool           // Whether it is a directory (commands/, skills/)
	IsVirtual   bool           // Whether it is a virtual node (section inside JSON)
	Size        int64          // Size in bytes (when exists)
	ModTime     time.Time      // Last modification time (when exists)
	Description string         // Human-readable description
	Children    []ConfigFile   // Child files when this is a directory
}

// ScanResult represents the complete scan result.
type ScanResult struct {
	Managed []ConfigFile // Managed scope files
	User    []ConfigFile // User scope files
	Project []ConfigFile // Project scope files
	RootDir string       // Detected project root (empty string if none)
}

// All returns all config files from every scope as a single slice.
func (r *ScanResult) All() []ConfigFile {
	all := make([]ConfigFile, 0, len(r.Managed)+len(r.User)+len(r.Project))
	all = append(all, r.Managed...)
	all = append(all, r.User...)
	all = append(all, r.Project...)
	return all
}
