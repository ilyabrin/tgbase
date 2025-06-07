package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Test setup
	testConfigContent := `
database:
  type: postgres
  postgres:
    dsn: "postgresql://user:pass@localhost:5432/db"
  sqlite:
    path: "./test.db"
redis:
  enabled: true
  addr: "localhost:6379"
  password: "testpass"
  db: 0
telegram:
  token: "test-token"
`
	tempFile := "test_config.yaml"
	err := os.WriteFile(tempFile, []byte(testConfigContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}
	defer os.Remove(tempFile)

	tests := []struct {
		name           string
		configPath     string
		envVars        map[string]string
		expectedConfig *Config
		expectError    bool
	}{
		{
			name:       "Valid config file",
			configPath: tempFile,
			envVars:    map[string]string{},
			expectedConfig: &Config{
				Database: struct {
					Type     string `yaml:"type"`
					Postgres struct {
						DSN string `yaml:"dsn"`
					} `yaml:"postgres"`
					SQLite struct {
						Path string `yaml:"path"`
					} `yaml:"sqlite"`
				}{
					Type: "postgres",
					Postgres: struct {
						DSN string `yaml:"dsn"`
					}{
						DSN: "postgresql://user:pass@localhost:5432/db",
					},
					SQLite: struct {
						Path string `yaml:"path"`
					}{
						Path: "./test.db",
					},
				},
				Redis: struct {
					Enabled  bool   `yaml:"enabled"`
					Addr     string `yaml:"addr"`
					Password string `yaml:"password"`
					DB       int    `yaml:"db"`
				}{
					Enabled:  true,
					Addr:     "localhost:6379",
					Password: "testpass",
					DB:       0,
				},
				Telegram: struct {
					Token string `yaml:"token"`
				}{
					Token: "test-token",
				},
			},
			expectError: false,
		},
		{
			name:       "Environment variables override",
			configPath: tempFile,
			envVars: map[string]string{
				"TELEGRAM_TOKEN": "env-token",
				"POSTGRES_DSN":   "postgresql://env:pass@localhost:5432/envdb",
				"REDIS_ADDR":     "env-redis:6379",
			},
			expectedConfig: &Config{
				Database: struct {
					Type     string `yaml:"type"`
					Postgres struct {
						DSN string `yaml:"dsn"`
					} `yaml:"postgres"`
					SQLite struct {
						Path string `yaml:"path"`
					} `yaml:"sqlite"`
				}{
					Type: "postgres",
					Postgres: struct {
						DSN string `yaml:"dsn"`
					}{
						DSN: "postgresql://env:pass@localhost:5432/envdb",
					},
					SQLite: struct {
						Path string `yaml:"path"`
					}{
						Path: "./test.db",
					},
				},
				Redis: struct {
					Enabled  bool   `yaml:"enabled"`
					Addr     string `yaml:"addr"`
					Password string `yaml:"password"`
					DB       int    `yaml:"db"`
				}{
					Enabled:  true,
					Addr:     "env-redis:6379",
					Password: "testpass",
					DB:       0,
				},
				Telegram: struct {
					Token string `yaml:"token"`
				}{
					Token: "env-token",
				},
			},
			expectError: false,
		},
		{
			name:           "Non-existent config file",
			configPath:     "non_existent_config.yaml",
			envVars:        map[string]string{},
			expectedConfig: nil,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
				defer os.Unsetenv(key)
			}

			// Load configuration
			cfg, err := LoadConfig(tt.configPath)

			// Check error expectation
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
				return
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// If we expect an error, no need to check the config
			if tt.expectError {
				return
			}

			// Compare configurations
			if cfg.Database.Type != tt.expectedConfig.Database.Type {
				t.Errorf("Database type mismatch: got %v, want %v", cfg.Database.Type, tt.expectedConfig.Database.Type)
			}
			if cfg.Database.Postgres.DSN != tt.expectedConfig.Database.Postgres.DSN {
				t.Errorf("Postgres DSN mismatch: got %v, want %v", cfg.Database.Postgres.DSN, tt.expectedConfig.Database.Postgres.DSN)
			}
			if cfg.Redis.Addr != tt.expectedConfig.Redis.Addr {
				t.Errorf("Redis addr mismatch: got %v, want %v", cfg.Redis.Addr, tt.expectedConfig.Redis.Addr)
			}
			if cfg.Telegram.Token != tt.expectedConfig.Telegram.Token {
				t.Errorf("Telegram token mismatch: got %v, want %v", cfg.Telegram.Token, tt.expectedConfig.Telegram.Token)
			}
		})
	}
}
