package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestIntegrationNoSandbox(t *testing.T) {
	skillPath := "testdata/test-skill"

	scenarios, err := discoverScenarios(skillPath)
	if err != nil {
		t.Fatalf("discoverScenarios() error: %v", err)
	}
	if len(scenarios) == 0 {
		t.Fatal("no scenarios discovered")
	}

	judgeResp := `[{"name":"criterion-a","pass":true,"reason":"well done"},{"name":"criterion-b","pass":true,"reason":"correct"}]`
	agent := &fixedCompleter{response: "I have completed the task. Created the files as requested with all necessary content to satisfy requirements."}
	judge := &fixedCompleter{response: judgeResp}

	runner := &Runner{
		Completer:   agent,
		JudgeModel:  judge,
		Concurrency: 1,
		SkillPath:   skillPath,
	}

	results, err := runner.Run(scenarios)
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("len(results) = %d, want 1", len(results))
	}

	r := results[0]
	if r.Name != "scenario-0" {
		t.Errorf("Name = %q, want %q", r.Name, "scenario-0")
	}
	if r.Score != 15 {
		t.Errorf("Score = %d, want 15", r.Score)
	}
	if r.MaxScore != 15 {
		t.Errorf("MaxScore = %d, want 15", r.MaxScore)
	}
	if r.Percentage != 100 {
		t.Errorf("Percentage = %f, want 100", r.Percentage)
	}
	if !r.Pass {
		t.Error("Pass = false, want true")
	}

	// Verify overall scoring
	eval := computeOverall(results, 95)
	if !eval.Pass {
		t.Error("overall Pass = false, want true")
	}
	if eval.OverallPercentage != 100 {
		t.Errorf("OverallPercentage = %f, want 100", eval.OverallPercentage)
	}
}

func TestIntegrationPartialFail(t *testing.T) {
	skillPath := "testdata/test-skill"

	scenarios, err := discoverScenarios(skillPath)
	if err != nil {
		t.Fatalf("discoverScenarios() error: %v", err)
	}

	judgeResp := `[{"name":"criterion-a","pass":true,"reason":"good"},{"name":"criterion-b","pass":false,"reason":"missing"}]`
	agent := &fixedCompleter{response: "I completed part of the task with enough content to pass the suspicious check easily."}
	judge := &fixedCompleter{response: judgeResp}

	runner := &Runner{
		Completer:   agent,
		JudgeModel:  judge,
		Concurrency: 1,
		SkillPath:   skillPath,
	}

	results, err := runner.Run(scenarios)
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	r := results[0]
	if r.Score != 10 {
		t.Errorf("Score = %d, want 10", r.Score)
	}
	if r.MaxScore != 15 {
		t.Errorf("MaxScore = %d, want 15", r.MaxScore)
	}
	if r.Pass {
		t.Error("Pass = true, want false (not 100%)")
	}

	eval := computeOverall(results, 95)
	if eval.Pass {
		t.Error("overall Pass = true, want false (66% < 95%)")
	}
}

func TestIntegrationWithSandboxRunner(t *testing.T) {
	skillPath := "testdata/test-skill"

	scenarios, err := discoverScenarios(skillPath)
	if err != nil {
		t.Fatalf("discoverScenarios() error: %v", err)
	}

	// Force sandbox mode
	for i := range scenarios {
		scenarios[i].Sandbox = true
	}

	mockSandbox := func(scenario Scenario, prompt string) (string, error) {
		return "Sandbox completed the task successfully with all required outputs generated.", nil
	}

	judgeResp := `[{"name":"criterion-a","pass":true,"reason":"done"},{"name":"criterion-b","pass":true,"reason":"done"}]`
	judge := &fixedCompleter{response: judgeResp}

	runner := &Runner{
		Completer:   nil, // should not be called
		JudgeModel:  judge,
		Concurrency: 1,
		SkillPath:   skillPath,
		Sandbox:     mockSandbox,
	}

	results, err := runner.Run(scenarios)
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("len(results) = %d, want 1", len(results))
	}
	if !results[0].Pass {
		t.Error("sandbox scenario should pass")
	}
}

func TestIntegrationJSONOutput(t *testing.T) {
	skillPath := "testdata/test-skill"

	scenarios, err := discoverScenarios(skillPath)
	if err != nil {
		t.Fatalf("discoverScenarios() error: %v", err)
	}

	judgeResp := `[{"name":"criterion-a","pass":true,"reason":"ok"},{"name":"criterion-b","pass":true,"reason":"ok"}]`
	agent := &fixedCompleter{response: "Completed the full task with detailed output for all criteria checks."}
	judge := &fixedCompleter{response: judgeResp}

	runner := &Runner{
		Completer:   agent,
		JudgeModel:  judge,
		Concurrency: 1,
		SkillPath:   skillPath,
	}

	results, err := runner.Run(scenarios)
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	eval := computeOverall(results, 95)

	// Write JSON output
	outDir := t.TempDir()
	outPath := filepath.Join(outDir, "results.json")
	cfg := Config{
		Backend:    "kiro",
		Model:      "test-model",
		JudgeModel: "test-judge",
	}

	if err := writeJSONOutput(outPath, "test-skill", "", cfg, eval, false); err != nil {
		t.Fatalf("writeJSONOutput() error: %v", err)
	}

	// Read and verify structure
	files, _ := filepath.Glob(filepath.Join(outDir, "*.json"))
	if len(files) == 0 {
		t.Fatal("no JSON output file created")
	}

	data, err := os.ReadFile(files[0])
	if err != nil {
		t.Fatalf("reading output: %v", err)
	}

	var out JSONOutput
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("parsing output JSON: %v", err)
	}

	if out.Skill != "test-skill" {
		t.Errorf("Skill = %q, want %q", out.Skill, "test-skill")
	}
	if !out.Pass {
		t.Error("Pass = false, want true")
	}
	if out.OverallPercentage != 100 {
		t.Errorf("OverallPercentage = %f, want 100", out.OverallPercentage)
	}
	if len(out.Scenarios) != 1 {
		t.Fatalf("len(Scenarios) = %d, want 1", len(out.Scenarios))
	}
	if len(out.Scenarios[0].Criteria) != 2 {
		t.Errorf("len(Criteria) = %d, want 2", len(out.Scenarios[0].Criteria))
	}
	if out.Config.Backend != "kiro" {
		t.Errorf("Config.Backend = %q, want %q", out.Config.Backend, "kiro")
	}
}

func TestIntegrationExitCode(t *testing.T) {
	tests := []struct {
		name      string
		mode      parseMode
		threshold int
		pct       float64
		wantExit  int
	}{
		{"skill pass", modeSkill, 95, 100, 0},
		{"skill fail", modeSkill, 95, 60, 1},
		{"scenario always 0", modeScenario, 95, 60, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eval := EvalResult{
				OverallPercentage: tt.pct,
				Threshold:         tt.threshold,
				Pass:              tt.pct >= float64(tt.threshold),
			}
			got := shouldExit(tt.mode, eval)
			if got != tt.wantExit {
				t.Errorf("shouldExit() = %d, want %d", got, tt.wantExit)
			}
		})
	}
}
