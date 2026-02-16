package model

import "time"

// Scope는 설정 파일의 적용 범위를 나타낸다.
type Scope int

const (
	ScopeManaged Scope = iota // 시스템 관리자 설정
	ScopeUser                 // 사용자 전역 설정
	ScopeProject              // 프로젝트별 설정
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

// FileType은 설정 파일의 형식을 나타낸다.
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

// ConfigCategory는 설정 파일의 기능별 분류를 나타낸다.
type ConfigCategory int

const (
	CategorySettings     ConfigCategory = iota // 동작 설정 (permissions, hooks 등)
	CategoryInstructions                       // 지시사항 (CLAUDE.md 계열)
	CategoryMCP                                // MCP 서버 설정
	CategoryPolicy                             // 관리자 정책
	CategoryCommands                           // 커스텀 슬래시 명령어
	CategorySkills                             // 에이전트 스킬
	CategoryAgents                             // 커스텀 에이전트 정의
	CategoryKeybindings                        // 키바인딩 설정
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
	default:
		return "Unknown"
	}
}

// ConfigFile은 스캔된 개별 설정 파일을 나타낸다.
type ConfigFile struct {
	Path        string         // 절대 경로
	Scope       Scope          // 소속 Scope
	FileType    FileType       // 파일 형식
	Category    ConfigCategory // 기능별 분류
	Exists      bool           // 파일 존재 여부
	IsDir       bool           // 디렉토리 여부 (commands/, skills/)
	Size        int64          // 바이트 크기 (존재 시)
	ModTime     time.Time      // 최종 수정 시간 (존재 시)
	Description string         // 사용자에게 보여줄 설명
	Children    []ConfigFile   // 디렉토리인 경우 하위 파일들
}

// ScanResult는 전체 스캔 결과를 나타낸다.
type ScanResult struct {
	Managed []ConfigFile // Managed Scope 파일들
	User    []ConfigFile // User Scope 파일들
	Project []ConfigFile // Project Scope 파일들
	RootDir string       // 감지된 프로젝트 루트 (없으면 빈 문자열)
}

// All은 모든 Scope의 파일을 하나의 슬라이스로 반환한다.
func (r *ScanResult) All() []ConfigFile {
	all := make([]ConfigFile, 0, len(r.Managed)+len(r.User)+len(r.Project))
	all = append(all, r.Managed...)
	all = append(all, r.User...)
	all = append(all, r.Project...)
	return all
}
