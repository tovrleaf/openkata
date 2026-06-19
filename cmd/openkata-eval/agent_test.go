package main

import (
	"strings"
	"testing"
)

func TestBuildAgentPrompt(t *testing.T) {
	skillPath := "testdata/test-skill"
	scenario := Scenario{
		Dir:  "testdata/test-skill/evals/scenario-0",
		Name: "scenario-0",
	}

	prompt, err := buildAgentPrompt(skillPath, scenario)
	if err != nil {
		t.Fatalf("buildAgentPrompt() error: %v", err)
	}

	t.Run("includes SKILL.md", func(t *testing.T) {
		if !strings.Contains(prompt, "Test Skill") {
			t.Error("buildAgentPrompt() missing SKILL.md content")
		}
	})

	t.Run("includes references", func(t *testing.T) {
		if !strings.Contains(prompt, "Reference content here.") {
			t.Error("buildAgentPrompt() missing references/guide.md content")
		}
	})

	t.Run("excludes ACKNOWLEDGMENTS.md", func(t *testing.T) {
		if strings.Contains(prompt, "Should be excluded.") {
			t.Error("buildAgentPrompt() should exclude references/ACKNOWLEDGMENTS.md")
		}
	})

	t.Run("includes scripts", func(t *testing.T) {
		if !strings.Contains(prompt, "echo \"test\"") {
			t.Error("buildAgentPrompt() missing scripts/run.sh content")
		}
	})

	t.Run("includes assets", func(t *testing.T) {
		if !strings.Contains(prompt, "asset data") {
			t.Error("buildAgentPrompt() missing assets/template.txt content")
		}
	})

	t.Run("includes task.md", func(t *testing.T) {
		if !strings.Contains(prompt, "Do the task described above.") {
			t.Error("buildAgentPrompt() missing task.md content")
		}
	})

	t.Run("includes inputs as fenced blocks", func(t *testing.T) {
		if !strings.Contains(prompt, "```src/main.js") {
			t.Error("buildAgentPrompt() missing inputs as fenced code block")
		}
		if !strings.Contains(prompt, "const x = 1;") {
			t.Error("buildAgentPrompt() missing inputs content")
		}
	})

	t.Run("excludes CHANGELOG.md", func(t *testing.T) {
		if strings.Contains(prompt, "Should not appear in prompt.") {
			t.Error("buildAgentPrompt() should exclude CHANGELOG.md")
		}
	})
}

func TestIsExcludedSkillFile(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{"CHANGELOG.md", "CHANGELOG.md", true},
		{"RATIONALE.md", "RATIONALE.md", true},
		{"tile.json", "tile.json", true},
		{".tesslignore", ".tesslignore", true},
		{".tessl-plugin dir", ".tessl-plugin/plugin.json", true},
		{"evals dir", "evals/scenario-0/task.md", true},
		{"SKILL.md", "SKILL.md", false},
		{"references file", "references/guide.md", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isExcludedSkillFile(tt.path)
			if got != tt.want {
				t.Errorf("isExcludedSkillFile(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}
