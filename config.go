package sshconf

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// SSHConfig represents a parsed SSH config file
type SSHConfig struct {
	lines []Line
}

// ParseConfigRaw parses an SSH config file without organizing the lines hierarchically
func ParseConfigRaw(content string) *SSHConfig {
	return &SSHConfig{
		lines: parseLines(content),
	}
}

// ParseConfigFile parses an SSH config file from the given path
func ParseConfigFile(path string) (*SSHConfig, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseConfigRaw(string(content)), nil
}

// ParseConfig parses an SSH config file and organizes it hierarchically
func ParseConfig(content string) *SSHConfig {
	return &SSHConfig{
		lines: OrganizeConfig(parseLines(content)),
	}
}

// renderLine renders a single Line struct back to a string, including its children
func renderLine(line Line) string {
	var builder strings.Builder

	// Render the current line
	builder.WriteString(line.Indent)
	builder.WriteString(line.Key)
	builder.WriteString(line.Sep)
	builder.WriteString(line.Value)
	builder.WriteString(line.TrailIndent)
	builder.WriteString(line.Comment)

	// Render children if any exist
	for _, child := range line.Children {
		builder.WriteString("\n")
		builder.WriteString(renderLine(child))
	}

	return builder.String()
}

// Render returns the config as a string
func (c *SSHConfig) Render() string {
	var builder strings.Builder
	for i, line := range c.lines {
		if i > 0 {
			builder.WriteString("\n")
		}
		builder.WriteString(renderLine(line))
	}
	return builder.String()
}

// Write writes the config to the given writer
func (c *SSHConfig) Write(w io.Writer) error {
	_, err := w.Write([]byte(c.Render()))
	return err
}

// WriteFile writes the config to the given file path
func (c *SSHConfig) WriteFile(path string) error {
	return os.WriteFile(path, []byte(c.Render()), 0644)
}

// Patch replaces a directive in the config with new content, or appends it if not found.
// The find parameter should be a single line (e.g., "Host example").
// The replacement parameter can be multiple lines of config.
func (c *SSHConfig) Patch(find, replacement string) error {
	// Parse the find line to match against
	findLine := ParseLine(find)
	if findLine.Key == "" {
		return fmt.Errorf("invalid find directive: must contain a key")
	}

	// Parse the replacement content
	replacementConfig := ParseConfig(replacement)

	// Try to find and replace the directive
	found := false
	for i, line := range c.lines {
		if line.Key == findLine.Key && line.Value == findLine.Value {
			// Replace the existing line and its children with the replacement
			c.lines = append(c.lines[:i], append(replacementConfig.lines, c.lines[i+1:]...)...)
			found = true
			break
		}
	}

	// If directive wasn't found, append it to the end
	if !found {
		c.lines = append(c.lines, replacementConfig.lines...)
	}

	return nil
}

// Delete removes a directive and its children from the config.
// The find parameter should be a single line (e.g., "Host example").
func (c *SSHConfig) Delete(find string) error {
	return c.Patch(find, "")
}
