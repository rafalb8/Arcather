package main

import (
	"os"
	"os/exec"

	"github.com/rafalb8/Arcather/internal/config"
	"github.com/rafalb8/Arcather/internal/game"
	"github.com/rafalb8/ln"
)

func main() {
	config.Init()

	if config.GameName == "" {
		config.GameName = getGameName(config.GameArgs)
	}

	ln.Info("Loading config: " + config.GameName)
	cfg, err := game.Load(config.GameName)
	if err != nil {
		ln.Fatal("Failed to load config", ln.Err(err))
	}

	// Pre-Game Sync
	before := game.Probe(cfg.SavePath)

	// Run the game
	ln.Info("Starting game: " + config.GameName)
	cmd := exec.Command(config.GameArgs[0], config.GameArgs[1:]...)
	if config.Verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	err = cmd.Run()
	if err != nil {
		ln.Error("Game process failed:", ln.Err(err))
	} else {
		ln.Info("Game process finished")
	}

	// Post-Game Sync
	diff := game.Diff(before, game.Probe(cfg.SavePath))
	if len(diff.Created) > 0 || len(diff.Deleted) > 0 || len(diff.Modified) > 0 {
		ln.Info("Uploading saves to the cloud")
		err = cfg.Upload()
		if err != nil {
			ln.Fatal("Failed to upload", ln.Err(err))
		}
	}

	ln.Info("Arcather finished")
}

func getGameName(execPath []string) string {
	return execPath[0]
}
