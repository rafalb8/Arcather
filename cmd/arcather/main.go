package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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
	case flag.Setup != "":
		setup(rclone.ToRemoteType(flag.Setup))

	case flag.Remotes:
		remotes()

	case len(flag.Launch) > 0:
		launch(flag.GameName, flag.Launch)

	default:
		err := rclone.Sync("./cmd/arcather", "arcather-test:Arcather", &rclone.Filter{ExcludeRule: []string{"main.go"}})
		if err != nil {
			panic(err)
		}
		fmt.Println("Usage: arcather -- <game_executable>")
		os.Exit(1)
	}
}

func setup(rtype rclone.RemoteType) {
	name := fmt.Sprintf("%s-%s-%s", rclone.RemotePrefix, rtype.String(), "test")
	fmt.Println(rclone.ConfigCreate(name, rtype))
}

func remotes() {
	remotes, err := rclone.ConfigRemotes()
	if err != nil {
		ln.Fatal("Failed to list remotes", ln.Err(err))
	}

	for _, remote := range remotes {
		// Print all remotes with arcather prefix
		if !strings.HasPrefix(remote, rclone.RemotePrefix) {
			continue
		}
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

	ln.Info("Uploading saves to the cloud")
	err = cfg.Sync()
	if err != nil {
		ln.Fatal("Failed to upload:", ln.Err(err))
	}

	ln.Info("Arcather finished")
}
