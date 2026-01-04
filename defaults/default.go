package defaults

import (
	"embed"
	"errors"
	"io/fs"
	"strings"

	"github.com/rafalb8/ln"
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
			ln.Fatal("No config found", ln.String("exe", executable))
		}
		return exec + ".toml"
	}
	return exec
}
