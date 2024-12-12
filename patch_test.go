package sshconf

import (
	"testing"
)

func TestSSHConfigPatch(t *testing.T) {
	tests := []struct {
		name        string
		initial     string
		find        string
		replacement string
		expected    string
		wantErr     bool
	}{
		{
			name: "Replace existing host",
			initial: `Host example
    User old-user
    Port 22

Host other
    User other-user`,
			find: "Host example",
			replacement: `Host example
    User new-user
    IdentityFile ~/.ssh/id_rsa`,
			expected: `Host example
    User new-user
    IdentityFile ~/.ssh/id_rsa

Host other
    User other-user`,
		},
		{
			name: "Append new host",
			initial: `Host example
    User old-user`,
			find: "Host new",
			replacement: `Host new
    User new-user`,
			expected: `Host example
    User old-user
Host new
    User new-user`,
		},
		{
			name:        "Invalid find directive",
			initial:     `Host example`,
			find:        "  # just a comment",
			replacement: `Host new`,
			wantErr:     true,
		},
		{
			name:        "Empty replacement",
			initial:     `Host example`,
			find:        "Host example",
			replacement: "",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := ParseConfig(tt.initial)
			err := config.Patch(tt.find, tt.replacement)

			if (err != nil) != tt.wantErr {
				t.Errorf("Patch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				got := config.Render()
				if got != tt.expected {
					t.Errorf("Patch() produced incorrect output.\nGot:\n%s\nWant:\n%s", got, tt.expected)
				}
			}
		})
	}
}

func TestSSHConfigDelete(t *testing.T) {
	tests := []struct {
		name    string
		config  string
		find    string
		want    string
		wantErr bool
	}{
		{
			name: "Delete existing host",
			config: `Host example
    HostName example.com

Host other
    HostName other.com`,
			find: "Host example",
			want: `
Host other
    HostName other.com`,
			wantErr: false,
		},
		{
			name: "Delete non-existent host",
			config: `Host example
    HostName example.com`,
			find: "Host nonexistent",
			want: `Host example
    HostName example.com`,
			wantErr: false,
		},
		{
			name: "Invalid find directive",
			config: `Host example
    HostName example.com`,
			find: "  ",
			want: `Host example
    HostName example.com`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ParseConfig(tt.config)
			err := c.Delete(tt.find)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && c.Render() != tt.want {
				t.Errorf("Delete() got = %v, want %v", c.Render(), tt.want)
			}
		})
	}
}
