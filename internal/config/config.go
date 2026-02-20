package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/rafalb8/ln"
)

// Flags
var (
	ConfigPath string
	GameName   string
	Verbose    bool
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

	cfgPath := ln.Must(os.UserConfigDir())

	// Flags
	flag.StringVar(&ConfigPath, "config-path", filepath.Join(cfgPath, "Arcather"), "path to Arcather config directory")
	flag.StringVar(&GameName, "game", "", "override game name")
	flag.BoolVar(&Verbose, "v", false, "show game logs")
	err := flag.CommandLine.Parse(os.Args[1:split])
	if err != nil {
		ln.Fatal("Failed to parse flags", ln.Err(err))
	}
	GameArgs = os.Args[split+1:]

	Environ = *strings.NewReplacer(
		"$HOME", ln.Must(os.UserHomeDir()),
		"$XDG_CONFIG_HOME", cfgPath,
	)
}
