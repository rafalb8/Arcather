package rclone

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	_ "github.com/rclone/rclone/backend/drive" // Google Drive
	_ "github.com/rclone/rclone/backend/local" // Local backend
	_ "github.com/rclone/rclone/fs/sync"       // Sync
	"github.com/rclone/rclone/librclone/librclone"
)

const RemotePrefix = "arcather-"

var (
	Init  = librclone.Initialize
	Close = librclone.Finalize
)

func ConfigRemotes() ([]Remote, error) {
	out, status := librclone.RPC("config/listremotes", "")
	if status != http.StatusOK {
		return nil, fmt.Errorf("rclone: %s", out)
	}

	resp := &struct{ Remotes []Remote }{}
	err := json.Unmarshal([]byte(out), resp)
	if err != nil {
		return nil, fmt.Errorf("rclone: %w", err)
	}

	remotes := make([]Remote, 0, len(resp.Remotes))
	for _, remote := range resp.Remotes {
		if strings.HasPrefix(remote, RemotePrefix) {
			remotes = append(remotes, remote)
		}
	}

	return remotes, nil
}

type Config struct {
	Type      RemoteType
	Scope     string
	TeamDrive string `json:"team_drive"`
	Token     json.RawMessage
}

func ConfigGet(name Remote) (*Config, error) {
	out, status := librclone.RPC("config/get", fmt.Sprintf(`{"name": "%s"}`, name))
	if status != http.StatusOK {
		return nil, fmt.Errorf("rclone: %s", out)
	}

	resp := &Config{}
	err := json.Unmarshal([]byte(out), resp)
	if err != nil {
		return nil, fmt.Errorf("rclone: %w", err)
	}
	return resp, nil
}

func ConfigCreate(name Remote, rtype RemoteType) error {
	req := &struct {
		Name       Remote            `json:"name"`
		Type       RemoteType        `json:"type"`
		Parameters map[string]string `json:"parameters"`
	}{
		Name: fmt.Sprint(RemotePrefix, name),
		Type: rtype,
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

func Sync(src, dst Remote, filter *Filter) error {
	req := &struct {
		Source             Remote `json:"srcFs"`
		Destination        Remote `json:"dstFs"`
		CreateEmptySrcDirs bool   `json:"createEmptySrcDirs"`

		Filter *Filter `json:"_filter,omitempty"`
	}{
		Source:             src,
		Destination:        dst,
		CreateEmptySrcDirs: true,

		Filter: filter,
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
