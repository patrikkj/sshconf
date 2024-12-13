package sshconf

import (
	"bufio"
	"regexp"
	"strings"
)

type LineNoChildren struct {
	Indent      string // The indentation of the line
	Key         string // The directive/keyword (e.g., "Host", "HostName", etc.)
	Sep         string // The separator between the key and the value (e.g., " ", "=")
	Value       string // The values associated with the directive
	TrailIndent string // The indentation of the trailing comment
	Comment     string // Any comment on the line, for lines with only a comment everything except indent and comment is empty
}

type Line struct {
	LineNoChildren
	Children []LineNoChildren // The children of the line (for Host and Match directives)
}

// ParseLine parses a single line of an SSH config file into a Line struct
func ParseLine(line string) Line {
	// Regex pattern to match SSH config line components
	parts := []string{
		`^(\s*)?`,                              // Group 1: Leading indentation (optional whitespace at start of line)
		`([^\s=#]+)?`,                          // Group 2: Key/directive (captures text until it hits whitespace, =, or #)
		`(\s*=?\s*)?`,                          // Group 3: Separator (whitespace before/after an optional = character)
		`((?:[^"#\s][^\s#]*|"[^"]*"|\s+?)*?)?`, // Group 4: Value (quoted or unquoted strings)
		`(\s*)?`,                               // Group 5: Trailing whitespace before any comment
		`(#.*)?$`,                              // Group 6: Comment (# followed by any text until end of line)
	}
	pattern := strings.Join(parts, "")
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(line)

	if matches == nil {
		// Return empty Line if no match (shouldn't happen with this pattern)
		return Line{}
	}

	return Line{
		LineNoChildren: LineNoChildren{
			Indent:      matches[1],
			Key:         matches[2],
			Sep:         matches[3],
			Value:       matches[4],
			TrailIndent: matches[5],
			Comment:     matches[6],
		},
	}
}

// parseLines parses an entire SSH config file content into a slice of Lines
func parseLines(content string) []Line {
	var lines []Line
	scanner := bufio.NewScanner(strings.NewReader(content))

	for scanner.Scan() {
		lines = append(lines, ParseLine(scanner.Text()))
	}

	return lines
}

// OrganizeConfig takes a flat slice of Lines and returns a hierarchical structure
// where Host and Match blocks have their parameters as children
func OrganizeConfig(lines []Line) []Line {
	var result []Line
	var currentParent *Line

	for _, line := range lines {
		// If line has no key (empty line or comment), add it to current parent or result
		if line.Key == "" {
			if currentParent != nil {
				currentParent.Children = append(currentParent.Children, line.LineNoChildren)
			} else {
				result = append(result, line)
			}
			continue
		}

		// Check if this is a Host or Match directive
		if strings.EqualFold(line.Key, "Host") || strings.EqualFold(line.Key, "Match") {
			result = append(result, line)
			currentParent = &result[len(result)-1]
		} else {
			// This is a parameter line
			if currentParent != nil {
				currentParent.Children = append(currentParent.Children, line.LineNoChildren)
			} else {
				result = append(result, line)
			}
		}
	}

	return cleanEmptyLines(result)
}

// cleanEmptyLines processes a slice of Lines and moves trailing empty lines
// from children to the parent level
func cleanEmptyLines(lines []Line) []Line {
	var result []Line

	for _, line := range lines {
		newLine := line
		var trailingEmptyLines []Line

		if len(line.Children) > 0 {
			for i := len(newLine.Children) - 1; i >= 0; i-- {
				child := newLine.Children[i]
				if child.Key == "" && child.Indent == "" {
					trailingEmptyLines = append(trailingEmptyLines, Line{LineNoChildren: child})
					newLine.Children = newLine.Children[:i]
				} else {
					break
				}
			}

			result = append(result, newLine)

			for i := len(trailingEmptyLines) - 1; i >= 0; i-- {
				result = append(result, trailingEmptyLines[i])
			}
		} else {
			result = append(result, line)
		}
	}

	return result
}
