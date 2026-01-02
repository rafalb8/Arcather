package game

import (
	"reflect"
	"testing"
	"time"
)

func TestDiff(t *testing.T) {
	now := time.Now()
	oneHourAgo := now.Add(-time.Hour)
	twoHoursAgo := now.Add(-2 * time.Hour)

	// Test case 1: Basic case with all types of changes
	t.Run("Basic", func(t *testing.T) {
		before := []Stat{
			{Path: "a", ModTime: oneHourAgo},
			{Path: "b", ModTime: oneHourAgo},
			{Path: "c", ModTime: oneHourAgo},
		}
		after := []Stat{
			{Path: "a", ModTime: oneHourAgo}, // Unchanged
			{Path: "b", ModTime: now},        // Modified
			{Path: "d", ModTime: now},        // Created
		}

		expected := Result{
			Unchanged: []Stat{{Path: "a", ModTime: oneHourAgo}},
			Modified:  []Stat{{Path: "b", ModTime: now}},
			Created:   []Stat{{Path: "d", ModTime: now}},
			Deleted:   []Stat{{Path: "c", ModTime: oneHourAgo}},
		}

		result := Diff(before, after)

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Diff() = %v, want %v", result, expected)
		}
	})

	// Test case 2: No changes
	t.Run("NoChanges", func(t *testing.T) {
		before := []Stat{
			{Path: "a", ModTime: oneHourAgo},
			{Path: "b", ModTime: twoHoursAgo},
		}
		after := []Stat{
			{Path: "a", ModTime: oneHourAgo},
			{Path: "b", ModTime: twoHoursAgo},
		}

		expected := Result{
			Unchanged: []Stat{
				{Path: "a", ModTime: oneHourAgo},
				{Path: "b", ModTime: twoHoursAgo},
			},
		}

		result := Diff(before, after)

		if !reflect.DeepEqual(result.Unchanged, expected.Unchanged) {
			t.Errorf("Unchanged: got %v, want %v", result.Unchanged, expected.Unchanged)
		}
		if len(result.Modified) != 0 {
			t.Errorf("Expected no modified files, got %v", result.Modified)
		}
		if len(result.Created) != 0 {
			t.Errorf("Expected no created files, got %v", result.Created)
		}
		if len(result.Deleted) != 0 {
			t.Errorf("Expected no deleted files, got %v", result.Deleted)
		}
	})

	// Test case 3: All files created
	t.Run("AllCreated", func(t *testing.T) {
		before := []Stat{}
		after := []Stat{
			{Path: "a", ModTime: now},
			{Path: "b", ModTime: now},
		}

		expected := Result{
			Created: []Stat{
				{Path: "a", ModTime: now},
				{Path: "b", ModTime: now},
			},
		}

		result := Diff(before, after)

		if !reflect.DeepEqual(result.Created, expected.Created) {
			t.Errorf("Created: got %v, want %v", result.Created, expected.Created)
		}
		if len(result.Unchanged) != 0 {
			t.Errorf("Expected no unchanged files, got %v", result.Unchanged)
		}
		if len(result.Modified) != 0 {
			t.Errorf("Expected no modified files, got %v", result.Modified)
		}
		if len(result.Deleted) != 0 {
			t.Errorf("Expected no deleted files, got %v", result.Deleted)
		}
	})

	// Test case 4: All files deleted
	t.Run("AllDeleted", func(t *testing.T) {
		before := []Stat{
			{Path: "a", ModTime: oneHourAgo},
			{Path: "b", ModTime: oneHourAgo},
		}
		after := []Stat{}

		expected := Result{
			Deleted: []Stat{
				{Path: "a", ModTime: oneHourAgo},
				{Path: "b", ModTime: oneHourAgo},
			},
		}

		result := Diff(before, after)

		if !reflect.DeepEqual(result.Deleted, expected.Deleted) {
			t.Errorf("Deleted: got %v, want %v", result.Deleted, expected.Deleted)
		}
		if len(result.Unchanged) != 0 {
			t.Errorf("Expected no unchanged files, got %v", result.Unchanged)
		}
		if len(result.Modified) != 0 {
			t.Errorf("Expected no modified files, got %v", result.Modified)
		}
		if len(result.Created) != 0 {
			t.Errorf("Expected no created files, got %v", result.Created)
		}
	})

	// Test case 5: All files modified
	t.Run("AllModified", func(t *testing.T) {
		before := []Stat{
			{Path: "a", ModTime: oneHourAgo},
			{Path: "b", ModTime: oneHourAgo},
		}
		after := []Stat{
			{Path: "a", ModTime: now},
			{Path: "b", ModTime: now},
		}

		expected := Result{
			Modified: []Stat{
				{Path: "a", ModTime: now},
				{Path: "b", ModTime: now},
			},
		}

		result := Diff(before, after)

		if !reflect.DeepEqual(result.Modified, expected.Modified) {
			t.Errorf("Modified: got %v, want %v", result.Modified, expected.Modified)
		}
		if len(result.Unchanged) != 0 {
			t.Errorf("Expected no unchanged files, got %v", result.Unchanged)
		}
		if len(result.Created) != 0 {
			t.Errorf("Expected no created files, got %v", result.Created)
		}
		if len(result.Deleted) != 0 {
			t.Errorf("Expected no deleted files, got %v", result.Deleted)
		}
	})
}
