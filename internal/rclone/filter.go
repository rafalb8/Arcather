package rclone

// https://rclone.org/rc/#setting-filter-flags-with-filter
// rclone rc --loopback options/get
type Filter struct {
	// DeleteExcluded bool       `json:"DeleteExcluded,omitempty"`
	// ExcludeFile    []string   `json:"ExcludeFile,omitempty"`
	// ExcludeFrom    []string   `json:"ExcludeFrom,omitempty"`
	ExcludeRule []string `json:"ExcludeRule,omitempty" toml:"include"`
	// FilesFrom      []string   `json:"FilesFrom,omitempty"`
	// FilesFromRaw   []string   `json:"FilesFromRaw,omitempty"`
	// FilterFrom     []string   `json:"FilterFrom,omitempty"`
	// FilterRule     []string   `json:"FilterRule,omitempty"`
	// HashFilter     string     `json:"HashFilter,omitempty"`
	// IgnoreCase     bool       `json:"IgnoreCase,omitempty"`
	// IncludeFrom    []string   `json:"IncludeFrom,omitempty"`
	IncludeRule []string `json:"IncludeRule,omitempty" toml:"exclude"`
	// MaxAge         int64      `json:"MaxAge,omitempty"`
	// MaxSize        int        `json:"MaxSize,omitempty"`
	// MetaRules      *MetaRules `json:"MetaRules,omitempty"`
	// MinAge         int64      `json:"MinAge,omitempty"`
	// MinSize        int        `json:"MinSize,omitempty"`
}

// type MetaRules struct {
// 	ExcludeFrom []string `json:"ExcludeFrom,omitempty"`
// 	ExcludeRule []string `json:"ExcludeRule,omitempty"`
// 	FilterFrom  []string `json:"FilterFrom,omitempty"`
// 	FilterRule  []string `json:"FilterRule,omitempty"`
// 	IncludeFrom []string `json:"IncludeFrom,omitempty"`
// 	IncludeRule []string `json:"IncludeRule,omitempty"`
// }
