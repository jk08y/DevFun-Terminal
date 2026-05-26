package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config holds all gsh settings.
type Config struct {
	Theme     string        `toml:"theme"`
	ShowGit   bool          `toml:"show_git"`
	ShowUser  bool          `toml:"show_user"`
	ShowHost  bool          `toml:"show_host"`
	History   HistoryConfig `toml:"history"`
	AliasFile string        `toml:"alias_file"`
	RCFile    string        `toml:"rc_file"`
}

// HistoryConfig controls command history behaviour.
type HistoryConfig struct {
	MaxSize int    `toml:"max_size"`
	File    string `toml:"file"`
}

// Default returns a Config with sensible out-of-the-box values.
func Default() *Config {
	home, _ := os.UserHomeDir()
	base := filepath.Join(home, ".config", "gsh")
	return &Config{
		Theme:    "dracula",
		ShowGit:  true,
		ShowUser: true,
		ShowHost: false,
		History: HistoryConfig{
			MaxSize: 10000,
			File:    filepath.Join(base, "history"),
		},
		AliasFile: filepath.Join(base, "aliases"),
		RCFile:    filepath.Join(base, "gshrc"),
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
		_ = writeDefaults(configFile)
		return cfg, nil
	}

	if _, err := toml.DecodeFile(configFile, cfg); err != nil {
		return cfg, fmt.Errorf("parsing config: %w", err)
	}

	// Expand ~ in path fields
	cfg.History.File = expandHome(cfg.History.File, home)
	cfg.AliasFile = expandHome(cfg.AliasFile, home)
	cfg.RCFile = expandHome(cfg.RCFile, home)

	return cfg, nil
}

func expandHome(path, home string) string {
	if len(path) >= 2 && path[:2] == "~/" {
		return filepath.Join(home, path[2:])
	}
	return path
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

alias_file = "~/.config/gsh/aliases"
rc_file    = "~/.config/gsh/gshrc"
`)
	return err
}
