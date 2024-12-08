package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	// Create a temporary test directory
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		setupEnv    func()
		setupFile   func(string)
		wantErr     bool
		expectedURL string
		expectedAddr string
	}{
		{
			name: "successful config load from file",
			setupFile: func(dir string) {
				content := []byte(`DB_URL=postgresql://test:test@localhost:5432/test_db
API_ADDR=:8080`)
				os.WriteFile(dir+"/app.env", content, 0644)
			},
			wantErr:      false,
			expectedURL:  "postgresql://test:test@localhost:5432/test_db",
			expectedAddr: ":8080",
		},
		{
			name: "successful config load from env",
			setupEnv: func() {
				os.Setenv("DB_URL", "postgresql://env:env@localhost:5432/env_db")
				os.Setenv("API_ADDR", ":9090")
			},
			setupFile: func(dir string) {},
			wantErr:   false,
			expectedURL: "postgresql://env:env@localhost:5432/env_db",
			expectedAddr: ":9090",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			if tt.setupEnv != nil {
				tt.setupEnv()
				defer func() {
					os.Unsetenv("DB_URL")
					os.Unsetenv("API_ADDR")
				}()
			}
			tt.setupFile(tmpDir)

			// Test
			cfg, err := newConfig(tmpDir)

			// Assertions
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedURL, cfg.DBURL)
			assert.Equal(t, tt.expectedAddr, cfg.APIAddr)
		})
	}
}
