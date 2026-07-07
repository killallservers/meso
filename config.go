package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config represents the .meso/config.toml structure
type Config struct {
	Project struct {
		Name        string `toml:"name"`
		Description string `toml:"description"`
		Agent       string `toml:"agent"`
	} `toml:"project"`

	Preferences struct {
		DefaultAgent string `toml:"default_agent"`
	} `toml:"preferences"`

	Metadata struct {
		Version      string `toml:"version"`
		CreatedAt    string `toml:"created_at"`
		LastModified string `toml:"last_modified"`
	} `toml:"metadata"`
}

const configVersion = "1.0"
const configDir = ".meso"
const configFile = "config.toml"

func getConfigPath() string {
	return filepath.Join(configDir, configFile)
}

func loadConfig() (*Config, error) {
	path := getConfigPath()
	if _, err := os.Stat(path); err != nil {
		return nil, nil // Config doesn't exist yet
	}

	var cfg Config
	_, err := toml.DecodeFile(path, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return &cfg, nil
}

func saveConfig(cfg *Config) error {
	// Create .meso directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	path := getConfigPath()

	// Write config file
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer f.Close()

	encoder := toml.NewEncoder(f)
	if err := encoder.Encode(cfg); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

func createProjectConfig(info PromptModel) error {
	cfg := &Config{}

	// Project info
	cfg.Project.Name = info.inputs[0]
	cfg.Project.Description = info.inputs[1]
	cfg.Project.Agent = info.agent

	// Preferences
	cfg.Preferences.DefaultAgent = info.agent

	// Metadata
	cfg.Metadata.Version = configVersion
	cfg.Metadata.CreatedAt = "Scaffolded with Meso"
	cfg.Metadata.LastModified = "Scaffolded with Meso"

	return saveConfig(cfg)
}
