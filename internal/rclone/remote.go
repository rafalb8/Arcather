package rclone

import (
	"strings"
)

type (
	Remote     = string
	RemoteType int8
)

const (
	Unsupported RemoteType = iota - 1
	Local
	SSH
	GoogleDrive
)

func ToRemoteType(t string) RemoteType {
	switch strings.ToLower(t) {
	case "local":
		return Local
	case "ssh":
		return SSH
	case "drive", "gdrive", "google", "gd":
		return GoogleDrive
	default:
		return Unsupported
	}
}

func (rt RemoteType) String() string {
	switch rt {
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

func (rt RemoteType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + rt.String() + `"`), nil
}

func (rt *RemoteType) UnmarshalJSON(data []byte) error {
	*rt = ToRemoteType(strings.Trim(string(data), `"`))
	return nil
}
