package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	AnthropicAPIKey string `json:"anthropic_api_key"`
	DefaultModel    string `json:"default_model"`
}

func DefaultConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".claude-setup", "config.json")
}

func ClaudeDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".claude")
}

func Load() (*Config, error) {
	path := DefaultConfigPath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{DefaultModel: "claude-sonnet-4-6"}, nil
		}
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if cfg.DefaultModel == "" {
		cfg.DefaultModel = "claude-sonnet-4-6"
	}
	return &cfg, nil
}

func (c *Config) Save() error {
	path := DefaultConfigPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
