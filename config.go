package sshconf

import (
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

// renderLine renders a single Line struct back to a string
func renderLine(line Line) string {
	var builder strings.Builder

	builder.WriteString(line.Indent)
	if line.Key != "" {
		builder.WriteString(line.Key)
		builder.WriteString(line.Sep)
		builder.WriteString(line.Value)
	}
	if line.TrailIndent != "" {
		builder.WriteString(line.TrailIndent)
	}
	if line.Comment != "" {
		builder.WriteString(line.Comment)
	}

	return builder.String()
}

// renderLines renders a slice of Lines, including their children
func renderLines(lines []Line) []string {
	var result []string

	for _, line := range lines {
		result = append(result, renderLine(line))
		if len(line.Children) > 0 {
			result = append(result, renderLines(line.Children)...)
		}
	}

	return result
}

// Render returns the config as a string
func (c *SSHConfig) Render() string {
	return strings.Join(renderLines(c.lines), "\n")
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
