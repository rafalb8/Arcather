package game

import (
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/rafalb8/Arcather/defaults"
	"github.com/rafalb8/Arcather/internal/config"
	"github.com/rafalb8/Arcather/internal/rclone"
	"github.com/rafalb8/ln"
)

var environ = strings.NewReplacer(
	"$HOME", ln.Must(os.UserHomeDir()),
	"$XDG_CONFIG_HOME", ln.Must(os.UserConfigDir()),
)

type Config struct {
	SavePath        string          `toml:"save_path"`
	Remote          []rclone.Remote `toml:"remotes"`
	IncludePatterns []string        `toml:"include_patterns"`
	ExcludePatterns []string        `toml:"exclude_patterns"`
}

func Load(name string) (*Config, error) {
	name += ".toml"
	fsys := os.DirFS(config.ConfigPath)
	_, err := fs.Stat(fsys, name)
	if os.IsNotExist(err) {
		fsys = defaults.Configs
		name = defaults.Select(name)
	}

	cfg := &Config{}
	_, err = toml.DecodeFS(fsys, name, cfg)
	if err != nil {
		return nil, fmt.Errorf("game.Load: failed to decode toml: %w", err)
	}

	cfg.SavePath = environ.Replace(cfg.SavePath)
	return cfg, nil
}

func (cfg *Config) Upload() error {
	panic("unimplemented")
}

func (cfg *Config) Download() error {
	panic("unimplemented")
}
