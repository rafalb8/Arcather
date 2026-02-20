package provider

import "strings"

type Type int8

const (
	Unsupported = iota - 1
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
	case "gdrive", "google", "gd":
		return GoogleDrive
	default:
		return Unsupported
	}
}
