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
	Name     string          `toml:"-"`
	SavePath string          `toml:"save_path"`
	Filters  Filters         `toml:"filters"`
	Remotes  []rclone.Remote `toml:"remotes"`
}

func Load(path, name string) (*Config, error) {
	fsys, name, err := resolveFS(path, name)
	if err != nil {
		return nil, fmt.Errorf("game.Load: %w", err)
	}

	cfg := &Config{Name: name}
	_, err = toml.DecodeFS(fsys, name+".toml", cfg)
	if err != nil {
		return nil, fmt.Errorf("game.Load: failed to decode toml: %w", err)
	}

	cfg.SavePath = environ.Replace(cfg.SavePath)
	return cfg, nil
}

func resolveFS(path, name string) (fs.FS, string, error) {
	// resolve alt name
	name, err := defaults.Select(name)
	if err != nil {
		return nil, "", err
	}

	// check local config
	localFS := os.DirFS(path)
	_, err = fs.Stat(localFS, name+".toml")
	if err == nil {
		return localFS, name, nil
	}

	if !errors.Is(err, fs.ErrNotExist) {
		return nil, "", fmt.Errorf("failed to check local config: %w", err)
	}

	// check the defaults
	_, err = fs.Stat(defaults.Configs, name+".toml")
	if err == nil {
		return defaults.Configs, name, nil
	}

	if errors.Is(err, fs.ErrNotExist) {
		return nil, "", fmt.Errorf("configuration file %s.toml not found locally or in defaults", name)
	}

	return nil, "", fmt.Errorf("failed to check embedded config: %w", err)
}
