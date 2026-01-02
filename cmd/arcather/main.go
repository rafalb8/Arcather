package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/rafalb8/Arcather/internal/config"
	"github.com/rafalb8/Arcather/internal/game"
)

func main() {
	config.Init()

	if config.GameName == "" {
		config.GameName = getGameName(config.GameArgs)
	}

	log.Println("INFO: Loading config for", config.GameName)
	cfg, err := game.Load(config.GameName)
	if err != nil {
		log.Fatalln("FATAL:", err)
	}

	// Pre-Game Sync
	before := game.Probe(cfg.SavePath)
	// log.Printf("Syncing latest save for %s from cloud...\n", config.GameName)
	// err = cfg.Download()
	// if err != nil {
	// 	log.Fatalln("FATAL:", err)
	// }

	// Run the game
	log.Println("INFO: Starting game:", config.GameName)
	cmd := exec.Command(config.GameArgs[0], config.GameArgs[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		log.Println("ERROR: Game process failed:", err)
	} else {
		log.Println("INFO: Game process finished.")
	}

	// Post-Game Sync
	diff := game.Diff(before, game.Probe(cfg.SavePath))
	if len(diff.Created) > 0 || len(diff.Deleted) > 0 || len(diff.Modified) > 0 {
		log.Printf("INFO: Uploading new save for %s to cloud...\n", config.GameName)
		err = cfg.Upload()
		if err != nil {
			log.Fatalln("FATAL:", err)
		}
	}

	log.Println("INFO: Arcather finished.")
}

func getGameName(execPath []string) string {
	return execPath[0]
}
