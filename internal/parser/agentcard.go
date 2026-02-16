package parser

import (
	"os"
	"strings"
)

// AgentMeta holds metadata extracted from an agent markdown file.
type AgentMeta struct {
	Name  string // Frontmatter name, or the name part of the first # heading.
	Role  string // The role part after the dash in the first # heading.
	Desc  string // Frontmatter description, or first paragraph (max 150 chars).
	Model string // Frontmatter model (e.g. opus, sonnet).
	Color string // Frontmatter color.
}

// SkillMeta holds metadata extracted from a skill SKILL.md file.
type SkillMeta struct {
	Name     string // Frontmatter name.
	Desc     string // Frontmatter description (max 150 chars).
	Category string // Frontmatter category.
	Tags     string // Frontmatter tags (comma-separated string).
}

const maxDescLen = 150

// ParseAgentMeta parses an agent markdown file and returns its metadata.
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

	// If name is missing from frontmatter, extract from the first # heading
	if meta.Name == "" {
		name, role := parseHeading(body)
		meta.Name = name
		meta.Role = role
	}

	// If description is missing from frontmatter, extract from the first paragraph
	if meta.Desc == "" {
		meta.Desc = truncate(firstParagraph(body), maxDescLen)
	}

	if meta.Name == "" {
		return nil
	}
	return meta
}

// ParseSkillMeta parses a skill SKILL.md file and returns its metadata.
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

	// If name is missing from frontmatter, extract from the first # heading
	if meta.Name == "" {
		meta.Name, _ = parseHeading(body)
	}

	// If description is missing from frontmatter, extract from the first paragraph
	if meta.Desc == "" {
		meta.Desc = truncate(firstParagraph(body), maxDescLen)
	}

	// If tags are missing, try parsing a tag list from the body
	if meta.Tags == "" {
		meta.Tags = parseTagList(body)
	}

	if meta.Name == "" {
		return nil
	}
	return meta
}

// parseFrontmatter parses key: value pairs between --- delimiters and returns the remaining body.
func parseFrontmatter(content string) (map[string]string, string) {
	fm := make(map[string]string)

	trimmed := strings.TrimSpace(content)
	if !strings.HasPrefix(trimmed, "---") {
		return fm, content
	}

	// Find the second --- after the first one
	rest := trimmed[3:]
	idx := strings.Index(rest, "\n---")
	if idx < 0 {
		return fm, content
	}

	block := rest[:idx]
	body := rest[idx+4:] // After "\n---"

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
		// Strip surrounding quotes
		val = strings.Trim(val, "\"'")
		if key != "" && val != "" {
			fm[key] = val
		}
	}

	return fm, body
}

// parseHeading finds the first # heading in body and parses the "name -- role" pattern.
func parseHeading(body string) (name, role string) {
	for _, line := range strings.Split(body, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") {
			heading := strings.TrimPrefix(line, "# ")
			// "name -- role" or "name - role" pattern
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

// firstParagraph returns the first non-empty paragraph from body.
// Lines starting with # (headings) are skipped.
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

// parseTagList finds "- tag" style lists in body and returns them as a comma-separated string.
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

// cleanDesc sanitizes a description string by removing escape characters and extracting the first sentence.
func cleanDesc(desc string) string {
	// Remove escaped newlines
	desc = strings.ReplaceAll(desc, "\\n", " ")
	// Collapse consecutive spaces
	for strings.Contains(desc, "  ") {
		desc = strings.ReplaceAll(desc, "  ", " ")
	}
	desc = strings.TrimSpace(desc)

	// Extract the first sentence (split at period)
	if idx := strings.Index(desc, ". "); idx > 0 && idx < maxDescLen {
		desc = desc[:idx+1]
	}

	return truncate(desc, maxDescLen)
}

// truncate trims a string to at most maxLen runes.
func truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "…"
}
