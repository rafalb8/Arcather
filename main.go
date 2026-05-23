package main

import (
	"fmt"

	"github.com/rafalb8/Arcather/internal/flag"
	"github.com/rafalb8/Arcather/internal/rclone"
	"github.com/rafalb8/Arcather/internal/runner"
	"github.com/rafalb8/ln"
)

func main() {
	ln.Default = ln.New(ln.WithMultiline(true))
	
	rclone.Init()
	defer rclone.Close()

	switch {
	case flag.Setup:
		setup(flag.SetupName, rclone.ToType(flag.SetupType))

	case flag.Remotes:
		remotes()
	}

	r, err := runner.New(ln.Default, flag.ConfigPath, flag.Verbose)
	if err != nil {
		ln.Fatal("Failed to init game runner", ln.Err(err))
	}

	err = r.Launch(flag.GameName, flag.GameArgs)
	if err != nil {
		ln.Fatal("Session crashed", ln.Err(err))
	}
}

func setup(name string, rtype rclone.Type) {
	err := rclone.AddRemote(name, rtype)
	if err != nil {
		ln.Fatal("Failed to setup remote", ln.Err(err))
	}
}

func remotes() {
	remotes, err := rclone.ListRemotes()
	if err != nil {
		ln.Fatal("Failed to list remotes", ln.Err(err))
	}

	for _, remote := range remotes {
		fmt.Printf("[%s] %s\n", remote.Type.String(), remote.Name)
	}
}
