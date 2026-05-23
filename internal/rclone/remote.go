package rclone

import (
	"fmt"
	"path"
	"strings"
)

type Remote struct {
	Type     Type   `json:"type"`
	Name     string `json:"name"` // e.g., "arcather-gdrive" (or "local" for local storage)
	Path     string `json:"path"` // e.g., "/backup/game-name" or "C:/Saves"
	Compress *bool  `json:"compress"`
}

// RcloneString generates the connection string or filesystem descriptor that rclone RPC expects
func (r Remote) Rclone() string {
	// if is local, return Path
	if r.Type == Local || r.Name == "" || strings.ToLower(r.Name) == "local" {
		return r.Path
	}

	// build the standard remote string (e.g., "arcather-gdrive:backup/game")
	rStr := r.Name + ":"
	if r.Compress != nil && *r.Compress {
		rStr = fmt.Sprintf(`:compress,remote='%s',mode='zstd',level='2':`, rStr)
	}

	// append path at the end
	if r.Path != "" {
		rStr += strings.TrimPrefix(r.Path, "/")
	}
	return rStr
}

func (r *Remote) FillDefault() {
	r.Path = "Arcather"
	r.Compress = new(true)
}

func (r Remote) JoinPath(p string) Remote {
	r.Path = path.Join(r.Path, p)
	return r
}
