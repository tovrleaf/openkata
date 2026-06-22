package templates

import (
	"encoding/json"
	"fmt"
	"path"
	"sort"
	"strings"
)

// SplitTags splits a comma-separated tags string into trimmed slices.
func SplitTags(tags string) []string {
	if tags == "" {
		return nil
	}
	parts := strings.Split(tags, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}

// TagClass returns a CSS class based on the tag prefix.
func TagClass(tag string) string {
	prefix, _, _ := strings.Cut(tag, ":")
	switch prefix {
	case "tool":
		return "badge badge-orange"
	case "language":
		return "badge badge-purple"
	default:
		return "badge badge-green"
	}
}

// FirstSentence returns text up to and including the first period.
func FirstSentence(s string) string {
	if i := strings.Index(s, "."); i >= 0 {
		return s[:i+1]
	}
	return s
}

// AfterFirstSentence returns text after the first period.
func AfterFirstSentence(s string) string {
	if i := strings.Index(s, "."); i >= 0 {
		return strings.TrimSpace(s[i+1:])
	}
	return ""
}

// fileContentsJSON serializes file contents map to JSON for embedding.
func fileContentsJSON(contents map[string]string) string {
	if contents == nil {
		return "{}"
	}
	b, err := json.Marshal(contents)
	if err != nil {
		return "{}"
	}
	return string(b)
}

// treeNode represents a file or directory in the tree.
type treeNode struct {
	name     string
	path     string
	children []*treeNode
	isDir    bool
}

// buildFileTree generates HTML for a file tree with connectors.
func buildFileTree(files []string, artifactType, version string) string {
	_ = artifactType
	_ = version
	root := &treeNode{isDir: true}
	for _, f := range files {
		parts := strings.Split(f, "/")
		insertNode(root, parts, f)
	}
	sortTree(root)
	var sb strings.Builder
	for i, child := range root.children {
		isLast := i == len(root.children)-1
		renderNode(&sb, child, "", isLast)
	}
	return sb.String()
}

func insertNode(parent *treeNode, parts []string, fullPath string) {
	if len(parts) == 0 {
		return
	}
	name := parts[0]
	var child *treeNode
	for _, c := range parent.children {
		if c.name == name {
			child = c
			break
		}
	}
	if child == nil {
		child = &treeNode{name: name}
		if len(parts) > 1 {
			child.isDir = true
			child.path = path.Join(path.Dir(fullPath), name)
		} else {
			child.path = fullPath
		}
		parent.children = append(parent.children, child)
	}
	if len(parts) > 1 {
		child.isDir = true
		insertNode(child, parts[1:], fullPath)
	}
}

func sortTree(node *treeNode) {
	sort.Slice(node.children, func(i, j int) bool {
		// Directories first, then alphabetical
		if node.children[i].isDir != node.children[j].isDir {
			return node.children[i].isDir
		}
		return node.children[i].name < node.children[j].name
	})
	for _, c := range node.children {
		if c.isDir {
			sortTree(c)
		}
	}
}

// fileID converts a filepath to a valid HTML id.
func fileID(path string) string {
	id := strings.ReplaceAll(path, "/", "-")
	id = strings.ReplaceAll(id, ".", "-")
	return id
}

func renderNode(sb *strings.Builder, node *treeNode, prefix string, isLast bool) {
	connector := "├── "
	if isLast {
		connector = "└── "
	}
	if node.isDir {
		sb.WriteString(fmt.Sprintf(`<details class="file-tree-dir"><summary><span class="file-tree-prefix">%s%s</span><span class="file-tree-dirname">%s/</span></summary>`, prefix, connector, node.name))
		childPrefix := prefix + "│   "
		if isLast {
			childPrefix = prefix + "    "
		}
		for i, child := range node.children {
			childIsLast := i == len(node.children)-1
			renderNode(sb, child, childPrefix, childIsLast)
		}
		sb.WriteString("</details>")
	} else {
		sb.WriteString(fmt.Sprintf(`<div class="file-tree-file"><span class="file-tree-prefix">%s%s</span><a class="file-tree-link" href="#file-%s">%s</a></div>`, prefix, connector, fileID(node.path), node.name))
	}
}

// SkillTabs returns the tab definitions for a skill detail page.
func SkillTabs(skill ArtifactDetail) []TabDef {
	tabs := []TabDef{{"overview", "Overview"}, {"files", "Files"}}
	if skill.Rationale != "" {
		tabs = append(tabs, TabDef{"rationale", "Rationale"})
	}
	tabs = append(tabs, TabDef{"changelog", "Changelog"})
	if len(skill.Models) > 0 || skill.TesslScore > 0 {
		tabs = append(tabs, TabDef{"benchmarks", "Benchmarks"})
	}
	if skill.Acknowledgments != "" {
		tabs = append(tabs, TabDef{"acknowledgments", "Acknowledgments"})
	}
	return tabs
}

// ScenarioBarClass returns CSS classes for the scenario color bar.
func ScenarioBarClass(score int) string {
	switch {
	case score >= 95:
		return "scenario-bar scenario-bar--green"
	case score >= 70:
		return "scenario-bar scenario-bar--yellow"
	default:
		return "scenario-bar scenario-bar--red"
	}
}

// CountPass returns the number of passing scenarios.
func CountPass(scenarios []BenchmarkScenario) int {
	n := 0
	for _, s := range scenarios {
		if s.Pass {
			n++
		}
	}
	return n
}