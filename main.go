package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rafalb8/Arcather/internal/flag"
	"github.com/rafalb8/Arcather/internal/game"
	"github.com/rafalb8/Arcather/internal/rclone"
	"github.com/rafalb8/ln"
)

func main() {
	rclone.Init()
	defer rclone.Close()
	flag.Init()

	switch {
	case flag.Setup:
		setup(flag.SetupName, rclone.ToRemoteType(flag.SetupType))

	case flag.Remotes:
		remotes()

	case len(flag.Launch) > 0:
		launch(flag.GameName, flag.Launch)
	}
}

func setup(name string, rtype rclone.RemoteType) {
	err := rclone.ConfigCreate(name, rtype)
	if err != nil {
		ln.Fatal("Failed to setup remote", ln.Err(err))
	}
}

func remotes() {
	remotes, err := rclone.ConfigRemotes()
	if err != nil {
		ln.Fatal("Failed to list remotes", ln.Err(err))
	}

	for _, remote := range remotes {
		fmt.Println(remote)
	}
}

func launch(name string, args []string) {
	if name == "" {
		name = filepath.Base(args[0])
	}

	ln.Info("Loading config: " + name)
	cfg, err := game.Load(flag.ConfigPath, name)
	if err != nil {
		ln.Fatal("Failed to load config", ln.Err(err))
	}

	// TODO: Pre-Game Sync

	ln.Info("Starting game: " + name)
	cmd := exec.Command(args[0], args[1:]...)
	if flag.Verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	err = cmd.Run()
	if err != nil {
		ln.Error("Game process failed:", ln.Err(err))
	} else {
		ln.Info("Game process finished")
	}

	ln.Info("Syncing saves with the cloud")
	err = cfg.Sync()
	if err != nil {
		ln.Fatal("Failed to upload:", ln.Err(err))
	}

	ln.Info("Arcather finished")
}
