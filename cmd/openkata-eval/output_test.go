package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestWriteJSONOutput(t *testing.T) {
	dir := t.TempDir()
	outPath := filepath.Join(dir, "results.json")

	cfg := Config{
		Backend:    "http",
		Model:      "gpt-4o",
		JudgeModel: "gpt-4o",
	}
	result := EvalResult{
		Scenarios: []ScenarioResult{
			{
				Name:        "scenario-0",
				Description: "Test scenario",
				Criteria: []CriterionResult{
					{Name: "crit-a", Pass: true, Score: 10, MaxScore: 10},
					{Name: "crit-b", Pass: false, Score: 0, MaxScore: 5, Reason: "missing"},
				},
				Score:      10,
				MaxScore:   15,
				Percentage: 66.7,
				Pass:       false,
			},
		},
		OverallPercentage: 66.7,
		Threshold:         95,
		Pass:              false,
	}

	err := writeJSONOutput(outPath, "commit-conventions", "", cfg, result, false)
	if err != nil {
		t.Fatalf("writeJSONOutput() error: %v", err)
	}

	// Find the written file (has timestamp appended)
	entries, _ := os.ReadDir(dir)
	if len(entries) == 0 {
		t.Fatal("no output file written")
	}

	data, err := os.ReadFile(filepath.Join(dir, entries[0].Name()))
	if err != nil {
		t.Fatalf("reading output: %v", err)
	}

	var out JSONOutput
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("parsing output JSON: %v", err)
	}

	if out.Skill != "commit-conventions" {
		t.Errorf("Skill = %q, want %q", out.Skill, "commit-conventions")
	}
	if out.Config.Backend != "http" {
		t.Errorf("Config.Backend = %q, want %q", out.Config.Backend, "http")
	}
	if len(out.Scenarios) != 1 {
		t.Fatalf("len(Scenarios) = %d, want 1", len(out.Scenarios))
	}
	if out.Scenarios[0].Score != 10 {
		t.Errorf("Scenarios[0].Score = %d, want 10", out.Scenarios[0].Score)
	}
	if out.Pass {
		t.Error("Pass = true, want false")
	}
}

func TestWriteJSONOutputDirectory(t *testing.T) {
	dir := t.TempDir()

	cfg := Config{Backend: "kiro", Model: "test", JudgeModel: "test"}
	result := EvalResult{Threshold: 95}

	err := writeJSONOutput(dir, "test-skill", "", cfg, result, false)
	if err != nil {
		t.Fatalf("writeJSONOutput() error: %v", err)
	}

	entries, _ := os.ReadDir(dir)
	if len(entries) != 1 {
		t.Fatalf("expected 1 file, got %d", len(entries))
	}
	if !strings.HasPrefix(entries[0].Name(), "results-") {
		t.Errorf("filename %q doesn't start with 'results-'", entries[0].Name())
	}
}

func TestResolveOutputPath(t *testing.T) {
	now := time.Date(2026, 6, 18, 11, 0, 0, 0, time.UTC)

	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "file with extension",
			path: "/tmp/nonexistent/results.json",
			want: "/tmp/nonexistent/results-2026-06-18T110000.json",
		},
		{
			name: "file without extension",
			path: "/tmp/nonexistent/output",
			want: "/tmp/nonexistent/output-2026-06-18T110000.json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolveOutputPath(tt.path, now)
			if got != tt.want {
				t.Errorf("resolveOutputPath(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}

	// Test directory case
	t.Run("directory", func(t *testing.T) {
		dir := t.TempDir()
		got := resolveOutputPath(dir, now)
		want := filepath.Join(dir, "results-2026-06-18T110000.json")
		if got != want {
			t.Errorf("resolveOutputPath(dir) = %q, want %q", got, want)
		}
	})
}

func TestParseSkillVersion(t *testing.T) {
	tests := []struct {
		name      string
		changelog string
		want      string
	}{
		{
			name: "version found",
			changelog: "# Changelog\n\n## [Unreleased]\n\n## [1.2.3] - 2026-06-01\n\n### Added\n- stuff\n",
			want: "1.2.3",
		},
		{
			name: "no unreleased, direct version",
			changelog: "## [0.5.0] - 2026-01-01\n",
			want: "0.5.0",
		},
		{
			name:      "no changelog file",
			changelog: "",
			want:      "",
		},
		{
			name:      "no version headings",
			changelog: "# Changelog\n\nNo versions yet.\n",
			want:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			if tt.changelog != "" {
				if err := os.WriteFile(filepath.Join(dir, "CHANGELOG.md"), []byte(tt.changelog), 0644); err != nil {
					t.Fatal(err)
				}
			}
			got := parseSkillVersion(dir)
			if got != tt.want {
				t.Errorf("parseSkillVersion(%q) = %q, want %q", tt.changelog, got, tt.want)
			}
		})
	}
}

func TestAutoResolvePath(t *testing.T) {
	now := time.Date(2026, 6, 21, 19, 43, 0, 0, time.UTC)
	got := autoResolvePath("skills/commit-conventions", "claude-sonnet-4.6", now)
	want := filepath.Join("skills", "commit-conventions", "evals", "results", "claude-sonnet-4.6", "2026-06-21T194300.json")
	if got != want {
		t.Errorf("autoResolvePath() = %q, want %q", got, want)
	}
}
