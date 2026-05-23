package rclone

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	_ "github.com/rclone/rclone/backend/compress" // Transparent compression wrapper
	_ "github.com/rclone/rclone/backend/drive"    // Google Drive
	_ "github.com/rclone/rclone/backend/local"    // Local backend
	_ "github.com/rclone/rclone/fs/sync"          // Sync
	"github.com/rclone/rclone/librclone/librclone"
)

// RPC docs: https://rclone.org/rc/

const RemotePrefix = "arcather-"

var (
	Init  = librclone.Initialize
	Close = librclone.Finalize
)

func AddRemote(name string, rtype Type) error {
	req := &struct {
		Name       string            `json:"name"`
		Type       Type        `json:"type"`
		Parameters map[string]string `json:"parameters"`
	}{
		Name:       fmt.Sprint(RemotePrefix, name),
		Type:       rtype,
		Parameters: make(map[string]string),
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("rclone: %w", err)
	}

	out, status := librclone.RPC("config/create", string(payload))
	if status != http.StatusOK {
		return fmt.Errorf("rclone: %s", out)
	}
	return nil
}

// rclone rc --loopback config/dump
func ListRemotes() ([]Remote, error) {
	out, status := librclone.RPC("config/dump", "")
	if status != http.StatusOK {
		return nil, fmt.Errorf("rclone: %s", out)
	}

	type rcfg struct{ Type Type }
	resp := map[string]rcfg{}
	err := json.Unmarshal([]byte(out), &resp)
	if err != nil {
		return nil, fmt.Errorf("rclone: %w", err)
	}

	remotes := make([]Remote, 0, len(resp))
	for name, cfg := range resp {
		if strings.HasPrefix(name, RemotePrefix) {
			r := Remote{Type: cfg.Type, Name: name}
			r.FillDefault()
			remotes = append(remotes, r)
		}
	}

	return remotes, nil
}

type syncReq struct {
	Source             string  `json:"srcFs"`
	Destination        string  `json:"dstFs"`
	CreateEmptySrcDirs bool    `json:"createEmptySrcDirs"`
	Filter             *Filter `json:"_filter,omitempty"`
}

// Sync the source to the destination, changing the destination only
// https://rclone.org/commands/rclone_sync/
func Sync(src, dst Remote, filter *Filter) error {
	req := &syncReq{
		Source:             src.Rclone(),
		Destination:        dst.Rclone(),
		CreateEmptySrcDirs: true,
		Filter:             filter,
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("rclone: %w", err)
	}

	out, status := librclone.RPC("sync/sync", string(payload))
	if status != http.StatusOK {
		return fmt.Errorf("rclone: %s", out)
	}
	return nil
}
