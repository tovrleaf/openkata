package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
)

// fixedCompleter always returns the same response.
type fixedCompleter struct {
	response string
	calls    atomic.Int32
}

func (f *fixedCompleter) Complete(system, user string) (string, error) {
	f.calls.Add(1)
	return f.response, nil
}

func setupTestScenario(t *testing.T, name string) string {
	t.Helper()
	dir := t.TempDir()

	skillDir := filepath.Join(dir, "test-skill")
	os.MkdirAll(filepath.Join(skillDir, "evals", name), 0755)
	os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("# Test Skill"), 0644)

	scenDir := filepath.Join(skillDir, "evals", name)
	os.WriteFile(filepath.Join(scenDir, "task.md"), []byte("Do the thing"), 0644)

	criteria := CriteriaFile{
		Checklist: []Criterion{
			{Name: "criterion-a", Description: "does A", MaxScore: 10},
			{Name: "criterion-b", Description: "does B", MaxScore: 5},
		},
	}
	data, _ := json.Marshal(criteria)
	os.WriteFile(filepath.Join(scenDir, "criteria.json"), data, 0644)

	return skillDir
}

func TestRunnerOrdering(t *testing.T) {
	skillDir := setupTestScenario(t, "scenario-0")

	scen1Dir := filepath.Join(skillDir, "evals", "scenario-1")
	os.MkdirAll(scen1Dir, 0755)
	os.WriteFile(filepath.Join(scen1Dir, "task.md"), []byte("Do another thing"), 0644)
	criteria := CriteriaFile{
		Checklist: []Criterion{
			{Name: "crit-x", Description: "does X", MaxScore: 8},
		},
	}
	data, _ := json.Marshal(criteria)
	os.WriteFile(filepath.Join(scen1Dir, "criteria.json"), data, 0644)

	judgeResp := `[{"name": "criterion-a", "pass": true, "reason": "good"}, {"name": "criterion-b", "pass": true, "reason": "good"}, {"name": "crit-x", "pass": true, "reason": "ok"}]`
	agent := &fixedCompleter{response: "Here is my detailed response with enough content to pass suspicious check"}
	judge := &fixedCompleter{response: judgeResp}

	runner := &Runner{
		Completer:   agent,
		JudgeModel:  judge,
		Concurrency: 2,
		SkillPath:   skillDir,
	}

	scenarios := []Scenario{
		{Dir: filepath.Join(skillDir, "evals", "scenario-0"), Name: "scenario-0", Description: "First"},
		{Dir: scen1Dir, Name: "scenario-1", Description: "Second"},
	}

	results, err := runner.Run(scenarios)
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("len(results) = %d, want 2", len(results))
	}
	if results[0].Name != "scenario-0" {
		t.Errorf("results[0].Name = %q, want %q", results[0].Name, "scenario-0")
	}
	if results[1].Name != "scenario-1" {
		t.Errorf("results[1].Name = %q, want %q", results[1].Name, "scenario-1")
	}
}

func TestRunnerConcurrencyLimit(t *testing.T) {
	skillDir := setupTestScenario(t, "scenario-0")

	judgeResp := `[{"name": "criterion-a", "pass": true, "reason": "ok"}, {"name": "criterion-b", "pass": true, "reason": "ok"}]`
	agent := &fixedCompleter{response: "Detailed response that passes the suspicious check easily"}
	judge := &fixedCompleter{response: judgeResp}

	runner := &Runner{
		Completer:   agent,
		JudgeModel:  judge,
		Concurrency: 1,
		SkillPath:   skillDir,
	}

	scenarios := []Scenario{
		{Dir: filepath.Join(skillDir, "evals", "scenario-0"), Name: "scenario-0", Description: "Test"},
	}

	results, err := runner.Run(scenarios)
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("len(results) = %d, want 1", len(results))
	}
	if agent.calls.Load() != 1 {
		t.Errorf("agent calls = %d, want 1", agent.calls.Load())
	}
}

func TestRunnerThrottleRetry(t *testing.T) {
	skillDir := setupTestScenario(t, "scenario-0")

	judgeResp := `[{"name": "criterion-a", "pass": true, "reason": "ok"}, {"name": "criterion-b", "pass": true, "reason": "ok"}]`

	var callCount atomic.Int32
	throttleAgent := &throttlingCompleter{
		failCount: 2,
		calls:     &callCount,
		response:  "Detailed response after retry that passes suspicious check",
	}
	judge := &fixedCompleter{response: judgeResp}

	runner := &Runner{
		Completer:   throttleAgent,
		JudgeModel:  judge,
		Concurrency: 1,
		SkillPath:   skillDir,
	}

	scenarios := []Scenario{
		{Dir: filepath.Join(skillDir, "evals", "scenario-0"), Name: "scenario-0", Description: "Retry test"},
	}

	results, err := runner.Run(scenarios)
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("len(results) = %d, want 1", len(results))
	}
	if callCount.Load() != 3 {
		t.Errorf("throttle agent calls = %d, want 3", callCount.Load())
	}
}

// throttlingCompleter fails with 429 for the first N calls, then succeeds.
type throttlingCompleter struct {
	failCount int
	calls     *atomic.Int32
	response  string
}

func (tc *throttlingCompleter) Complete(system, user string) (string, error) {
	n := int(tc.calls.Add(1))
	if n <= tc.failCount {
		return "", fmt.Errorf("HTTP 429: rate limit exceeded")
	}
	return tc.response, nil
}
