package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/rafalb8/SaveSync/internal/config"
	"github.com/rafalb8/SaveSync/internal/game"
)

func main() {
	if config.GameName == "" {
		config.GameName = getGameName(config.GameArgs)
	}

	cfg, err := game.Load(config.GameName)
	if err != nil {
		panic(err)
	}

	// Pre-Game Sync
	// fmt.Printf("Syncing latest save for %s from cloud...\n", config.GameName)
	// err = cfg.Download()
	// if err != nil {
	// 	panic(err)
	// }

	// Run the game
	fmt.Println("Starting game:", config.GameName)
	cmd := exec.Command(config.GameArgs[0], config.GameArgs[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Println("Game process failed:", err)
	}

	// Post-Game Sync
	fmt.Printf("Game closed. Uploading new save for %s to cloud...\n", config.GameName)
	err = cfg.Upload()
	if err != nil {
		panic(err)
	}

	fmt.Println("SaveSync finished")
}

func getGameName(execPath []string) string {
	return execPath[0]
}
