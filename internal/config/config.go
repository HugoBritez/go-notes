package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	NotesRoot string `json:"notes_root"`
	Editor    string `json:"editor"`
}

func GetConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "go-notes", "config.json")
}

func Load() (*Config, error) {
	path := GetConfigPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = json.Unmarshal(data, &cfg)
	return &cfg, err
}

func (c *Config) Save() error {
	path := GetConfigPath()

	os.MkdirAll(filepath.Dir(path), 0o755)

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}
