package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Criterion represents a single eval criterion from criteria.json.
type Criterion struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	MaxScore    int    `json:"max_score"`
}

// CriteriaFile represents the criteria.json file.
type CriteriaFile struct {
	Context   string      `json:"context"`
	Type      string      `json:"type"`
	Checklist []Criterion `json:"checklist"`
}

// Verdict is the judge's result for one criterion.
type Verdict struct {
	Name   string `json:"name"`
	Pass   bool   `json:"pass"`
	Reason string `json:"reason"`
}

// judgeSystemPrompt is the system prompt for the judge model.
const judgeSystemPrompt = `You are an eval judge. You evaluate an AI agent's response against a checklist of criteria.

For each criterion, determine independently whether the agent's response satisfies it.

Return ONLY a JSON array of objects with these fields:
- "name": the criterion name (exact match)
- "pass": true if satisfied, false otherwise
- "reason": brief explanation

Example output:
[{"name": "criterion-a", "pass": true, "reason": "clearly demonstrated"}, {"name": "criterion-b", "pass": false, "reason": "not addressed"}]

Do not include any text outside the JSON array.`

// loadCriteria reads and parses criteria.json from a scenario directory.
func loadCriteria(scenarioDir string) (*CriteriaFile, error) {
	data, err := os.ReadFile(filepath.Join(scenarioDir, "criteria.json"))
	if err != nil {
		return nil, fmt.Errorf("reading criteria.json: %w", err)
	}
	var cf CriteriaFile
	if err := json.Unmarshal(data, &cf); err != nil {
		return nil, fmt.Errorf("parsing criteria.json: %w", err)
	}
	return &cf, nil
}

// buildJudgeUserMessage constructs the user message for the judge call.
func buildJudgeUserMessage(agentResponse string, criteria *CriteriaFile) string {
	var sb strings.Builder
	sb.WriteString("## Agent Response\n\n")
	sb.WriteString(agentResponse)
	sb.WriteString("\n\n## Evaluation Criteria\n\n")
	if criteria.Context != "" {
		sb.WriteString("Context: ")
		sb.WriteString(criteria.Context)
		sb.WriteString("\n\n")
	}
	for _, c := range criteria.Checklist {
		sb.WriteString(fmt.Sprintf("- **%s** (max %d points): %s\n", c.Name, c.MaxScore, c.Description))
	}
	return sb.String()
}

// judgeEval sends the agent response to the judge and parses verdicts.
// Retries once on parse failure.
func judgeEval(completer Completer, agentResponse string, criteria *CriteriaFile) ([]Verdict, error) {
	userMsg := buildJudgeUserMessage(agentResponse, criteria)

	verdicts, err := callJudge(completer, userMsg)
	if err != nil {
		// Retry once
		verdicts, err = callJudge(completer, userMsg)
		if err != nil {
			return nil, fmt.Errorf("judge parse failed after retry: %w", err)
		}
	}

	// Sanity check: empty/trivial agent response + all pass = suspicious
	if isSuspicious(agentResponse, verdicts) {
		return nil, fmt.Errorf("suspicious result: trivial agent response (%d chars) with all criteria passing", len(agentResponse))
	}

	return verdicts, nil
}

// callJudge makes a single judge completion call and parses the JSON response.
func callJudge(completer Completer, userMsg string) ([]Verdict, error) {
	raw, err := completer.Complete(judgeSystemPrompt, userMsg)
	if err != nil {
		return nil, err
	}
	return parseJudgeResponse(raw)
}

// parseJudgeResponse extracts a JSON array of verdicts from raw judge output.
func parseJudgeResponse(raw string) ([]Verdict, error) {
	// Strip ANSI and header lines
	cleaned := cleanKiroOutput(raw)

	// Find first [ or {
	start := -1
	for i, c := range cleaned {
		if c == '[' || c == '{' {
			start = i
			break
		}
	}
	if start == -1 {
		return nil, fmt.Errorf("no JSON found in judge response")
	}

	jsonStr := cleaned[start:]

	// Try parsing as array first
	var verdicts []Verdict
	if err := json.Unmarshal([]byte(jsonStr), &verdicts); err == nil {
		return verdicts, nil
	}

	// Try trimming after last ] or }
	if end := strings.LastIndex(jsonStr, "]"); end >= 0 {
		var v []Verdict
		if err := json.Unmarshal([]byte(jsonStr[:end+1]), &v); err == nil {
			return v, nil
		}
	}

	return nil, fmt.Errorf("failed to parse judge JSON")
}

// isSuspicious returns true if the agent response looks trivial but all criteria pass.
func isSuspicious(agentResponse string, verdicts []Verdict) bool {
	if len(strings.TrimSpace(agentResponse)) >= 20 {
		return false
	}
	for _, v := range verdicts {
		if !v.Pass {
			return false
		}
	}
	return len(verdicts) > 0
}
