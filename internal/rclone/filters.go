package rclone

type Filter struct {
	// DeleteExcluded bool       `json:",omitempty"`
	// ExcludeFile    []string   `json:",omitempty"`
	// ExcludeFrom    []string   `json:",omitempty"`
	ExcludeRule []string `json:",omitempty"`
	// FilesFrom      []string   `json:",omitempty"`
	// FilesFromRaw   []string   `json:",omitempty"`
	// FilterFrom     []string   `json:",omitempty"`
	// FilterRule     []string   `json:",omitempty"`
	// HashFilter     string     `json:",omitempty"`
	// IgnoreCase     bool       `json:",omitempty"`
	// IncludeFrom    []string   `json:",omitempty"`
	IncludeRule []string `json:",omitempty"`
	// MaxAge         int64      `json:",omitempty"`
	// MaxSize        int        `json:",omitempty"`
	// MetaRules      *MetaRules `json:",omitempty"`
	// MinAge         int64      `json:",omitempty"`
	// MinSize        int        `json:",omitempty"`
}

// type MetaRules struct {
// 	ExcludeFrom []string `json:",omitempty"`
// 	ExcludeRule []string `json:",omitempty"`
// 	FilterFrom  []string `json:",omitempty"`
// 	FilterRule  []string `json:",omitempty"`
// 	IncludeFrom []string `json:",omitempty"`
// 	IncludeRule []string `json:",omitempty"`
// }
