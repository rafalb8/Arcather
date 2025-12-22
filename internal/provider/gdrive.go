package provider

type googleDrive struct {
	Token    string
	BasePath string `toml:"base_path"`
}
