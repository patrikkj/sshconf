package provider

import (
	"fmt"
	"io"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

// INI represents a basic INI-style configuration file
type INI struct {
	Properties []*Property `@@*`
	Sections   []*Section  `@@*`
}

type Section struct {
	Identifier string      `"[" @Ident "]"`
	Properties []*Property `@@*`
}

type Property struct {
	Key   string `@Ident "="`
	Value *Value `@@`
}

type Value struct {
	String *string  `  @String`
	Float  *float64 `| @Float`
	Int    *int     `| @Int`
}

var (
	iniLexer = lexer.MustSimple([]lexer.SimpleRule{
		{"Whitespace", `\s+`},
		{"String", `"[^"]*"`}, // Only match quoted strings
		{"Float", `\d*\.\d+`},
		{"Int", `\d+`},
		{"Ident", `[a-zA-Z_][a-zA-Z0-9_]*`}, // Match identifiers
		{"Punct", `[=\[\]]`},
	})
	iniParser, _ = participle.Build[INI](
		participle.Lexer(iniLexer),
	)
)

// ParseINI parses an INI-style config from a reader
func ParseINI(r io.Reader) (*INI, error) {
	return iniParser.Parse("", r)
}

// DebugTokens prints all tokens from the lexer without consuming them
func DebugTokens(r io.Reader) error {
	tokens, err := iniLexer.Lex("", r)
	if err != nil {
		return fmt.Errorf("lexing error: %w", err)
	}

	// Create reverse mapping from token type to name
	typeToName := make(map[lexer.TokenType]string)
	for name, tokenType := range iniLexer.Symbols() {
		typeToName[tokenType] = name
	}

	const format = "Token: Type = %-10s Value = %-20q Pos = %v\n"

	for {
		token, err := tokens.Next()
		if err != nil {
			return fmt.Errorf("token iteration error: %w", err)
		}
		if token.Type == lexer.EOF {
			break
		}
		tokenName := typeToName[token.Type]
		fmt.Printf(format, tokenName, token.Value, token.Pos)
	}
	return nil
}
