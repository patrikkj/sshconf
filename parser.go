package provider

import (
	"io"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

// Config represents a complete SSH config file
type Config struct {
	Lines []*Line `@@*`
}

// Line represents either a comment, directive, or empty line
type Line struct {
	LeadingSpace *string    `(@Whitespace)?`
	Comment      *string    `( @Comment @Newline`
	TrailingNL   *string    `)?`
	Directive    *Directive `| @@`
	Empty        *string    `| @EmptyLine`
}

// Directive represents a configuration directive and its children
type Directive struct {
	Key           string       `@Ident`
	Space         *string      `(@Whitespace)?`
	Value         *string      `(@Value)?`
	TrailingSpace *string      `(@Whitespace)?`
	Comment       *string      `(@Comment)?`
	TrailingNL    string       `@Newline`
	Children      []*Directive `( @@ )*`
}

var (
	sshLexer = lexer.MustSimple([]lexer.SimpleRule{
		{"Comment", `#[^\n]*`},
		{"EmptyLine", `\n\s*\n`},
		{"Whitespace", `[\t ]+`},
		{"Newline", `\n`},
		{"Ident", `[A-Za-z][A-Za-z0-9_-]*`},
		{"Value", `[^\s#\n][^\n#]*`},
	})

	parser = participle.MustBuild[Config](
		participle.Lexer(sshLexer),
		participle.UseLookahead(2),
	)
)

// ParseConfig parses an SSH config from a reader
func ParseConfig(r io.Reader) (*Config, error) {
	return parser.Parse("", r)
}

// String converts the Config back to its string representation
func (c *Config) String() string {
	var result string
	for i, line := range c.Lines {
		if i > 0 && line.LeadingSpace == nil {
			result += "\n"
		}

		if line.LeadingSpace != nil {
			result += *line.LeadingSpace
		}

		if line.Comment != nil {
			result += *line.Comment
			if line.TrailingNL != nil {
				result += *line.TrailingNL
			}
		} else if line.Empty != nil {
			result += *line.Empty
		} else if line.Directive != nil {
			result += formatDirective(line.Directive, 0)
		}
	}
	return result
}

// formatDirective formats a directive with proper indentation
func formatDirective(d *Directive, indent int) string {
	indentStr := strings.Repeat("    ", indent)
	result := indentStr + d.Key

	if d.Space != nil {
		result += *d.Space
	}
	if d.Value != nil {
		result += *d.Value
	}
	if d.TrailingSpace != nil {
		result += *d.TrailingSpace
	}
	if d.Comment != nil {
		result += *d.Comment
	}
	if d.TrailingNL != "" {
		result += d.TrailingNL
	}

	for _, child := range d.Children {
		result += formatDirective(child, indent+1)
	}

	return result
}
