package parser

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/chroma/v2/quick"
)

// FormatJSON pretty-prints JSON/JSONC content and applies syntax highlighting.
// Returns the original text on failure.
func FormatJSON(raw string) string {
	// Convert JSONC to JSON (strip comments and trailing commas)
	cleaned := StripJSONC(raw)

	// Pretty-print
	var parsed any
	if err := json.Unmarshal([]byte(cleaned), &parsed); err != nil {
		// Parse failed; attempt syntax highlighting on the original
		return highlightJSON(raw)
	}

	pretty, err := json.MarshalIndent(parsed, "", "  ")
	if err != nil {
		return highlightJSON(raw)
	}

	return highlightJSON(string(pretty))
}

// highlightJSON applies JSON syntax highlighting using Chroma.
func highlightJSON(src string) string {
	var buf bytes.Buffer
	err := quick.Highlight(&buf, src, "json", "terminal256", "monokai")
	if err != nil {
		return src
	}
	return buf.String()
}

// StripJSONC strips comments and trailing commas from JSONC.
func StripJSONC(s string) string {
	var result strings.Builder
	runes := []rune(s)
	i := 0
	inString := false

	for i < len(runes) {
		ch := runes[i]

		// Handle characters inside a string literal
		if inString {
			result.WriteRune(ch)
			if ch == '\\' && i+1 < len(runes) {
				i++
				result.WriteRune(runes[i])
			} else if ch == '"' {
				inString = false
			}
			i++
			continue
		}

		// Start of a string literal
		if ch == '"' {
			inString = true
			result.WriteRune(ch)
			i++
			continue
		}

		// Single-line comment
		if ch == '/' && i+1 < len(runes) && runes[i+1] == '/' {
			for i < len(runes) && runes[i] != '\n' {
				i++
			}
			continue
		}

		// Block comment
		if ch == '/' && i+1 < len(runes) && runes[i+1] == '*' {
			i += 2
			for i+1 < len(runes) && !(runes[i] == '*' && runes[i+1] == '/') {
				i++
			}
			i += 2 // Skip past */
			continue
		}

		// Trailing comma: ,] or ,}
		if ch == ',' {
			// Skip whitespace/newlines after comma and check for ] or }
			j := i + 1
			for j < len(runes) && (runes[j] == ' ' || runes[j] == '\t' || runes[j] == '\n' || runes[j] == '\r') {
				j++
			}
			if j < len(runes) && (runes[j] == ']' || runes[j] == '}') {
				// Omit trailing comma
				i++
				continue
			}
		}

		result.WriteRune(ch)
		i++
	}

	return result.String()
}
