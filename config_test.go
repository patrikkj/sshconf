package sshconf

import (
	"strings"
	"testing"
)

func TestSSHConfigRender(t *testing.T) {
	input := `# Global settings
Host example
    HostName example.com
    User myuser
    # Port comment
    Port 22

Host *
    Port 22
    User default`

	config := ParseConfig(input)
	rendered := config.Render()

	// The rendered output should match the input exactly
	if rendered != input {
		t.Errorf("Render() produced different output than input.\nGot:\n%s\nWant:\n%s", rendered, input)
	}
}

func TestSSHConfigWrite(t *testing.T) {
	input := `Host example
    HostName example.com
    User myuser`

	config := ParseConfig(input)
	var buf strings.Builder
	err := config.Write(&buf)
	if err != nil {
		t.Fatalf("Write() failed: %v", err)
	}

	if buf.String() != input {
		t.Errorf("Write() produced different output than input.\nGot:\n%s\nWant:\n%s", buf.String(), input)
	}
}
