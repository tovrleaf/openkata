package main

import (
	"fmt"
	"testing"
)

// mockCompleter returns preset responses for testing.
type mockCompleter struct {
	responses []string
	calls     int
}

func (m *mockCompleter) Complete(system, user string) (string, error) {
	if m.calls >= len(m.responses) {
		return "", fmt.Errorf("no more mock responses")
	}
	resp := m.responses[m.calls]
	m.calls++
	return resp, nil
}

func TestParseJudgeResponse(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		want    int
		wantErr bool
	}{
		{
			name: "clean JSON array",
			raw:  `[{"name": "a", "pass": true, "reason": "good"}]`,
			want: 1,
		},
		{
			name: "with preamble text",
			raw:  "Here is my evaluation:\n[{\"name\": \"a\", \"pass\": true, \"reason\": \"ok\"}]",
			want: 1,
		},
		{
			name: "with ANSI codes",
			raw:  "\x1b[32m[{\"name\": \"a\", \"pass\": false, \"reason\": \"bad\"}]\x1b[0m",
			want: 1,
		},
		{
			name: "with header lines",
			raw:  "> Model: claude\n[{\"name\": \"x\", \"pass\": true, \"reason\": \"yes\"}]",
			want: 1,
		},
		{
			name: "trailing text after JSON",
			raw:  `[{"name": "a", "pass": true, "reason": "ok"}] some extra text`,
			want: 1,
		},
		{
			name:    "no JSON at all",
			raw:     "This response has no structured data",
			wantErr: true,
		},
		{
			name: "multiple verdicts",
			raw:  `[{"name": "a", "pass": true, "reason": "x"}, {"name": "b", "pass": false, "reason": "y"}]`,
			want: 2,
		},
		{
			name: "wrapped in markdown code fences",
			raw:  "```json\n[{\"name\": \"a\", \"pass\": true, \"reason\": \"ok\"}]\n```",
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verdicts, err := parseJudgeResponse(tt.raw)
			if tt.wantErr {
				if err == nil {
					t.Error("parseJudgeResponse() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("parseJudgeResponse() error: %v", err)
			}
			if len(verdicts) != tt.want {
				t.Errorf("parseJudgeResponse() got %d verdicts, want %d", len(verdicts), tt.want)
			}
		})
	}
}

func TestJudgeEval_HappyPath(t *testing.T) {
	mock := &mockCompleter{
		responses: []string{
			`[{"name": "criterion-a", "pass": true, "reason": "good"}, {"name": "criterion-b", "pass": false, "reason": "missing"}]`,
		},
	}

	criteria := &CriteriaFile{
		Context: "test context",
		Checklist: []Criterion{
			{Name: "criterion-a", Description: "Does A", MaxScore: 10},
			{Name: "criterion-b", Description: "Does B", MaxScore: 5},
		},
	}

	verdicts, err := judgeEval(mock, "A detailed agent response that is more than twenty characters.", criteria)
	if err != nil {
		t.Fatalf("judgeEval() error: %v", err)
	}
	if len(verdicts) != 2 {
		t.Fatalf("judgeEval() got %d verdicts, want 2", len(verdicts))
	}
	if !verdicts[0].Pass {
		t.Error("judgeEval() verdicts[0] should pass")
	}
	if verdicts[1].Pass {
		t.Error("judgeEval() verdicts[1] should fail")
	}
}

func TestJudgeEval_RetryOnParseFailure(t *testing.T) {
	mock := &mockCompleter{
		responses: []string{
			"This is not JSON at all",
			`[{"name": "a", "pass": true, "reason": "recovered"}]`,
		},
	}

	criteria := &CriteriaFile{
		Checklist: []Criterion{
			{Name: "a", Description: "test", MaxScore: 10},
		},
	}

	verdicts, err := judgeEval(mock, "A sufficiently long agent response for testing purposes.", criteria)
	if err != nil {
		t.Fatalf("judgeEval() should succeed on retry, got: %v", err)
	}
	if len(verdicts) != 1 {
		t.Fatalf("judgeEval() got %d verdicts, want 1", len(verdicts))
	}
	if mock.calls != 2 {
		t.Errorf("judgeEval() made %d calls, want 2 (retry)", mock.calls)
	}
}

func TestJudgeEval_SuspiciousResult(t *testing.T) {
	mock := &mockCompleter{
		responses: []string{
			`[{"name": "a", "pass": true, "reason": "ok"}]`,
		},
	}

	criteria := &CriteriaFile{
		Checklist: []Criterion{
			{Name: "a", Description: "test", MaxScore: 10},
		},
	}

	// Short agent response + all pass = suspicious
	_, err := judgeEval(mock, "short", criteria)
	if err == nil {
		t.Error("judgeEval() should fail with suspicious result for trivial agent response")
	}
	if err != nil && !contains(err.Error(), "suspicious") {
		t.Errorf("judgeEval() error should mention 'suspicious', got: %v", err)
	}
}

func TestIsSuspicious(t *testing.T) {
	tests := []struct {
		name     string
		response string
		verdicts []Verdict
		want     bool
	}{
		{
			name:     "short response all pass",
			response: "ok",
			verdicts: []Verdict{{Pass: true}},
			want:     true,
		},
		{
			name:     "long response all pass",
			response: "This is a long enough response with good content",
			verdicts: []Verdict{{Pass: true}},
			want:     false,
		},
		{
			name:     "short response some fail",
			response: "ok",
			verdicts: []Verdict{{Pass: true}, {Pass: false}},
			want:     false,
		},
		{
			name:     "empty verdicts",
			response: "ok",
			verdicts: []Verdict{},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isSuspicious(tt.response, tt.verdicts)
			if got != tt.want {
				t.Errorf("isSuspicious(%q, %v) = %v, want %v", tt.response, tt.verdicts, got, tt.want)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsStr(s, substr))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
