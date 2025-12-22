package provider

type Type int8

const (
	Unsupported = iota - 1
	Local
	SSH
	GoogleDrive
)