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
	GameName   string
	Verbose    bool
	Launch     []string
)

var (
	SetupName string
	SetupType string
)

var Remotes bool

func Init() {
	if len(os.Args) == 1 {
		os.Args = append(os.Args, "-h")
	}

	split := slices.Index(os.Args, "--")
	if split == -1 {
		split = len(os.Args)
	}

	var err error
	args := os.Args[1:split]
	if len(args) > 0 {
		switch args[0] {
		case "setup":
			err = setupFlags(args[1:])
		case "remotes":
			Remotes = true
		default:
			err = mainFlags(args)
		}
	}

	if err != nil && err != pflag.ErrHelp {
		ln.Fatal("Failed to parse flags", ln.Err(err))
	}

	if split != len(os.Args) {
		Launch = os.Args[split+1:]
	}
}

func mainFlags(args []string) error {
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
	return fs.Parse(args)
}

func setupFlags(args []string) error {
	fs := pflag.NewFlagSet("arcather setup", pflag.ContinueOnError)

	fs.StringVarP(&SetupName, "name", "n", "", "remote name")
	fs.StringVarP(&SetupType, "type", "t", "", "remote type (sftp, webdav, etc.)")

	fs.Usage = func() {
		fmt.Println("Usage: arcather setup [options]")
		fmt.Println("\nOptions:")
		fs.PrintDefaults()
	}
	return fs.Parse(args)
}
