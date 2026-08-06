package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{
			name: "valid configuration",
			content: `
server:
  port: 8080

loadBalancer:
  strategy: round-robin

upstreams:
  - id: api-1
    url: http://localhost:9001
`,
			wantErr: false,
		},
		{
			name: "invalid yaml",
			content: `
server:
  port: [
`,
			wantErr: true,
		},
		{
			name: "invalid configuration",
			content: `
server:
  port: 0

loadBalancer:
  strategy: round-robin

upstreams:
  - id: api-1
    url: http://localhost:9001
`,
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tempDir := t.TempDir()

			configPath := filepath.Join(tempDir, "config.yaml")

			if err := os.WriteFile(
				configPath,
				[]byte(test.content),
				0644,
			); err != nil {
				t.Fatalf("failed to create config file: %v", err)
			}

			_, err := Load(configPath)

			if test.wantErr && err == nil {
				t.Fatal("expected load to fail")
			}

			if !test.wantErr && err != nil {
				t.Fatalf("expected load to succeed: %v", err)
			}
		})
	}
}
