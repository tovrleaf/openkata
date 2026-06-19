package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// FileEntry holds a file path and its content.
type FileEntry struct {
	Path    string
	Content string
}

// DiffSummary describes workspace changes.
type DiffSummary struct {
	Added     []FileEntry
	Modified  []FileEntry
	Deleted   []string
	Unchanged []string
}

// snapshotDir reads all files recursively into a map of relative path → content.
func snapshotDir(dir string) (map[string]string, error) {
	files := make(map[string]string)
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(dir, path)
		rel = filepath.ToSlash(rel)
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		files[rel] = string(data)
		return nil
	})
	return files, err
}

// computeDiff produces a DiffSummary from before and after snapshots.
func computeDiff(before, after map[string]string) DiffSummary {
	var diff DiffSummary

	for path, content := range after {
		prev, existed := before[path]
		if !existed {
			diff.Added = append(diff.Added, FileEntry{Path: path, Content: content})
		} else if content != prev {
			diff.Modified = append(diff.Modified, FileEntry{Path: path, Content: content})
		} else {
			diff.Unchanged = append(diff.Unchanged, path)
		}
	}

	for path := range before {
		if _, exists := after[path]; !exists {
			diff.Deleted = append(diff.Deleted, path)
		}
	}

	sort.Slice(diff.Added, func(i, j int) bool { return diff.Added[i].Path < diff.Added[j].Path })
	sort.Slice(diff.Modified, func(i, j int) bool { return diff.Modified[i].Path < diff.Modified[j].Path })
	sort.Strings(diff.Deleted)
	sort.Strings(diff.Unchanged)

	return diff
}

// formatDiffForJudge formats agent stdout and diff for the judge's context.
func formatDiffForJudge(stdout string, diff DiffSummary) string {
	var sb strings.Builder

	sb.WriteString("## Agent Output\n\n")
	sb.WriteString(stdout)
	sb.WriteString("\n\n## Workspace Changes\n\n")

	if len(diff.Added) > 0 {
		sb.WriteString("### Added Files\n\n")
		for _, f := range diff.Added {
			sb.WriteString(fmt.Sprintf("**%s**\n```\n%s\n```\n\n", f.Path, f.Content))
		}
	}

	if len(diff.Modified) > 0 {
		sb.WriteString("### Modified Files\n\n")
		for _, f := range diff.Modified {
			sb.WriteString(fmt.Sprintf("**%s**\n```\n%s\n```\n\n", f.Path, f.Content))
		}
	}

	if len(diff.Deleted) > 0 {
		sb.WriteString("### Deleted Files\n\n")
		for _, p := range diff.Deleted {
			sb.WriteString(fmt.Sprintf("- %s\n", p))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}
