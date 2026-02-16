package parser

import (
	"os"
	"strings"
)

// AgentMeta — 에이전트 md에서 추출한 메타데이터
type AgentMeta struct {
	Name  string // frontmatter name 또는 첫 # 헤딩의 이름 부분
	Role  string // 첫 # 헤딩의 "— 역할" 부분
	Desc  string // frontmatter description 또는 첫 단락 (최대 150자)
	Model string // frontmatter model (opus/sonnet 등)
	Color string // frontmatter color
}

// SkillMeta — 스킬 SKILL.md에서 추출한 메타데이터
type SkillMeta struct {
	Name     string // frontmatter name
	Desc     string // frontmatter description (최대 150자)
	Category string // frontmatter category
	Tags     string // frontmatter tags (쉼표 구분 문자열)
}

const maxDescLen = 150

// ParseAgentMeta는 에이전트 md 파일을 파싱하여 메타데이터를 반환한다.
func ParseAgentMeta(path string) *AgentMeta {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	content := string(data)
	if strings.TrimSpace(content) == "" {
		return nil
	}

	fm, body := parseFrontmatter(content)

	meta := &AgentMeta{
		Name:  fm["name"],
		Model: fm["model"],
		Color: fm["color"],
	}
	if desc := fm["description"]; desc != "" {
		meta.Desc = cleanDesc(desc)
	}

	// frontmatter에 name이 없으면 첫 # 헤딩에서 추출
	if meta.Name == "" {
		name, role := parseHeading(body)
		meta.Name = name
		meta.Role = role
	}

	// frontmatter에 description이 없으면 첫 단락에서 추출
	if meta.Desc == "" {
		meta.Desc = truncate(firstParagraph(body), maxDescLen)
	}

	if meta.Name == "" {
		return nil
	}
	return meta
}

// ParseSkillMeta는 스킬 SKILL.md 파일을 파싱하여 메타데이터를 반환한다.
func ParseSkillMeta(path string) *SkillMeta {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	content := string(data)
	if strings.TrimSpace(content) == "" {
		return nil
	}

	fm, body := parseFrontmatter(content)

	meta := &SkillMeta{
		Name:     fm["name"],
		Category: fm["category"],
		Tags:     fm["tags"],
	}
	if desc := fm["description"]; desc != "" {
		meta.Desc = cleanDesc(desc)
	}

	// frontmatter에 name이 없으면 첫 # 헤딩에서 추출
	if meta.Name == "" {
		meta.Name, _ = parseHeading(body)
	}

	// frontmatter에 description이 없으면 첫 단락에서 추출
	if meta.Desc == "" {
		meta.Desc = truncate(firstParagraph(body), maxDescLen)
	}

	// tags가 없으면 body에서 태그 리스트 파싱 시도
	if meta.Tags == "" {
		meta.Tags = parseTagList(body)
	}

	if meta.Name == "" {
		return nil
	}
	return meta
}

// parseFrontmatter는 --- 구분자 사이의 key: value를 파싱하고 나머지 body를 반환한다.
func parseFrontmatter(content string) (map[string]string, string) {
	fm := make(map[string]string)

	trimmed := strings.TrimSpace(content)
	if !strings.HasPrefix(trimmed, "---") {
		return fm, content
	}

	// 첫 번째 --- 이후에서 두 번째 --- 찾기
	rest := trimmed[3:]
	idx := strings.Index(rest, "\n---")
	if idx < 0 {
		return fm, content
	}

	block := rest[:idx]
	body := rest[idx+4:] // "\n---" 이후

	for _, line := range strings.Split(block, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		colonIdx := strings.Index(line, ":")
		if colonIdx < 0 {
			continue
		}
		key := strings.TrimSpace(line[:colonIdx])
		val := strings.TrimSpace(line[colonIdx+1:])
		// 따옴표 제거
		val = strings.Trim(val, "\"'")
		if key != "" && val != "" {
			fm[key] = val
		}
	}

	return fm, body
}

// parseHeading은 body에서 첫 # 헤딩을 찾아 "name — role" 패턴을 파싱한다.
func parseHeading(body string) (name, role string) {
	for _, line := range strings.Split(body, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") {
			heading := strings.TrimPrefix(line, "# ")
			// "name — role" 또는 "name - role" 패턴
			for _, sep := range []string{" — ", " - ", " – "} {
				if idx := strings.Index(heading, sep); idx > 0 {
					return strings.TrimSpace(heading[:idx]), strings.TrimSpace(heading[idx+len(sep):])
				}
			}
			return strings.TrimSpace(heading), ""
		}
	}
	return "", ""
}

// firstParagraph은 body에서 첫 비어있지 않은 단락을 반환한다.
// #으로 시작하는 헤딩은 건너뛴다.
func firstParagraph(body string) string {
	var para []string
	inPara := false

	for _, line := range strings.Split(body, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			if inPara {
				break
			}
			continue
		}
		if strings.HasPrefix(trimmed, "#") {
			if inPara {
				break
			}
			continue
		}
		inPara = true
		para = append(para, trimmed)
	}
	return strings.Join(para, " ")
}

// parseTagList는 body에서 "- tag" 형태의 리스트를 찾아 쉼표 구분 문자열로 반환한다.
func parseTagList(body string) string {
	var tags []string
	inTagSection := false

	for _, line := range strings.Split(body, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(strings.ToLower(trimmed), "tag") && strings.HasPrefix(trimmed, "#") {
			inTagSection = true
			continue
		}
		if inTagSection {
			if strings.HasPrefix(trimmed, "- ") {
				tag := strings.TrimSpace(strings.TrimPrefix(trimmed, "- "))
				if tag != "" {
					tags = append(tags, tag)
				}
			} else if trimmed == "" {
				continue
			} else if strings.HasPrefix(trimmed, "#") {
				break
			}
		}
	}
	return strings.Join(tags, ", ")
}

// cleanDesc는 description 문자열에서 이스케이프 문자를 정리하고 첫 문장만 추출한다.
func cleanDesc(desc string) string {
	// 이스케이프된 줄바꿈 제거
	desc = strings.ReplaceAll(desc, "\\n", " ")
	// 연속 공백 정리
	for strings.Contains(desc, "  ") {
		desc = strings.ReplaceAll(desc, "  ", " ")
	}
	desc = strings.TrimSpace(desc)

	// 첫 문장 추출 (마침표 기준)
	if idx := strings.Index(desc, ". "); idx > 0 && idx < maxDescLen {
		desc = desc[:idx+1]
	}

	return truncate(desc, maxDescLen)
}

// truncate는 문자열을 maxLen 이하로 잘라낸다.
func truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "…"
}
