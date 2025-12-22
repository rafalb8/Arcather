package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// Flags
var (
	ConfigPath string
	GameName   string
)

var (
	GameArgs []string
	Environ  strings.Replacer
)

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

	// Flags
	flag.StringVar(&ConfigPath, "config-path", filepath.Join(cfgPath, "SaveSync"), "path to SaveSync config directory")
	flag.StringVar(&GameName, "game", "", "override game name")
	flag.CommandLine.Parse(os.Args[:split])

	GameArgs = os.Args[split+1:]

	Environ = *strings.NewReplacer(
		"$HOME", Must(os.UserHomeDir()),
		"$XDG_CONFIG_HOME", cfgPath,
	)
}

func Must[T any](x T, err error) T {
	if err != nil {
		panic(err)
	}
	return x
}
