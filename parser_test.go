package sshconf

import (
	"reflect"
	"testing"
)

func TestParseLine(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Line
	}{
		{
			name:  "Simple key-value",
			input: "Host example.com",
			expected: Line{
				LineNoChildren: LineNoChildren{
					Indent: "",
					Key:    "Host",
					Sep:    " ",
					Value:  "example.com",
				},
			},
		},
		{
			name:  "Simple key-value with multiple values",
			input: "Host example.com example.com.staging",
			expected: Line{
				LineNoChildren: LineNoChildren{
					Indent: "",
					Key:    "Host",
					Sep:    " ",
					Value:  "example.com example.com.staging",
				},
			},
		},
		{
			name:  "With comment",
			input: "Host example.com example.com.staging  # My server",
			expected: Line{
				LineNoChildren: LineNoChildren{
					Indent:      "",
					Key:         "Host",
					Sep:         " ",
					Value:       "example.com example.com.staging",
					TrailIndent: "  ",
					Comment:     "# My server",
				},
			},
		},
		{
			name:  "Indented with equals",
			input: "    IdentityFile = ~/.ssh/id_rsa",
			expected: Line{
				LineNoChildren: LineNoChildren{
					Indent: "    ",
					Key:    "IdentityFile",
					Sep:    " = ",
					Value:  "~/.ssh/id_rsa",
				},
			},
		},
		{
			name:  "Comment only",
			input: "# Just a comment",
			expected: Line{
				LineNoChildren: LineNoChildren{
					Comment: "# Just a comment",
				},
			},
		},
		{
			name:     "Empty line",
			input:    "",
			expected: Line{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseLine(tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ParseLine() = %+v, want %+v", got, tt.expected)
			}
		})
	}
}

func TestParseConfig(t *testing.T) {
	input := `# SSH Config
Host example.com
    HostName example.com
    User myuser
    IdentityFile = ~/.ssh/id_rsa  # My key

Host *.staging
    User staging-user`

	expected := []Line{
		{
			LineNoChildren: LineNoChildren{
				Indent:  "",
				Comment: "# SSH Config",
			},
		},
		{
			LineNoChildren: LineNoChildren{
				Indent: "",
				Key:    "Host",
				Sep:    " ",
				Value:  "example.com",
			},
		},
		{
			LineNoChildren: LineNoChildren{
				Indent: "    ",
				Key:    "HostName",
				Sep:    " ",
				Value:  "example.com",
			},
		},
		{
			LineNoChildren: LineNoChildren{
				Indent: "    ",
				Key:    "User",
				Sep:    " ",
				Value:  "myuser",
			},
		},
		{
			LineNoChildren: LineNoChildren{
				Indent:      "    ",
				Key:         "IdentityFile",
				Sep:         " = ",
				Value:       "~/.ssh/id_rsa",
				TrailIndent: "  ",
				Comment:     "# My key",
			},
		},
		{
			LineNoChildren: LineNoChildren{},
		},
		{
			LineNoChildren: LineNoChildren{
				Indent: "",
				Key:    "Host",
				Sep:    " ",
				Value:  "*.staging",
			},
		},
		{
			LineNoChildren: LineNoChildren{
				Indent: "    ",
				Key:    "User",
				Sep:    " ",
				Value:  "staging-user",
			},
		},
	}

	got := parseLines(input)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("parseLines() = %+v, want %+v", got, expected)
	}
}

func TestOrganizeConfig(t *testing.T) {
	input := `# Global settings
Host example
    HostName example.com
    User myuser
    # Port comment
    Port 22

Host *
    Port 22
    User default`

	lines := parseLines(input)
	organized := OrganizeConfig(lines)

	expected := []Line{
		{
			LineNoChildren: LineNoChildren{
				Indent:  "",
				Comment: "# Global settings",
			},
		},
		{
			LineNoChildren: LineNoChildren{
				Indent: "",
				Key:    "Host",
				Sep:    " ",
				Value:  "example",
			},
			Children: []LineNoChildren{
				{
					Indent: "    ",
					Key:    "HostName",
					Sep:    " ",
					Value:  "example.com",
				},
				{
					Indent: "    ",
					Key:    "User",
					Sep:    " ",
					Value:  "myuser",
				},
				{
					Indent:  "    ",
					Comment: "# Port comment",
				},
				{
					Indent: "    ",
					Key:    "Port",
					Sep:    " ",
					Value:  "22",
				},
			},
		},
		{
			LineNoChildren: LineNoChildren{},
		},
		{
			LineNoChildren: LineNoChildren{
				Indent: "",
				Key:    "Host",
				Sep:    " ",
				Value:  "*",
			},
			Children: []LineNoChildren{
				{
					Indent: "    ",
					Key:    "Port",
					Sep:    " ",
					Value:  "22",
				},
				{
					Indent: "    ",
					Key:    "User",
					Sep:    " ",
					Value:  "default",
				},
			},
		},
	}

	if !reflect.DeepEqual(organized, expected) {
		t.Errorf("OrganizeConfig result does not match expected structure.\nGot: %+v\nWant: %+v", organized, expected)
	}
}

func TestOrganizeConfigEmptyLines(t *testing.T) {
	input := `Host example
    HostName example.com
    User myuser

Host other
    Port 22

`
	lines := parseLines(input)
	organized := OrganizeConfig(lines)

	expected := []Line{
		{
			LineNoChildren: LineNoChildren{
				Indent: "",
				Key:    "Host",
				Sep:    " ",
				Value:  "example",
			},
			Children: []LineNoChildren{
				{
					Indent: "    ",
					Key:    "HostName",
					Sep:    " ",
					Value:  "example.com",
				},
				{
					Indent: "    ",
					Key:    "User",
					Sep:    " ",
					Value:  "myuser",
				},
			},
		},
		{
			LineNoChildren: LineNoChildren{},
		},
		{
			LineNoChildren: LineNoChildren{
				Indent: "",
				Key:    "Host",
				Sep:    " ",
				Value:  "other",
			},
			Children: []LineNoChildren{
				{
					Indent: "    ",
					Key:    "Port",
					Sep:    " ",
					Value:  "22",
				},
			},
		},
		{
			LineNoChildren: LineNoChildren{},
		},
	}

	if !reflect.DeepEqual(organized, expected) {
		t.Errorf("OrganizeConfig empty lines handling failed.\nGot: %+v\nWant: %+v", organized, expected)
	}
}
