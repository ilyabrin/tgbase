package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database struct {
		Type     string `yaml:"type"`
		Postgres struct {
			DSN string `yaml:"dsn"`
		} `yaml:"postgres"`
		SQLite struct {
			Path string `yaml:"path"`
		} `yaml:"sqlite"`
	} `yaml:"database"`
	Redis struct {
		Enabled  bool   `yaml:"enabled"`
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	} `yaml:"redis"`
	Telegram struct {
		Token    string  `yaml:"token"`
		AdminIDs []int64 `yaml:"admin_ids"`
	} `yaml:"telegram"`
}

// MustLoad loads the config and panics on error. Useful in main().
func MustLoad(path string) *Config {
	cfg, err := LoadConfig(path)
	if err != nil {
		panic("failed to load config: " + err.Error())
	}
	return cfg
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Перезапись конфигурации из переменных окружения
	if token := os.Getenv("TELEGRAM_TOKEN"); token != "" {
		cfg.Telegram.Token = token
	}
	if dsn := os.Getenv("POSTGRES_DSN"); dsn != "" {
		cfg.Database.Postgres.DSN = dsn
	}
	if addr := os.Getenv("REDIS_ADDR"); addr != "" {
		cfg.Redis.Addr = addr
	}

	return &cfg, nil
}
