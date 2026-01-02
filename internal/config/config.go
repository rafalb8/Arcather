package config

import (
	"flag"
	"fmt"
	"log"
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

func Init() {
	// Find game/arcather args split
	split := slices.Index(os.Args, "--")
	if split == -1 {
		fmt.Println("Usage: arcather -- <game_executable>")
		os.Exit(1)
	}

	cfgPath, err := os.UserConfigDir()
	if err != nil {
		log.Fatalln("FATAL:", err)
	}

	// Flags
	flag.StringVar(&ConfigPath, "config-path", filepath.Join(cfgPath, "Arcather"), "path to Arcather config directory")
	flag.StringVar(&GameName, "game", "", "override game name")
	err = flag.CommandLine.Parse(os.Args[1:split])
	if err != nil {
		log.Fatalln("FATAL:", err)
	}
	GameArgs = os.Args[split+1:]

	Environ = *strings.NewReplacer(
		"$HOME", Must(os.UserHomeDir()),
		"$XDG_CONFIG_HOME", cfgPath,
	)
}

func Must[T any](x T, err error) T {
	if err != nil {
		log.Fatalln("FATAL:", err)
	}
	return x
}
