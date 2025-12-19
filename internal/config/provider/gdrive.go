package provider

type GoogleDrive struct {
	Token    string
	BasePath string `toml:"base_path"`
}