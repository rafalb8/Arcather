package game

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/rafalb8/Arcather/defaults"
	"github.com/rafalb8/Arcather/internal/rclone"
	"github.com/rafalb8/ln"
)

var environ = strings.NewReplacer(
	"$HOME", ln.Must(os.UserHomeDir()),
	"$CONFIG", ln.Must(os.UserConfigDir()),
)

type Config struct {
	Name     string  `toml:"-"` // Name of the game
	SavePath string  `toml:"save_path"`
	Filters  Filters `toml:"filters"`

	Remotes []rclone.Remote `toml:"remotes"`
}

func Load(path, name string) (*Config, error) {
	fsys := os.DirFS(path)
	_, err := fs.Stat(fsys, name+".toml")
	if err != nil {
		fsys = defaults.Configs
		name = defaults.Select(name)
	}

	cfg := &Config{Name: name}
	_, err = toml.DecodeFS(fsys, name+".toml", cfg)
	if err != nil {
		return nil, fmt.Errorf("game.Load: failed to decode toml: %w", err)
	}

	cfg.SavePath = environ.Replace(cfg.SavePath)
	return cfg, nil
}

func (cfg *Config) Sync() error {
	if len(cfg.Remotes) == 0 {
		var err error
		cfg.Remotes, err = rclone.ListRemotes()
		if err != nil {
			return fmt.Errorf("game.Sync: failed to load remotes: %w", err)
		}
	}

	local := rclone.Remote{
		Type: rclone.Local,
		Path: cfg.SavePath,
	}

	errs := []error{}
	for _, remote := range cfg.Remotes {
		err := rclone.Sync(local, remote.JoinPath(cfg.Name), cfg.Filters.Rclone())
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
