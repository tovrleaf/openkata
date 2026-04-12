package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

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

	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}

func listSkillsTool() mcp.Tool {
	return mcp.NewTool("list_skills",
		mcp.WithDescription("List available OpenKata skills with their descriptions"),
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
		mcp.WithDescription("Install an OpenKata skill into a target project"),
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
		if _, err := os.Stat(filepath.Join(src, "SKILL.md")); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("skill %q not found", skill)), nil
		}

		dest := filepath.Join(targetDir, ".agents", "skills", skill)
		if err := copyDir(src, dest); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to install: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Installed %q to %s", skill, dest)), nil
	}
}

type skillInfo struct {
	Name        string `json:"name"`
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
		desc := extractDescription(string(data))
		skills = append(skills, skillInfo{Name: e.Name(), Description: desc})
	}
	return skills, nil
}

func extractDescription(content string) string {
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
		if !strings.HasPrefix(trimmed, "description:") {
			continue
		}
		value := strings.TrimSpace(strings.TrimPrefix(trimmed, "description:"))
		value = strings.Trim(value, `"'`)

		// Single-line value
		if value != "" && value != ">" && value != "|" {
			return value
		}

		// Multi-line folded/literal scalar: collect indented continuation lines
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

func copyDir(src, dest string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, _ := filepath.Rel(src, path)
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
