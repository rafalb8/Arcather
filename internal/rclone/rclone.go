package rclone

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rclone/rclone/librclone/librclone"
)

var (
	Init  = librclone.Initialize
	Close = librclone.Finalize
)

func ListRemotes() ([]string, error) {
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
