package rclone

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/rclone/rclone/backend/drive" // Google Drive
	"github.com/rclone/rclone/librclone/librclone"
)

var (
	Init  = librclone.Initialize
	Close = librclone.Finalize
)

func ConfigRemotes() ([]string, error) {
	out, status := librclone.RPC("config/listremotes", "")
	if status != http.StatusOK {
		return nil, fmt.Errorf("rclone: %s", out)
	}

	resp := &struct {
		Remotes []string
	}{}

	err := json.Unmarshal([]byte(out), resp)
	if err != nil {
		return nil, fmt.Errorf("rclone: %w", err)
	}
	return resp.Remotes, nil
}

type Config struct {
	Scope     string
	TeamDrive string `json:"team_drive"`
	Token     json.RawMessage
	Type      Remote
}

func ConfigGet(name string) (any, error) {
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

func ConfigCreate(name string, remote Remote) error {
	req := &struct {
		Name       string            `json:"name"`
		Parameters map[string]string `json:"parameters"`
		Type       Remote            `json:"type"`
	}{
		Name: name,
		Type: remote,
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
