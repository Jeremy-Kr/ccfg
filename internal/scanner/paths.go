package scanner

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/jeremy-kr/ccfg/internal/model"
)

// FileEntry는 스캔 대상 파일 하나를 정의한다.
type FileEntry struct {
	RelPath     string               // 기준 디렉토리로부터의 상대 경로
	Description string               // 사용자에게 보여줄 설명
	Category    model.ConfigCategory // 기능별 분류
	IsDir       bool                 // 디렉토리 스캔 여부
}

// GetUserHomeDir는 현재 사용자의 홈 디렉토리를 반환한다.
func GetUserHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return home
}

// ManagedPaths는 시스템 관리자가 관리하는 설정 파일 경로를 반환한다.
func ManagedPaths() (string, []FileEntry) {
	var base string
	switch runtime.GOOS {
	case "darwin":
		base = "/Library/Application Support/ClaudeCode"
	case "linux":
		base = "/etc/claude-code"
	default:
		return "", nil
	}

	return base, []FileEntry{
		{RelPath: "managed_settings.json", Description: "관리형 설정", Category: model.CategorySettings},
		{RelPath: "policies.json", Description: "정책 파일", Category: model.CategoryPolicy},
	}
}

// UserPaths는 사용자 전역 설정 파일 경로를 반환한다.
func UserPaths() (string, []FileEntry) {
	home := GetUserHomeDir()
	if home == "" {
		return "", nil
	}

	return home, []FileEntry{
		{RelPath: filepath.Join(".claude", "settings.json"), Description: "사용자 전역 설정", Category: model.CategorySettings},
		{RelPath: filepath.Join(".claude", "settings.local.json"), Description: "사용자 로컬 설정", Category: model.CategorySettings},
		{RelPath: ".claude.json", Description: "레거시 전역 설정", Category: model.CategorySettings},
		{RelPath: filepath.Join(".claude", "CLAUDE.md"), Description: "사용자 전역 지시사항", Category: model.CategoryInstructions},
		{RelPath: ".mcp.json", Description: "MCP 서버 전역 설정", Category: model.CategoryMCP},
		{RelPath: filepath.Join(".claude", "commands"), Description: "커스텀 명령어", Category: model.CategoryCommands, IsDir: true},
		{RelPath: filepath.Join(".claude", "skills"), Description: "에이전트 스킬", Category: model.CategorySkills, IsDir: true},
		{RelPath: filepath.Join(".claude", "keybindings.json"), Description: "키바인딩 설정", Category: model.CategoryKeybindings},
	}
}

// ProjectPaths는 프로젝트별 설정 파일 경로를 반환한다.
func ProjectPaths(root string) (string, []FileEntry) {
	if root == "" {
		return "", nil
	}

	return root, []FileEntry{
		{RelPath: filepath.Join(".claude", "settings.json"), Description: "프로젝트 설정", Category: model.CategorySettings},
		{RelPath: filepath.Join(".claude", "settings.local.json"), Description: "프로젝트 로컬 설정", Category: model.CategorySettings},
		{RelPath: "CLAUDE.md", Description: "프로젝트 지시사항", Category: model.CategoryInstructions},
		{RelPath: filepath.Join(".claude", "CLAUDE.md"), Description: "프로젝트 지시사항 (대체 위치)", Category: model.CategoryInstructions},
		{RelPath: ".mcp.json", Description: "MCP 서버 프로젝트 설정", Category: model.CategoryMCP},
		{RelPath: filepath.Join(".claude", "commands"), Description: "프로젝트 명령어", Category: model.CategoryCommands, IsDir: true},
		{RelPath: filepath.Join(".claude", "skills"), Description: "프로젝트 스킬", Category: model.CategorySkills, IsDir: true},
	}
}
