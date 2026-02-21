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
	"$XDG_CONFIG_HOME", ln.Must(os.UserConfigDir()),
)

type Config struct {
	Name string `toml:"-"`

	SavePath string          `toml:"save_path"`
	Filters  Filters         `toml:"filters"`
	Remote   []rclone.Remote `toml:"remotes"`
}

func Load(path, name string) (*Config, error) {
	fsys := os.DirFS(path)
	_, err := fs.Stat(fsys, name+".toml")
	if os.IsNotExist(err) {
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

func (cfg *Config) RemotePath(remote rclone.Remote) string {
	return fmt.Sprintf("%s:Arcather/%s", remote, cfg.Name)
}

func (cfg *Config) Sync() error {
	errs := []error{}
	for _, remote := range cfg.Remote {
		err := rclone.Sync(cfg.SavePath, cfg.RemotePath(remote), nil)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

type Filters struct {
	Include []string `toml:"include"`
	Exclude []string `toml:"exclude"`
}

func (f Filters) Rclone() *rclone.Filter {
	if len(f.Include)+len(f.Exclude) == 0 {
		return nil
	}

	return &rclone.Filter{
		IncludeRule: f.Include,
		ExcludeRule: f.Exclude,
	}
}
