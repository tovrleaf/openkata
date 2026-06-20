package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ANSI color codes.
const (
	colorGreen = "\033[32m"
	colorRed   = "\033[31m"
	colorReset = "\033[0m"
)

// printScenarioResult prints a single scenario result to stdout with colors.
func printScenarioResult(r ScenarioResult) {
	fmt.Printf("  %s: %s\n", r.Name, r.Description)
	for _, c := range r.Criteria {
		if c.Pass {
			fmt.Printf("    %s✓%s %s (%d/%d)\n", colorGreen, colorReset, c.Name, c.Score, c.MaxScore)
		} else {
			reason := ""
			if c.Reason != "" {
				reason = " — " + c.Reason
			}
			fmt.Printf("    %s✗%s %s (%d/%d)%s\n", colorRed, colorReset, c.Name, c.Score, c.MaxScore, reason)
		}
	}
	passLabel := fmt.Sprintf("%sPASS%s", colorGreen, colorReset)
	if !r.Pass {
		passLabel = fmt.Sprintf("%sFAIL%s", colorRed, colorReset)
	}
	fmt.Printf("    Score: %d/%d (%.0f%%) %s\n\n", r.Score, r.MaxScore, r.Percentage, passLabel)
}

// printOverall prints the overall result line.
func printOverall(result EvalResult) {
	passLabel := fmt.Sprintf("%sPASS%s", colorGreen, colorReset)
	if !result.Pass {
		passLabel = fmt.Sprintf("%sFAIL%s", colorRed, colorReset)
	}
	fmt.Printf("  Overall: %.0f%% %s (threshold: %d%%)\n", result.OverallPercentage, passLabel, result.Threshold)
}

// JSONOutput is the JSON report structure.
type JSONOutput struct {
	Skill             string           `json:"skill"`
	Timestamp         string           `json:"timestamp"`
	Config            JSONOutputConfig  `json:"config"`
	Threshold         int              `json:"threshold"`
	Scenarios         []JSONScenario   `json:"scenarios"`
	OverallPercentage float64          `json:"overall_percentage"`
	Pass              bool             `json:"pass"`
}

// JSONOutputConfig holds the config section of JSON output.
type JSONOutputConfig struct {
	Backend    string `json:"backend"`
	AgentModel string `json:"agent_model"`
	JudgeModel string `json:"judge_model"`
}

// JSONScenario is a scenario in JSON output.
type JSONScenario struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Criteria    []JSONCriterion `json:"criteria"`
	Score       int             `json:"score"`
	MaxScore    int             `json:"max_score"`
	Percentage  float64         `json:"percentage"`
	Pass        bool            `json:"pass"`
}

// JSONCriterion is a criterion in JSON output.
type JSONCriterion struct {
	Name     string `json:"name"`
	Pass     bool   `json:"pass"`
	Score    int    `json:"score"`
	MaxScore int    `json:"max_score"`
	Reason   string `json:"reason"`
}

// writeJSONOutput writes the eval result to a JSON file.
func writeJSONOutput(outputPath string, skill string, cfg Config, result EvalResult) error {
	now := time.Now().UTC()

	out := JSONOutput{
		Skill:     skill,
		Timestamp: now.Format(time.RFC3339),
		Config: JSONOutputConfig{
			Backend:    cfg.Backend,
			AgentModel: cfg.Model,
			JudgeModel: cfg.JudgeModel,
		},
		Threshold:         result.Threshold,
		OverallPercentage: result.OverallPercentage,
		Pass:              result.Pass,
	}

	for _, s := range result.Scenarios {
		js := JSONScenario{
			Name:        s.Name,
			Description: s.Description,
			Score:       s.Score,
			MaxScore:    s.MaxScore,
			Percentage:  s.Percentage,
			Pass:        s.Pass,
		}
		for _, c := range s.Criteria {
			js.Criteria = append(js.Criteria, JSONCriterion{
				Name:     c.Name,
				Pass:     c.Pass,
				Score:    c.Score,
				MaxScore: c.MaxScore,
				Reason:   c.Reason,
			})
		}
		out.Scenarios = append(out.Scenarios, js)
	}

	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return err
	}

	filePath := resolveOutputPath(outputPath, now)
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}
	return os.WriteFile(filePath, append(data, '\n'), 0644)
}

// resolveOutputPath determines the final file path for JSON output.
// If outputPath is a directory, generates a timestamped filename inside it.
// Otherwise appends timestamp before the extension.
func resolveOutputPath(outputPath string, now time.Time) string {
	ts := now.Format("2006-01-02T150405")

	info, err := os.Stat(outputPath)
	if err == nil && info.IsDir() {
		return filepath.Join(outputPath, "results-"+ts+".json")
	}

	ext := filepath.Ext(outputPath)
	base := strings.TrimSuffix(outputPath, ext)
	if ext == "" {
		ext = ".json"
	}
	return base + "-" + ts + ext
}
