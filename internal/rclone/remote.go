package rclone

import (
	"strings"
)

type Remote int8

const (
	Unsupported = iota - 1
	Local
	SSH
	GoogleDrive
)

func ToRemote(t string) Remote {
	switch strings.ToLower(t) {
	case "local":
		return Local
	case "ssh":
		return SSH
	case "gdrive", "google", "gd", "drive":
		return GoogleDrive
	default:
		return Unsupported
	}
}

func (r Remote) String() string {
	switch r {
	case Local:
		return "local"
	case SSH:
		return "ssh"
	case GoogleDrive:
		return "drive"
	default:
		return ""
	}
}

func (r Remote) MarshalJSON() ([]byte, error) {
	return []byte(`"` + r.String() + `"`), nil
}

func (r *Remote) UnmarshalJSON(data []byte) error {
	*r = ToRemote(strings.Trim(string(data), `"`))
	return nil
}
