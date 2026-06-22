package exclude

import "strings"

// FromInstall returns true if the file should not be included
// in archive downloads, MCP installs, or eval agent copies.
func FromInstall(path string) bool {
	switch path {
	case "CHANGELOG.md", "RATIONALE.md", "tile.json", ".tesslignore", "references/ACKNOWLEDGMENTS.md":
		return true
	}
	if strings.HasPrefix(path, "evals/") || strings.HasPrefix(path, ".tessl-plugin/") {
		return true
	}
	return false
}

// FromFilesTab returns true if the file should not appear in
// the web UI Files tab. Superset of FromInstall.
func FromFilesTab(path string) bool {
	return FromInstall(path)
}
