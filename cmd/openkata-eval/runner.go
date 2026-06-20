package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// SandboxRunner is a function type for running sandbox scenarios (stubbed for now).
type SandboxRunner func(scenario Scenario, prompt string) (string, error)

// Runner orchestrates scenario execution.
type Runner struct {
	Completer   Completer
	JudgeModel  Completer
	Concurrency int
	SkillPath   string
	Sandbox     SandboxRunner
}

// scenarioOutcome holds the result and index for ordering.
type scenarioOutcome struct {
	index  int
	result ScenarioResult
	err    error
}

// Run executes all scenarios with bounded concurrency, printing results in order.
func (r *Runner) Run(scenarios []Scenario) ([]ScenarioResult, error) {
	outcomes := make([]scenarioOutcome, len(scenarios))
	sem := make(chan struct{}, r.Concurrency)
	var wg sync.WaitGroup

	for i, s := range scenarios {
		wg.Add(1)
		go func(idx int, sc Scenario) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			result, err := r.runScenario(sc)
			outcomes[idx] = scenarioOutcome{index: idx, result: result, err: err}
		}(i, s)
	}
	wg.Wait()

	// Collect results in order, print each
	var results []ScenarioResult
	for _, o := range outcomes {
		if o.err != nil {
			// Create a failed result for errors
			results = append(results, ScenarioResult{
				Name:        scenarios[o.index].Name,
				Description: scenarios[o.index].Description + " (error: " + o.err.Error() + ")",
			})
		} else {
			results = append(results, o.result)
		}
		printScenarioResult(results[len(results)-1])
	}
	return results, nil
}

// runScenario executes a single scenario with retry on throttling.
func (r *Runner) runScenario(s Scenario) (ScenarioResult, error) {
	criteria, err := loadCriteria(s.Dir)
	if err != nil {
		return ScenarioResult{}, err
	}

	prompt, err := buildAgentPrompt(r.SkillPath, s)
	if err != nil {
		return ScenarioResult{}, err
	}

	var agentResponse string
	if s.Sandbox && r.Sandbox != nil {
		agentResponse, err = r.Sandbox(s, prompt)
	} else {
		agentResponse, err = r.completeWithRetry(prompt)
	}
	if err != nil {
		return ScenarioResult{}, fmt.Errorf("agent completion: %w", err)
	}

	verdicts, err := judgeEval(r.JudgeModel, agentResponse, criteria)
	if err != nil {
		return ScenarioResult{}, fmt.Errorf("judge: %w", err)
	}

	return scoreScenario(s.Name, s.Description, verdicts, criteria), nil
}

// completeWithRetry retries on throttling errors with exponential backoff.
func (r *Runner) completeWithRetry(prompt string) (string, error) {
	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		resp, err := r.Completer.Complete("", prompt)
		if err == nil {
			return resp, nil
		}
		if !isThrottleError(err) {
			return "", err
		}
		lastErr = err
		time.Sleep(time.Duration(1<<uint(attempt)) * time.Second)
	}
	return "", fmt.Errorf("throttled after 3 attempts: %w", lastErr)
}

// isThrottleError checks if an error indicates rate limiting.
func isThrottleError(err error) bool {
	msg := err.Error()
	return strings.Contains(msg, "429") || strings.Contains(msg, "throttl") || strings.Contains(msg, "rate limit")
}
