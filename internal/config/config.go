package config

import (
	"fmt"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Provider int8

const (
	Unsupported = iota - 1
	Local
	SSH
	GoogleDrive
)

type Game struct {
	SavePath        string     `toml:"save_path"`
	Providers       []Provider `toml:"providers"`
	IncludePatterns []string   `toml:"include_patterns"`
	ExcludePatterns []string   `toml:"exclude_patterns"`
}

func LoadGame() (*Game, error) {
	path := filepath.Join(ConfigPath, GameName+".toml")

	cfg := &Game{}
	_, err := toml.DecodeFile(path, cfg)
	if err != nil {
		return nil, fmt.Errorf("config.Load: failed to load file %s: %w", path, err)
	}

	return cfg, nil
}
