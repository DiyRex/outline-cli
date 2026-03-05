package config

import (
	"testing"
)

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name:    "valid config",
			config:  Config{URL: "https://example.com", APIKey: "ol_api_test123"},
			wantErr: false,
		},
		{
			name:    "missing URL",
			config:  Config{APIKey: "ol_api_test123"},
			wantErr: true,
		},
		{
			name:    "missing API key",
			config:  Config{URL: "https://example.com"},
			wantErr: true,
		},
		{
			name:    "both missing",
			config:  Config{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigDir(t *testing.T) {
	dir := ConfigDir()
	if dir == "" {
		t.Error("ConfigDir() should not be empty")
	}
}

func TestConfigPath(t *testing.T) {
	path := ConfigPath()
	if path == "" {
		t.Error("ConfigPath() should not be empty")
	}
}
