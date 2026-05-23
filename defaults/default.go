package defaults

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"strings"
)

//go:embed *.toml
var Configs embed.FS

var Alternatives = map[string]string{
	"vintagestory.exe": "vintagestory",
}

// Select looks up a configuration by executable name.
// It returns the normalized name or an error if not found.
func Select(executable string) (string, error) {
	exec := strings.ToLower(executable)

	_, err := fs.Stat(Configs, exec+".toml")
	if err == nil {
		return exec, nil
	}

	if errors.Is(err, fs.ErrNotExist) {
		alt, found := Alternatives[exec]
		if !found {
			return "", fmt.Errorf("no default configuration found for executable: %s", executable)
		}
		return alt, nil
	}

	return "", fmt.Errorf("failed to stat embedded config: %w", err)
}
