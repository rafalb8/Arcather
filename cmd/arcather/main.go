package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/rafalb8/Arcather/internal/config"
	"github.com/rafalb8/Arcather/internal/game"
	"github.com/rafalb8/Arcather/internal/provider"
	"github.com/rafalb8/Arcather/internal/rclone"
	"github.com/rafalb8/ln"
)

func main() {
	rclone.Init()
	defer rclone.Close()
	config.Init()

	switch {
	case config.Setup != "":
		setup(provider.ToType(config.Setup))
	case len(config.Launch) > 0:
		launch(config.GameName, config.Launch)
	default:
		fmt.Println("Usage: arcather -- <game_executable>")
		os.Exit(1)
	}
}

func setup(p provider.Type) {
	fmt.Println(rclone.ListRemotes())
}

func launch(name string, args []string) {
	if name == "" {
		name = args[0]
	}

	ln.Info("Loading config: " + name)
	cfg, err := game.Load(name)
	if err != nil {
		ln.Fatal("Failed to load config", ln.Err(err))
	}

	// Pre-Game Sync
	before := game.Probe(cfg.SavePath)

	// Run the game
	ln.Info("Starting game: " + name)
	cmd := exec.Command(args[0], args[1:]...)
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
