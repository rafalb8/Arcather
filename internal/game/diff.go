package game

import (
	"io/fs"
	"os"
	"time"

	"github.com/rafalb8/ln"
)

type Stat struct {
	Path    string
	ModTime time.Time
	// In the future, this will also have a Hash property.
	Hash string
}

func Probe(path string) []Stat {
	stats := []Stat{}
	fs.WalkDir(os.DirFS(path), ".", func(path string, d fs.DirEntry, err error) error {
		if d.Type().IsDir() {
			// Enter dir
			return nil
		}

		info, err := d.Info()
		if err != nil {
			ln.Error("Error getting file info", ln.Err(err))
			return nil
		}

		stats = append(stats, Stat{
			Path:    path,
			ModTime: info.ModTime(),
		})
		return nil
	})
	return stats
}

type Result struct {
	Unchanged []Stat
	Modified  []Stat
	Created   []Stat
	Deleted   []Stat
}

func Diff(before, after []Stat) Result {
	beforeMap := make(map[string]Stat)
	for _, s := range before {
		beforeMap[s.Path] = s
	}

	afterMap := make(map[string]Stat)
	for _, s := range after {
		afterMap[s.Path] = s
	}

	result := Result{}

	for path, beforeStat := range beforeMap {
		if afterStat, ok := afterMap[path]; ok {
			if beforeStat.ModTime.Equal(afterStat.ModTime) {
				result.Unchanged = append(result.Unchanged, beforeStat)
				continue
			}
			result.Modified = append(result.Modified, afterStat)
			continue
		}
		result.Deleted = append(result.Deleted, beforeStat)
	}

	for path, afterStat := range afterMap {
		if _, ok := beforeMap[path]; !ok {
			result.Created = append(result.Created, afterStat)
		}
	}

	return result
}
