package config

import (
	"flag"
	"os"
	"path/filepath"
	"slices"

	"github.com/rafalb8/ln"
)

// Flags
var (
	ConfigPath string
	GameName   string
	Verbose    bool
)

// Modes
var (
	Setup   string
	Remotes bool
	Launch  []string
)

func Init() {
	// Find game/arcather args split
	split := slices.Index(os.Args, "--")
	if split == -1 {
		split = len(os.Args)
	}

	cfgPath := ln.Must(os.UserConfigDir())

	// Flags
	flag.StringVar(&ConfigPath, "config-path", filepath.Join(cfgPath, "Arcather"), "path to Arcather config directory")
	flag.StringVar(&GameName, "game", "", "override game name")
	flag.BoolVar(&Verbose, "v", false, "show game logs")

	// Modes
	flag.StringVar(&Setup, "setup", "", "setup remote")
	flag.BoolVar(&Remotes, "remotes", false, "list remotes")

	err := flag.CommandLine.Parse(os.Args[1:split])
	if err != nil {
		ln.Fatal("Failed to parse flags", ln.Err(err))
	}

	if split != len(os.Args) {
		Launch = os.Args[split+1:]
	}
}
