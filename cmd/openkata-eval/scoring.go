package main

// ScenarioResult holds the eval result for a single scenario.
type ScenarioResult struct {
	Name        string
	Description string
	Criteria    []CriterionResult
	Score       int
	MaxScore    int
	Percentage  float64
	Pass        bool
}

// CriterionResult holds the result for a single criterion.
type CriterionResult struct {
	Name     string
	Pass     bool
	Score    int
	MaxScore int
	Reason   string
}

// EvalResult holds the overall eval results.
type EvalResult struct {
	Scenarios         []ScenarioResult
	OverallPercentage float64
	Threshold         int
	Pass              bool
}

// scoreScenario computes the score for a single scenario.
func scoreScenario(name, description string, verdicts []Verdict, criteria *CriteriaFile) ScenarioResult {
	result := ScenarioResult{
		Name:        name,
		Description: description,
	}

	for _, c := range criteria.Checklist {
		cr := CriterionResult{
			Name:     c.Name,
			MaxScore: c.MaxScore,
		}
		// Find matching verdict
		for _, v := range verdicts {
			if v.Name == c.Name {
				cr.Pass = v.Pass
				cr.Reason = v.Reason
				if v.Pass {
					cr.Score = c.MaxScore
				}
				break
			}
		}
		result.Criteria = append(result.Criteria, cr)
		result.Score += cr.Score
		result.MaxScore += cr.MaxScore
	}

	if result.MaxScore > 0 {
		result.Percentage = float64(result.Score) / float64(result.MaxScore) * 100
	}
	result.Pass = result.Percentage >= 100 // individual scenarios pass at 100%

	return result
}

// computeOverall computes the overall eval result from scenario results.
func computeOverall(scenarios []ScenarioResult, threshold int) EvalResult {
	result := EvalResult{
		Scenarios: scenarios,
		Threshold: threshold,
	}

	if len(scenarios) == 0 {
		return result
	}

	var totalPct float64
	for _, s := range scenarios {
		totalPct += s.Percentage
	}
	result.OverallPercentage = totalPct / float64(len(scenarios))
	result.Pass = result.OverallPercentage >= float64(threshold)

	return result
}

// shouldExit returns the exit code based on mode and results.
// Single scenario mode always returns 0.
// Full skill mode returns 1 if below threshold.
func shouldExit(mode parseMode, result EvalResult) int {
	if mode == modeScenario {
		return 0
	}
	if !result.Pass {
		return 1
	}
	return 0
}
