package provider

import (
	"fmt"
	"strings"
	"testing"
)

func TestSimpleParser(t *testing.T) {
	input := `size=10  

[section]
key="value"
`

	t.Log("Tokens:")
	err := DebugTokens(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}

	t.Log("\nParsing result:")
	ini, err := ParseINI(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}

	// Print the parsed structure in a readable format
	t.Logf("Properties: %d", len(ini.Properties))
	for _, p := range ini.Properties {
		t.Logf("  %s = %v", p.Key, formatValue(p.Value))
	}
	t.Logf("Sections: %d", len(ini.Sections))
	for _, s := range ini.Sections {
		t.Logf("  [%s]", s.Identifier)
		for _, p := range s.Properties {
			t.Logf("    %s = %v", p.Key, formatValue(p.Value))
		}
	}
}

// formatValue returns a string representation of a Value
func formatValue(v *Value) string {
	if v.String != nil {
		return *v.String
	}
	if v.Int != nil {
		return fmt.Sprintf("%d", *v.Int)
	}
	if v.Float != nil {
		return fmt.Sprintf("%f", *v.Float)
	}
	return "<nil>"
}
