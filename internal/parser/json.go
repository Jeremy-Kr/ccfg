package parser

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/chroma/v2/quick"
)

// FormatJSON은 JSON/JSONC 내용을 pretty-print하고 구문 강조를 적용한다.
// 실패하면 원본 텍스트를 반환한다.
func FormatJSON(raw string) string {
	// JSONC → JSON 변환 (주석, trailing comma 제거)
	cleaned := stripJSONC(raw)

	// pretty-print
	var parsed any
	if err := json.Unmarshal([]byte(cleaned), &parsed); err != nil {
		// 파싱 실패 → 원본에 구문 강조만 시도
		return highlightJSON(raw)
	}

	pretty, err := json.MarshalIndent(parsed, "", "  ")
	if err != nil {
		return highlightJSON(raw)
	}

	return highlightJSON(string(pretty))
}

// highlightJSON은 Chroma로 JSON 구문 강조를 적용한다.
func highlightJSON(src string) string {
	var buf bytes.Buffer
	err := quick.Highlight(&buf, src, "json", "terminal256", "monokai")
	if err != nil {
		return src
	}
	return buf.String()
}

// stripJSONC는 JSONC의 주석과 trailing comma를 제거한다.
func stripJSONC(s string) string {
	var result strings.Builder
	runes := []rune(s)
	i := 0
	inString := false

	for i < len(runes) {
		ch := runes[i]

		// 문자열 내부 처리
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

		// 문자열 시작
		if ch == '"' {
			inString = true
			result.WriteRune(ch)
			i++
			continue
		}

		// 한 줄 주석
		if ch == '/' && i+1 < len(runes) && runes[i+1] == '/' {
			for i < len(runes) && runes[i] != '\n' {
				i++
			}
			continue
		}

		// 블록 주석
		if ch == '/' && i+1 < len(runes) && runes[i+1] == '*' {
			i += 2
			for i+1 < len(runes) && !(runes[i] == '*' && runes[i+1] == '/') {
				i++
			}
			i += 2 // */ 건너뛰기
			continue
		}

		// trailing comma: ,] 또는 ,}
		if ch == ',' {
			// 쉼표 뒤의 공백/개행을 건너뛰고 ] 또는 }가 오는지 확인
			j := i + 1
			for j < len(runes) && (runes[j] == ' ' || runes[j] == '\t' || runes[j] == '\n' || runes[j] == '\r') {
				j++
			}
			if j < len(runes) && (runes[j] == ']' || runes[j] == '}') {
				// trailing comma 생략
				i++
				continue
			}
		}

		result.WriteRune(ch)
		i++
	}

	return result.String()
}
