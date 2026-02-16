package parser

import "encoding/json"

// HookEntry represents an individual event within the hooks section of settings.json.
type HookEntry struct {
	Event    string   // Event name (e.g. "SessionStart").
	Count    int      // Number of registered commands.
	Commands []string // List of command strings.
}

// MCPServerEntry represents an individual MCP server from settings.json or .mcp.json.
type MCPServerEntry struct {
	Name    string // Server name.
	Type    string // Transport type (e.g. "stdio", "sse").
	Command string // Execution command.
}

// ParseSettingsHooks parses the hooks key from raw settings.json (JSONC) content.
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
			// May be a single object instead of an array
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

// ParseMCPServers parses the mcpServers key from raw JSON/JSONC content.
// Works with both settings.json and .mcp.json.
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
