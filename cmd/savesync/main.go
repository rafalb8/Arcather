package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/rafalb8/SaveSync/internal/config"
)

func main() {
	if config.GameName == "" {
		config.GameName = getGameName(config.GameArgs)
	}

	cfg, err := config.LoadGame()
	if err != nil {
		panic(err)
	}

	// Pre-Game Sync
	fmt.Printf("Syncing latest save for %s from cloud...\n", config.GameName)
	// performDownloadSync(gameConfig, cloudConfig)
	_ = cfg

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
	// performUploadSync(gameConfig, cloudConfig)

	fmt.Println("SaveSync finished")
}

func getGameName(execPath []string) string {
	return execPath[0]
}
