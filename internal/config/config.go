package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config holds all gsh settings.
type Config struct {
	Theme    string        `toml:"theme"`
	ShowGit  bool          `toml:"show_git"`
	ShowUser bool          `toml:"show_user"`
	ShowHost bool          `toml:"show_host"`
	History  HistoryConfig `toml:"history"`
}

// HistoryConfig controls command history behaviour.
type HistoryConfig struct {
	MaxSize int    `toml:"max_size"`
	File    string `toml:"file"`
}

// Default returns a Config with sensible out-of-the-box values.
func Default() *Config {
	home, _ := os.UserHomeDir()
	return &Config{
		Theme:    "dracula",
		ShowGit:  true,
		ShowUser: true,
		ShowHost: false,
		History: HistoryConfig{
			MaxSize: 10000,
			File:    filepath.Join(home, ".config", "gsh", "history"),
		},
	}
}

// Load reads ~/.config/gsh/config.toml, creating it with defaults if absent.
func Load() (*Config, error) {
	cfg := Default()

	home, err := os.UserHomeDir()
	if err != nil {
		return cfg, nil
	}

	configDir := filepath.Join(home, ".config", "gsh")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return cfg, fmt.Errorf("creating config dir: %w", err)
	}

	configFile := filepath.Join(configDir, "config.toml")

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		_ = writeDefaults(configFile) // best-effort
		return cfg, nil
	}

	if _, err := toml.DecodeFile(configFile, cfg); err != nil {
		return cfg, fmt.Errorf("parsing config: %w", err)
	}

	// Expand ~ in history file path
	if len(cfg.History.File) > 1 && cfg.History.File[:2] == "~/" {
		cfg.History.File = filepath.Join(home, cfg.History.File[2:])
	}

	return cfg, nil
}

func writeDefaults(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = fmt.Fprintf(f, `# gsh configuration
# See https://github.com/jk08y/gsh for documentation

theme     = "dracula"   # dracula | nord | catppuccin | onedark | tokyo-night
show_git  = true        # show git branch in prompt
show_user = true        # show username in prompt
show_host = false       # show hostname in prompt

[history]
max_size = 10000
file     = "~/.config/gsh/history"
`)
	return err
}
