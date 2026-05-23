package flag

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/rafalb8/ln"
	"github.com/spf13/pflag"
)

var (
	ConfigPath string
	Verbose    bool

	GameName string
	GameArgs []string
)

var (
	Setup     bool
	SetupName string
	SetupType string
)

var Remotes bool

func init() {
	if len(os.Args) == 1 {
		os.Args = append(os.Args, "-h")
	}

	split := slices.Index(os.Args, "--")
	if split == -1 {
		split = len(os.Args)
	}

	args := os.Args[1:split]

	var sub string
	if len(args) > 0 {
		sub = args[0]
	}

	switch sub {
	case "setup":
		Setup = true
		setupFlags(args[1:])
	case "remotes":
		Remotes = true
	default:
		mainFlags(args)
	}

	if split != len(os.Args) {
		GameArgs = os.Args[split+1:]
	}
}

func mainFlags(args []string) {
	fs := pflag.NewFlagSet("Arcather", pflag.ContinueOnError)

	cfgPath := ln.Must(os.UserConfigDir())

	fs.StringVarP(&ConfigPath, "config-path", "c", filepath.Join(cfgPath, "Arcather"), "path to Arcather config directory")
	fs.StringVarP(&GameName, "game", "n", "", "override game name")
	fs.BoolVarP(&Verbose, "verbose", "v", false, "show game logs")

	fs.Usage = func() {
		fmt.Println("Usage:")
		fmt.Println("  arcather <command> [options]")
		fmt.Println("  arcather [options] -- <game_executable> [args...]")

		fmt.Println("\nCommands:")
		fmt.Println("  setup           Setup a new remote")
		fmt.Println("  remotes         List Arcather remotes")

		fmt.Println("\nOptions:")
		fs.PrintDefaults()
	}
	err := fs.Parse(args)
	if err != nil {
		if err == pflag.ErrHelp {
			os.Exit(0)
		}
		ln.Fatal("Failed to parse flags", ln.Err(err))
	}
}

func setupFlags(args []string) {
	fs := pflag.NewFlagSet("arcather setup", pflag.ContinueOnError)
	fs.Usage = func() {
		fmt.Println("Usage: arcather setup <name> <type>")
		fmt.Println("\nArguments:")
		fmt.Println("  <name>          The name of the remote")
		fmt.Println("  <type>          The type (local, ssh, gdrive, etc.)")
		fs.PrintDefaults()
	}

	err := fs.Parse(args)
	if err != nil {
		if err == pflag.ErrHelp {
			os.Exit(0)
		}
		ln.Fatal("Failed to parse flags", ln.Err(err))
	}

	args = fs.Args()
	if len(args) < 2 {
		fs.Usage()
		os.Exit(1)
	}

	SetupName = args[0]
	SetupType = args[1]
}
