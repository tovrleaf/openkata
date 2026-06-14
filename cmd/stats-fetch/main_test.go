package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadCursor(t *testing.T) {
	t.Run("returns default when no file", func(t *testing.T) {
		origDir, _ := os.Getwd()
		tmpDir := t.TempDir()
		os.Chdir(tmpDir)
		defer os.Chdir(origDir)

		c := loadCursor("2026-05-15T00:00:00Z")
		if c.Downloads != "2026-05-15T00:00:00Z" {
			t.Errorf("loadCursor() Downloads = %q, want %q", c.Downloads, "2026-05-15T00:00:00Z")
		}
	})

	t.Run("reads existing cursor", func(t *testing.T) {
		origDir, _ := os.Getwd()
		tmpDir := t.TempDir()
		os.Chdir(tmpDir)
		defer os.Chdir(origDir)

		os.MkdirAll(filepath.Dir(cursorFile), 0o755)
		data, _ := json.Marshal(cursor{
			Downloads: "2026-06-10T00:00:00Z",
			Metrics:   "2026-06-10T00:00:00Z",
			Paths:     "2026-06-10T00:00:00Z",
		})
		os.WriteFile(cursorFile, data, 0o644)

		c := loadCursor("2026-05-15T00:00:00Z")
		if c.Downloads != "2026-06-10T00:00:00Z" {
			t.Errorf("loadCursor() Downloads = %q, want %q", c.Downloads, "2026-06-10T00:00:00Z")
		}
	})
}

func TestSaveCursor(t *testing.T) {
	origDir, _ := os.Getwd()
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	os.MkdirAll(filepath.Dir(cursorFile), 0o755)
	c := cursor{Downloads: "2026-06-14T00:00:00Z", Metrics: "2026-06-14T00:00:00Z", Paths: "2026-06-14T00:00:00Z"}
	saveCursor(c)

	data, err := os.ReadFile(cursorFile)
	if err != nil {
		t.Fatalf("saveCursor() file not created: %v", err)
	}

	var loaded cursor
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("saveCursor() invalid JSON: %v", err)
	}
	if loaded.Downloads != c.Downloads {
		t.Errorf("saveCursor() Downloads = %q, want %q", loaded.Downloads, c.Downloads)
	}
}

func TestLoadJSON(t *testing.T) {
	origDir, _ := os.Getwd()
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	t.Run("returns zero value for missing file", func(t *testing.T) {
		result := loadJSON[[]downloadEvent]("nonexistent.json")
		if result != nil {
			t.Errorf("loadJSON() = %v, want nil", result)
		}
	})

	t.Run("reads valid file", func(t *testing.T) {
		os.WriteFile("test.json", []byte(`[{"artifact":"skills/test","timestamp":"2026-06-14T00:00:00Z"}]`), 0o644)
		result := loadJSON[[]downloadEvent]("test.json")
		if len(result) != 1 {
			t.Errorf("loadJSON() len = %d, want 1", len(result))
		}
		if result[0].Artifact != "skills/test" {
			t.Errorf("loadJSON() artifact = %q, want %q", result[0].Artifact, "skills/test")
		}
	})
}
