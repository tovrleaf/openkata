package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const source = "github.com/tovrleaf/openkata"

var skipFiles = map[string]bool{
	"CHANGELOG.md": true,
}

func main() {
	skillsDir := os.Getenv("OPENKATA_SKILLS_DIR")
	if skillsDir == "" {
		skillsDir = "skills"
	}

	s := server.NewMCPServer(
		"openkata",
		"0.1.0",
		server.WithToolCapabilities(false),
	)

	s.AddTool(listSkillsTool(), listSkillsHandler(skillsDir))
	s.AddTool(installSkillTool(), installSkillHandler(skillsDir))

	addr := os.Getenv("OPENKATA_ADDR")
	if addr != "" {
		fmt.Fprintf(os.Stderr, "OpenKata MCP server listening on %s\n", addr)
		httpServer := server.NewStreamableHTTPServer(s)
		if err := httpServer.Start(addr); err != nil {
			fmt.Fprintf(os.Stderr, "server error: %v\n", err)
			os.Exit(1)
		}
	} else {
		if err := server.ServeStdio(s); err != nil {
			fmt.Fprintf(os.Stderr, "server error: %v\n", err)
			os.Exit(1)
		}
	}
}

// --- Tools ---

func listSkillsTool() mcp.Tool {
	return mcp.NewTool("list_skills",
		mcp.WithDescription("List available OpenKata skills with their descriptions and versions"),
	)
}

func listSkillsHandler(skillsDir string) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		skills, err := discoverSkills(skillsDir)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to list skills: %v", err)), nil
		}
		out, _ := json.MarshalIndent(skills, "", "  ")
		return mcp.NewToolResultText(string(out)), nil
	}
}

func installSkillTool() mcp.Tool {
	return mcp.NewTool("install_skill",
		mcp.WithDescription("Install an OpenKata skill into a target project. Copies skill files and writes a .manifest.json with version and origin."),
		mcp.WithString("skill",
			mcp.Required(),
			mcp.Description("Name of the skill to install"),
		),
		mcp.WithString("target_dir",
			mcp.Required(),
			mcp.Description("Absolute path to the target project root"),
		),
	)
}

func installSkillHandler(skillsDir string) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		skill, err := req.RequireString("skill")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		targetDir, err := req.RequireString("target_dir")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		src := filepath.Join(skillsDir, skill)
		skillMD := filepath.Join(src, "SKILL.md")
		data, err := os.ReadFile(skillMD)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("skill %q not found", skill)), nil
		}

		dest := filepath.Join(targetDir, ".agents", "skills", skill)
		if err := copyDir(src, dest); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to install: %v", err)), nil
		}

		version := extractFrontmatterField(string(data), "version")
		manifest := manifest{
			Name:        skill,
			Version:     version,
			Source:      source,
			InstalledAt: time.Now().UTC().Format(time.RFC3339),
		}
		if err := writeManifest(dest, manifest); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("installed but failed to write manifest: %v", err)), nil
		}

		msg := fmt.Sprintf("Installed %q v%s to %s", skill, version, dest)
		return mcp.NewToolResultText(msg), nil
	}
}

// --- Manifest ---

type manifest struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Source      string `json:"source"`
	InstalledAt string `json:"installedAt"`
}

func writeManifest(skillDir string, m manifest) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(skillDir, ".manifest.json"), append(data, '\n'), 0644)
}

// --- Discovery ---

type skillInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

func discoverSkills(skillsDir string) ([]skillInfo, error) {
	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		return nil, err
	}

	var skills []skillInfo
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		data, err := os.ReadFile(filepath.Join(skillsDir, e.Name(), "SKILL.md"))
		if err != nil {
			continue
		}
		content := string(data)
		skills = append(skills, skillInfo{
			Name:        e.Name(),
			Version:     extractFrontmatterField(content, "version"),
			Description: extractDescription(content),
		})
	}
	return skills, nil
}

// --- Frontmatter parsing ---

func extractFrontmatterField(content, field string) string {
	if !strings.HasPrefix(content, "---") {
		return ""
	}
	end := strings.Index(content[3:], "---")
	if end < 0 {
		return ""
	}
	frontmatter := content[3 : end+3]
	lines := strings.Split(frontmatter, "\n")

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		prefix := field + ":"
		if !strings.HasPrefix(trimmed, prefix) {
			continue
		}
		value := strings.TrimSpace(strings.TrimPrefix(trimmed, prefix))
		value = strings.Trim(value, `"'`)

		if value != "" && value != ">" && value != "|" {
			return value
		}

		// Multi-line folded/literal scalar
		var parts []string
		for j := i + 1; j < len(lines); j++ {
			l := lines[j]
			if len(l) == 0 || l[0] == ' ' || l[0] == '\t' {
				parts = append(parts, strings.TrimSpace(l))
			} else {
				break
			}
		}
		return strings.Join(parts, " ")
	}
	return ""
}

func extractDescription(content string) string {
	return extractFrontmatterField(content, "description")
}

// --- File copying ---

func copyDir(src, dest string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, _ := filepath.Rel(src, path)

		if skipFiles[d.Name()] {
			return nil
		}

		target := filepath.Join(dest, rel)

		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0644)
	})
}
