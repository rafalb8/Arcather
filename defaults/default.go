package defaults

import (
	"embed"
	"errors"
	"io/fs"
	"log"
	"strings"
)

//go:embed *.toml
var Configs embed.FS

var Alternatives = map[string]string{
	"vintagestory.exe": "vintagestory",
}

func Select(executable string) string {
	exec := strings.ToLower(executable)
	_, err := fs.Stat(Configs, exec)
	if errors.Is(err, fs.ErrNotExist) {
		exec, found := Alternatives[strings.TrimSuffix(exec, ".toml")]
		if !found {
			log.Fatalln("FATAL: No config for", executable, "found")
		}
		return exec + ".toml"
	}
	return exec
}
