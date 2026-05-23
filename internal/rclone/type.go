package rclone

import "strings"

type Type int8

const (
	Unsupported Type = iota - 1
	Local
	SSH
	GoogleDrive
)

func ToType(t string) Type {
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

func (rt Type) String() string {
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

func (rt Type) MarshalText() ([]byte, error) {
	return []byte(rt.String()), nil
}

func (rt *Type) UnmarshalText(text []byte) error {
	*rt = ToType(string(text))
	return nil
}

func (rt Type) MarshalJSON() ([]byte, error) {
	return []byte(`"` + rt.String() + `"`), nil
}

func (rt *Type) UnmarshalJSON(data []byte) error {
	*rt = ToType(strings.Trim(string(data), `"`))
	return nil
}
