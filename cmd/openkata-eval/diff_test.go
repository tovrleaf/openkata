package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSnapshotDir(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "sub"), 0755)
	os.WriteFile(filepath.Join(dir, "a.txt"), []byte("hello"), 0644)
	os.WriteFile(filepath.Join(dir, "sub/b.txt"), []byte("world"), 0644)

	snap, err := snapshotDir(dir)
	if err != nil {
		t.Fatalf("snapshotDir() error: %v", err)
	}
	if len(snap) != 2 {
		t.Fatalf("snapshotDir() got %d files, want 2", len(snap))
	}
	if snap["a.txt"] != "hello" {
		t.Errorf("snap[a.txt] = %q, want %q", snap["a.txt"], "hello")
	}
	if snap["sub/b.txt"] != "world" {
		t.Errorf("snap[sub/b.txt] = %q, want %q", snap["sub/b.txt"], "world")
	}
}

func TestSnapshotDirEmpty(t *testing.T) {
	dir := t.TempDir()
	snap, err := snapshotDir(dir)
	if err != nil {
		t.Fatalf("snapshotDir() error: %v", err)
	}
	if len(snap) != 0 {
		t.Errorf("snapshotDir() got %d files, want 0", len(snap))
	}
}

func TestComputeDiff(t *testing.T) {
	before := map[string]string{
		"keep.txt":    "same",
		"modify.txt":  "old",
		"deleted.txt": "gone",
	}
	after := map[string]string{
		"keep.txt":   "same",
		"modify.txt": "new",
		"added.txt":  "fresh",
	}

	diff := computeDiff(before, after)

	if len(diff.Added) != 1 || diff.Added[0].Path != "added.txt" {
		t.Errorf("Added = %v, want [added.txt]", diff.Added)
	}
	if diff.Added[0].Content != "fresh" {
		t.Errorf("Added[0].Content = %q, want %q", diff.Added[0].Content, "fresh")
	}

	if len(diff.Modified) != 1 || diff.Modified[0].Path != "modify.txt" {
		t.Errorf("Modified = %v, want [modify.txt]", diff.Modified)
	}
	if diff.Modified[0].Content != "new" {
		t.Errorf("Modified[0].Content = %q, want %q", diff.Modified[0].Content, "new")
	}

	if len(diff.Deleted) != 1 || diff.Deleted[0] != "deleted.txt" {
		t.Errorf("Deleted = %v, want [deleted.txt]", diff.Deleted)
	}

	if len(diff.Unchanged) != 1 || diff.Unchanged[0] != "keep.txt" {
		t.Errorf("Unchanged = %v, want [keep.txt]", diff.Unchanged)
	}
}

func TestComputeDiffNoChanges(t *testing.T) {
	snap := map[string]string{"a.txt": "content"}
	diff := computeDiff(snap, snap)

	if len(diff.Added) != 0 || len(diff.Modified) != 0 || len(diff.Deleted) != 0 {
		t.Errorf("expected no changes, got added=%d modified=%d deleted=%d",
			len(diff.Added), len(diff.Modified), len(diff.Deleted))
	}
	if len(diff.Unchanged) != 1 {
		t.Errorf("Unchanged = %d, want 1", len(diff.Unchanged))
	}
}

func TestFormatDiffForJudge(t *testing.T) {
	diff := DiffSummary{
		Added:    []FileEntry{{Path: "new.md", Content: "# New"}},
		Modified: []FileEntry{{Path: "old.md", Content: "# Updated"}},
		Deleted:  []string{"removed.txt"},
	}

	result := formatDiffForJudge("agent output here", diff)

	if !strings.Contains(result, "## Agent Output") {
		t.Error("missing Agent Output section")
	}
	if !strings.Contains(result, "agent output here") {
		t.Error("missing agent stdout")
	}
	if !strings.Contains(result, "### Added Files") {
		t.Error("missing Added Files section")
	}
	if !strings.Contains(result, "new.md") {
		t.Error("missing added file name")
	}
	if !strings.Contains(result, "### Modified Files") {
		t.Error("missing Modified Files section")
	}
	if !strings.Contains(result, "### Deleted Files") {
		t.Error("missing Deleted Files section")
	}
	if !strings.Contains(result, "removed.txt") {
		t.Error("missing deleted file name")
	}
}

func TestFormatDiffForJudgeEmpty(t *testing.T) {
	diff := DiffSummary{}
	result := formatDiffForJudge("output", diff)

	if !strings.Contains(result, "## Agent Output") {
		t.Error("missing Agent Output header")
	}
	if strings.Contains(result, "### Added") {
		t.Error("should not contain Added section when empty")
	}
}
