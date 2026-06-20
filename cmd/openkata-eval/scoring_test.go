package main

import (
	"math"
	"testing"
)

func TestScoreScenario(t *testing.T) {
	tests := []struct {
		name       string
		verdicts   []Verdict
		criteria   *CriteriaFile
		wantScore  int
		wantMax    int
		wantPct    float64
		wantPass   bool
	}{
		{
			name: "all pass",
			verdicts: []Verdict{
				{Name: "a", Pass: true},
				{Name: "b", Pass: true},
			},
			criteria: &CriteriaFile{
				Checklist: []Criterion{
					{Name: "a", MaxScore: 10},
					{Name: "b", MaxScore: 5},
				},
			},
			wantScore: 15,
			wantMax:   15,
			wantPct:   100,
			wantPass:  true,
		},
		{
			name: "partial pass",
			verdicts: []Verdict{
				{Name: "a", Pass: true},
				{Name: "b", Pass: false, Reason: "missing"},
			},
			criteria: &CriteriaFile{
				Checklist: []Criterion{
					{Name: "a", MaxScore: 10},
					{Name: "b", MaxScore: 10},
				},
			},
			wantScore: 10,
			wantMax:   20,
			wantPct:   50,
			wantPass:  false,
		},
		{
			name:     "all fail",
			verdicts: []Verdict{{Name: "a", Pass: false}},
			criteria: &CriteriaFile{
				Checklist: []Criterion{{Name: "a", MaxScore: 10}},
			},
			wantScore: 0,
			wantMax:   10,
			wantPct:   0,
			wantPass:  false,
		},
		{
			name:     "missing verdict defaults to fail",
			verdicts: []Verdict{},
			criteria: &CriteriaFile{
				Checklist: []Criterion{{Name: "a", MaxScore: 10}},
			},
			wantScore: 0,
			wantMax:   10,
			wantPct:   0,
			wantPass:  false,
		},
		{
			name:     "empty criteria",
			verdicts: []Verdict{},
			criteria: &CriteriaFile{
				Checklist: []Criterion{},
			},
			wantScore: 0,
			wantMax:   0,
			wantPct:   0,
			wantPass:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := scoreScenario("test", "desc", tt.verdicts, tt.criteria)
			if result.Score != tt.wantScore {
				t.Errorf("scoreScenario() Score = %d, want %d", result.Score, tt.wantScore)
			}
			if result.MaxScore != tt.wantMax {
				t.Errorf("scoreScenario() MaxScore = %d, want %d", result.MaxScore, tt.wantMax)
			}
			if math.Abs(result.Percentage-tt.wantPct) > 0.01 {
				t.Errorf("scoreScenario() Percentage = %f, want %f", result.Percentage, tt.wantPct)
			}
			if result.Pass != tt.wantPass {
				t.Errorf("scoreScenario() Pass = %v, want %v", result.Pass, tt.wantPass)
			}
		})
	}
}

func TestComputeOverall(t *testing.T) {
	tests := []struct {
		name      string
		scenarios []ScenarioResult
		threshold int
		wantPct   float64
		wantPass  bool
	}{
		{
			name: "above threshold",
			scenarios: []ScenarioResult{
				{Percentage: 100},
				{Percentage: 96},
			},
			threshold: 95,
			wantPct:   98,
			wantPass:  true,
		},
		{
			name: "below threshold",
			scenarios: []ScenarioResult{
				{Percentage: 100},
				{Percentage: 80},
			},
			threshold: 95,
			wantPct:   90,
			wantPass:  false,
		},
		{
			name: "exactly at threshold",
			scenarios: []ScenarioResult{
				{Percentage: 95},
			},
			threshold: 95,
			wantPct:   95,
			wantPass:  true,
		},
		{
			name:      "no scenarios",
			scenarios: []ScenarioResult{},
			threshold: 95,
			wantPct:   0,
			wantPass:  false,
		},
		{
			name: "single scenario perfect",
			scenarios: []ScenarioResult{
				{Percentage: 100},
			},
			threshold: 95,
			wantPct:   100,
			wantPass:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := computeOverall(tt.scenarios, tt.threshold)
			if math.Abs(result.OverallPercentage-tt.wantPct) > 0.01 {
				t.Errorf("computeOverall() Percentage = %f, want %f", result.OverallPercentage, tt.wantPct)
			}
			if result.Pass != tt.wantPass {
				t.Errorf("computeOverall() Pass = %v, want %v", result.Pass, tt.wantPass)
			}
		})
	}
}

func TestShouldExit(t *testing.T) {
	tests := []struct {
		name string
		mode parseMode
		pass bool
		want int
	}{
		{
			name: "skill mode pass",
			mode: modeSkill,
			pass: true,
			want: 0,
		},
		{
			name: "skill mode fail",
			mode: modeSkill,
			pass: false,
			want: 1,
		},
		{
			name: "scenario mode always 0",
			mode: modeScenario,
			pass: false,
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EvalResult{Pass: tt.pass}
			got := shouldExit(tt.mode, result)
			if got != tt.want {
				t.Errorf("shouldExit(%v, {Pass: %v}) = %d, want %d", tt.mode, tt.pass, got, tt.want)
			}
		})
	}
}
