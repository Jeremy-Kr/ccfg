package parser

import "encoding/json"

// HookEntry는 settings.json의 hooks 내 개별 이벤트를 나타낸다.
type HookEntry struct {
	Event    string   // 이벤트 이름 (예: "SessionStart")
	Count    int      // 등록된 커맨드 수
	Commands []string // 커맨드 문자열 목록
}

// MCPServerEntry는 settings.json 또는 .mcp.json의 개별 MCP 서버를 나타낸다.
type MCPServerEntry struct {
	Name    string // 서버 이름
	Type    string // 전송 타입 (예: "stdio", "sse")
	Command string // 실행 커맨드
}

// ParseSettingsHooks는 settings.json(JSONC) 원본에서 hooks 키를 파싱한다.
func ParseSettingsHooks(raw string) []HookEntry {
	cleaned := StripJSONC(raw)

	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(cleaned), &obj); err != nil {
		return nil
	}

	hooksRaw, ok := obj["hooks"]
	if !ok {
		return nil
	}

	var hooks map[string]json.RawMessage
	if err := json.Unmarshal(hooksRaw, &hooks); err != nil {
		return nil
	}

	var entries []HookEntry
	for event, cmdRaw := range hooks {
		var cmds []map[string]any
		if err := json.Unmarshal(cmdRaw, &cmds); err != nil {
			// 단일 객체일 수 있음
			var single map[string]any
			if err := json.Unmarshal(cmdRaw, &single); err != nil {
				continue
			}
			cmds = []map[string]any{single}
		}

		var commands []string
		for _, c := range cmds {
			if cmd, ok := c["command"].(string); ok {
				commands = append(commands, cmd)
			}
		}

		entries = append(entries, HookEntry{
			Event:    event,
			Count:    len(cmds),
			Commands: commands,
		})
	}
	return entries
}

// ParseMCPServers는 JSON(JSONC) 원본에서 mcpServers 키를 파싱한다.
// settings.json과 .mcp.json 양쪽에서 사용 가능하다.
func ParseMCPServers(raw string) []MCPServerEntry {
	cleaned := StripJSONC(raw)

	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(cleaned), &obj); err != nil {
		return nil
	}

	serversRaw, ok := obj["mcpServers"]
	if !ok {
		return nil
	}

	var servers map[string]json.RawMessage
	if err := json.Unmarshal(serversRaw, &servers); err != nil {
		return nil
	}

	var entries []MCPServerEntry
	for name, srvRaw := range servers {
		var srv map[string]any
		if err := json.Unmarshal(srvRaw, &srv); err != nil {
			entries = append(entries, MCPServerEntry{Name: name})
			continue
		}

		entry := MCPServerEntry{Name: name}
		if t, ok := srv["type"].(string); ok {
			entry.Type = t
		}
		if c, ok := srv["command"].(string); ok {
			entry.Command = c
		}
		entries = append(entries, entry)
	}
	return entries
}
