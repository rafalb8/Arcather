package defaults

import (
	"embed"
)

//go:embed *.toml
var Configs embed.FS
