package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
		validate func(*testing.T, *Config)
	}{
		{
			name: "valid configuration",
			setup: func() {
				os.Setenv("OPENAI_API_KEY", "test-key-123")
				os.Setenv("PORT", "3000")
				os.Setenv("OPENAI_MODEL", "gpt-4o-mini")
			},
			wantErr: false,
			validate: func(t *testing.T, cfg *Config) {
				if cfg.OpenAIAPIKey != "test-key-123" {
					t.Errorf("OpenAIAPIKey = %v, want test-key-123", cfg.OpenAIAPIKey)
				}
				if cfg.Port != "3000" {
					t.Errorf("Port = %v, want 3000", cfg.Port)
				}
				if cfg.OpenAIModel != "gpt-4o-mini" {
					t.Errorf("OpenAIModel = %v, want gpt-4o-mini", cfg.OpenAIModel)
				}
			},
		},
		{
			name: "missing API key",
			setup: func() {
				os.Unsetenv("OPENAI_API_KEY")
			},
			wantErr: true,
		},
		{
			name: "default values",
			setup: func() {
				os.Setenv("OPENAI_API_KEY", "test-key")
				os.Unsetenv("PORT")
				os.Unsetenv("OPENAI_MODEL")
			},
			wantErr: false,
			validate: func(t *testing.T, cfg *Config) {
				if cfg.Port != "8080" {
					t.Errorf("Port = %v, want 8080 (default)", cfg.Port)
				}
				if cfg.OpenAIModel != "gpt-4o" {
					t.Errorf("OpenAIModel = %v, want gpt-4o (default)", cfg.OpenAIModel)
				}
				if cfg.EnableStreaming != true {
					t.Errorf("EnableStreaming = %v, want true (default)", cfg.EnableStreaming)
				}
			},
		},
		{
			name: "boolean environment variables",
			setup: func() {
				os.Setenv("OPENAI_API_KEY", "test-key")
				os.Setenv("ENABLE_STREAMING", "false")
				os.Setenv("ENABLE_VECTOR_STORE", "true")
			},
			wantErr: false,
			validate: func(t *testing.T, cfg *Config) {
				if cfg.EnableStreaming != false {
					t.Errorf("EnableStreaming = %v, want false", cfg.EnableStreaming)
				}
				if cfg.EnableVectorStore != true {
					t.Errorf("EnableVectorStore = %v, want true", cfg.EnableVectorStore)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up environment before each test
			os.Unsetenv("OPENAI_API_KEY")
			os.Unsetenv("PORT")
			os.Unsetenv("OPENAI_MODEL")
			os.Unsetenv("ENABLE_STREAMING")
			os.Unsetenv("ENABLE_VECTOR_STORE")

			tt.setup()

			cfg, err := Load()
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.validate != nil {
				tt.validate(t, cfg)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: &Config{
				OpenAIAPIKey: "test-key",
			},
			wantErr: false,
		},
		{
			name: "missing API key",
			cfg: &Config{
				OpenAIAPIKey: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
