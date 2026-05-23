package runner

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rafalb8/Arcather/internal/game"
	"github.com/rafalb8/Arcather/internal/rclone"
	"github.com/rafalb8/ln"
)

type Runner struct {
	log     ln.Logger
	remotes []rclone.Remote
	cfgPath string
	verbose bool
}

func New(log ln.Logger, cfgPath string, verbose bool) (*Runner, error) {
	remotes, err := rclone.ListRemotes()
	if err != nil {
		return nil, fmt.Errorf("runner.New: failed to load default remotes: %w", err)
	}
	return &Runner{
		log:     log,
		remotes: remotes,
		cfgPath: cfgPath,
		verbose: verbose,
	}, nil
}

// Launch handles the entire lifecycle: config loading, pre-sync, execution, and post-sync.
func (r *Runner) Launch(name string, args []string) error {
	if name == "" {
		name = filepath.Base(args[0])
	}

	r.log.Info("Loading config: " + name)
	cfg, err := game.Load(r.cfgPath, name)
	if err != nil {
		return fmt.Errorf("runner.Launch: failed to load config: %w", err)
	}

	r.log.Info("Syncing saves from cloud...")
	err = r.Download(cfg)
	if err != nil {
		// log the error, but don't stop
		r.log.Error("Pre-game sync failed:", ln.Err(err))
	}

	r.log.Info("Starting game process: " + strings.Join(args, " "))
	cmd := exec.Command(args[0], args[1:]...)
	if r.verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	err = cmd.Run()
	if err != nil {
		ln.Error("Game process failed:", ln.Err(err))
	} else {
		ln.Info("Game process finished")
	}

	r.log.Info("Syncing saves to cloud...")
	err = r.Sync(cfg)
	if err != nil {
		return fmt.Errorf("runner.Launch: post-game sync failed: %w", err)
	}

	r.log.Info("Arcather finished session")
	return nil
}

func (r *Runner) Download(cfg *game.Config) error {
	ln.Warn("Download not implemented!")
	return nil
}

// Sync handles the execution logic for a specific game config
func (r *Runner) Sync(cfg *game.Config) error {
	remotes := cfg.Remotes

	// if the specific game config has no remotes, use defaults
	if len(remotes) == 0 {
		remotes = r.remotes
	}

	local := rclone.Remote{
		Type: rclone.Local,
		Path: cfg.SavePath,
	}

	var errs []error
	for _, dst := range remotes {
		dst = dst.JoinPath(cfg.Name)

		r.log.Debug("Syncing saves",
			ln.String("src", local.Rclone()),
			ln.String("dst", dst.Rclone()),
		)

		err := rclone.Sync(local, dst, cfg.Filters.Rclone())
		if err != nil {
			errs = append(errs, fmt.Errorf("syncer: failed syncing to %s: %w", dst.Path, err))
		}
	}

	return errors.Join(errs...)
}
