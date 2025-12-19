package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"
)

// Flags
var (
	ConfigPath string
	GameName   string
)

var GameArgs []string

func init() {
	// Find game/savesync args split
	split := slices.Index(os.Args, "--")

	if split == -1 {
		fmt.Println("Usage: savesync -- <game_executable>")
		os.Exit(1)
	}

	cfgPath, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	cfgPath = filepath.Join(cfgPath, "SaveSync")

	// Flags
	flag.StringVar(&ConfigPath, "config-path", cfgPath, "path to SaveSync config directory")
	flag.StringVar(&GameName, "game", "", "override game name")
	flag.CommandLine.Parse(os.Args[:split])

	GameArgs = os.Args[split+1:]
}
