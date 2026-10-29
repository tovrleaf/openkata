package evalresults

import (
	"fmt"
	"testing"

	"github.com/tovrleaf/openkata/internal/versionindex"
)

// reportJSON builds a runner report for a skill version.
func reportJSON(version, label string, overall float64, scenarios ...string) []byte {
	s := ""
	for i, sc := range scenarios {
		if i > 0 {
			s += ","
		}
		s += sc
	}
	return []byte(fmt.Sprintf(
		`{"skill":"x","skill_version":%q,"config":{"model_label":%q},"overall_percentage":%v,"scenarios":[%s]}`,
		version, label, overall, s))
}

func scenarioJSON(name string, pct float64, pass bool) string {
	return fmt.Sprintf(`{"name":%q,"description":"%s desc","percentage":%v,"pass":%t}`, name, name, pct, pass)
}

func TestForVersion(t *testing.T) {
	tests := []struct {
		name    string
		files   []File
		version string
		want    map[string]versionindex.Model
	}{
		{
			name:    "no files",
			files:   nil,
			version: "1.0.0",
			want:    nil,
		},
		{
			name: "one model, rounded",
			files: []File{{
				Path: "claude-sonnet-4.6/2026-06-01T100000.json",
				Data: reportJSON("1.0.0", "Sonnet 4.6", 95.5,
					scenarioJSON("scenario-1", 100, true),
					scenarioJSON("scenario-2", 66.67, false)),
			}},
			version: "1.0.0",
			want: map[string]versionindex.Model{
				"claude-sonnet-4.6": {
					Label:         "Sonnet 4.6",
					Effectiveness: 96,
					Scenarios: []versionindex.Scenario{
						{Name: "scenario-1", Description: "scenario-1 desc", Pass: true, Score: 100},
						{Name: "scenario-2", Description: "scenario-2 desc", Pass: false, Score: 67},
					},
				},
			},
		},
		{
			name: "rounds down below the half",
			files: []File{{
				Path: "m/a.json",
				Data: reportJSON("1.0.0", "M", 95.4),
			}},
			version: "1.0.0",
			want:    map[string]versionindex.Model{"m": {Label: "M", Effectiveness: 95}},
		},
		{
			name: "models are independent",
			files: []File{
				{Path: "a/2026-06-01T100000.json", Data: reportJSON("1.0.0", "A", 90)},
				{Path: "b/2026-06-02T100000.json", Data: reportJSON("1.0.0", "B", 80)},
			},
			version: "1.0.0",
			want: map[string]versionindex.Model{
				"a": {Label: "A", Effectiveness: 90},
				"b": {Label: "B", Effectiveness: 80},
			},
		},
		{
			name: "the newest file name wins, whatever its score",
			files: []File{
				{Path: "m/2026-06-01T100000.json", Data: reportJSON("1.0.0", "M", 99)},
				{Path: "m/2026-06-09T100000.json", Data: reportJSON("1.0.0", "M", 70)},
				{Path: "m/2026-06-05T100000.json", Data: reportJSON("1.0.0", "M", 85)},
			},
			version: "1.0.0",
			want:    map[string]versionindex.Model{"m": {Label: "M", Effectiveness: 70}},
		},
		{
			name: "a report for another version does not count",
			files: []File{
				{Path: "m/2026-06-01T100000.json", Data: reportJSON("1.0.0", "M", 99)},
			},
			version: "1.1.0",
			want:    nil,
		},
		{
			name: "a stale newest report hides the model even if an older one matches",
			files: []File{
				{Path: "m/2026-06-01T100000.json", Data: reportJSON("1.1.0", "M", 99)},
				{Path: "m/2026-06-09T100000.json", Data: reportJSON("1.0.0", "M", 70)},
			},
			version: "1.1.0",
			want:    nil,
		},
		{
			name: "one model can count while another does not",
			files: []File{
				{Path: "a/x.json", Data: reportJSON("1.0.0", "A", 90)},
				{Path: "b/x.json", Data: reportJSON("0.9.0", "B", 80)},
			},
			version: "1.0.0",
			want:    map[string]versionindex.Model{"a": {Label: "A", Effectiveness: 90}},
		},
		{
			name: "a malformed newest report leaves the model out",
			files: []File{
				{Path: "m/2026-06-01T100000.json", Data: reportJSON("1.0.0", "M", 99)},
				{Path: "m/2026-06-09T100000.json", Data: []byte("not json")},
			},
			version: "1.0.0",
			want:    nil,
		},
		{
			name: "files that are not model reports are ignored",
			files: []File{
				{Path: "m/notes.txt", Data: reportJSON("1.0.0", "M", 99)},
				{Path: "top-level.json", Data: reportJSON("1.0.0", "M", 99)},
				{Path: "m/sub/deep.json", Data: reportJSON("1.0.0", "M", 99)},
				{Path: "", Data: reportJSON("1.0.0", "M", 99)},
			},
			version: "1.0.0",
			want:    nil,
		},
		{
			name: "no version matches a report without one",
			files: []File{
				{Path: "m/x.json", Data: reportJSON("", "M", 88)},
			},
			version: "",
			want:    map[string]versionindex.Model{"m": {Label: "M", Effectiveness: 88}},
		},
		{
			name: "windows separators are accepted",
			files: []File{
				{Path: `m\x.json`, Data: reportJSON("1.0.0", "M", 88)},
			},
			version: "1.0.0",
			want:    map[string]versionindex.Model{"m": {Label: "M", Effectiveness: 88}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ForVersion(tt.files, tt.version)
			if len(got) != len(tt.want) {
				t.Fatalf("ForVersion(%q) has %d models %v, want %d %v", tt.version, len(got), got, len(tt.want), tt.want)
			}
			for model, want := range tt.want {
				g, ok := got[model]
				if !ok {
					t.Errorf("ForVersion(%q) is missing model %q", tt.version, model)
					continue
				}
				if g.Label != want.Label || g.Effectiveness != want.Effectiveness || len(g.Scenarios) != len(want.Scenarios) {
					t.Errorf("ForVersion(%q)[%q] = %+v, want %+v", tt.version, model, g, want)
					continue
				}
				for i, sc := range want.Scenarios {
					if g.Scenarios[i] != sc {
						t.Errorf("ForVersion(%q)[%q].Scenarios[%d] = %+v, want %+v", tt.version, model, i, g.Scenarios[i], sc)
					}
				}
			}
		})
	}
}

func TestForVersionIgnoresInputOrder(t *testing.T) {
	a := File{Path: "m/2026-06-01T100000.json", Data: reportJSON("1.0.0", "M", 60)}
	b := File{Path: "m/2026-06-09T100000.json", Data: reportJSON("1.0.0", "M", 80)}
	for _, files := range [][]File{{a, b}, {b, a}} {
		got := ForVersion(files, "1.0.0")
		if got["m"].Effectiveness != 80 {
			t.Errorf("ForVersion(%v) effectiveness = %d, want 80 whatever the order", []string{files[0].Path, files[1].Path}, got["m"].Effectiveness)
		}
	}
}
