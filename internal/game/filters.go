package game

import "github.com/rafalb8/arcather/internal/rclone"

type Filters struct {
	Include []string `toml:"include"`
	Exclude []string `toml:"exclude"`
}

func (f Filters) Rclone() *rclone.Filter {
	if len(f.Include)+len(f.Exclude) == 0 {
		return nil
	}

	return &rclone.Filter{
		IncludeRule: f.Include,
		ExcludeRule: f.Exclude,
	}
}
