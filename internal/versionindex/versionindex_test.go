package versionindex

import (
	"bytes"
	"strings"
	"testing"
)

// golden is versions.json as the generator writes it: every field in use,
// in the order of the structs.
const golden = `{
  "skills": {
    "grill-me": {
      "version": "1.0.1",
      "description": "Challenge a plan.",
      "tags": "category:review, category:architecture",
      "models": {
        "claude-sonnet-4.6": {
          "label": "Sonnet 4.6",
          "effectiveness": 96,
          "scenarios": [
            {
              "name": "scenario-1",
              "description": "Finds the gaps",
              "pass": true,
              "score": 100
            }
          ]
        }
      },
      "tessl_score": 100,
      "published": true
    }
  },
  "rules": {
    "git-naming": {
      "version": "1.0.0",
      "description": "Branch and commit naming.",
      "tags": "category:conventions, tool:git"
    }
  },
  "profiles": {}
}
`

func TestParseAndMarshalRoundTrip(t *testing.T) {
	idx, err := Parse([]byte(golden))
	if err != nil {
		t.Fatalf("Parse(golden) error = %v", err)
	}
	got, err := idx.Marshal()
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	if !bytes.Equal(got, []byte(golden)) {
		t.Errorf("Marshal(Parse(golden)) differs from golden:\n%s\nwant:\n%s", got, golden)
	}
}

func TestParseFields(t *testing.T) {
	idx, err := Parse([]byte(golden))
	if err != nil {
		t.Fatalf("Parse(golden) error = %v", err)
	}
	a := idx.Skills["grill-me"]
	if a.Version != "1.0.1" || a.TesslScore != 100 || !a.Published {
		t.Errorf("Skills[grill-me] = %+v, want version 1.0.1, score 100, published", a)
	}
	m := a.Models["claude-sonnet-4.6"]
	if m.Label != "Sonnet 4.6" || m.Effectiveness != 96 || len(m.Scenarios) != 1 || !m.Scenarios[0].Pass {
		t.Errorf("Models[claude-sonnet-4.6] = %+v, want label, 96 and one passing scenario", m)
	}
	if idx.Rules["git-naming"].Models != nil {
		t.Errorf("Rules[git-naming].Models = %v, want nil for a rule", idx.Rules["git-naming"].Models)
	}
}

func TestParseMissingSections(t *testing.T) {
	tests := []struct {
		name string
		in   string
	}{
		{"empty object", `{}`},
		{"only skills", `{"skills": {"a": {"version": "1.0.0", "description": "d"}}}`},
		{"null sections", `{"skills": null, "rules": null, "profiles": null}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx, err := Parse([]byte(tt.in))
			if err != nil {
				t.Fatalf("Parse(%q) error = %v", tt.in, err)
			}
			if idx.Skills == nil || idx.Rules == nil || idx.Profiles == nil {
				t.Errorf("Parse(%q) has a nil map: %+v, want empty maps", tt.in, idx)
			}
		})
	}
}

func TestParseIgnoresUnknownFields(t *testing.T) {
	in := `{"skills": {"a": {"version": "1.0.0", "description": "d", "future": 1}}, "extra": true}`
	idx, err := Parse([]byte(in))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if idx.Skills["a"].Version != "1.0.0" {
		t.Errorf("Skills[a].Version = %q, want 1.0.0", idx.Skills["a"].Version)
	}
}

func TestParseErrors(t *testing.T) {
	for _, in := range []string{"", "not json", `[]`, `{"skills": []}`, `{"skills": {"a": "x"}}`} {
		t.Run(in, func(t *testing.T) {
			if _, err := Parse([]byte(in)); err == nil {
				t.Errorf("Parse(%q) error = nil, want an error", in)
			}
		})
	}
}

func TestOf(t *testing.T) {
	idx := New()
	idx.Skills["s"] = Artifact{Version: "1.0.0"}
	idx.Rules["r"] = Artifact{Version: "2.0.0"}
	idx.Profiles["p"] = Artifact{Version: "3.0.0"}

	tests := []struct {
		artifactType string
		name         string
		want         string
	}{
		{"skills", "s", "1.0.0"},
		{"rules", "r", "2.0.0"},
		{"profiles", "p", "3.0.0"},
	}
	for _, tt := range tests {
		if got := idx.Of(tt.artifactType)[tt.name].Version; got != tt.want {
			t.Errorf("Of(%q)[%q].Version = %q, want %q", tt.artifactType, tt.name, got, tt.want)
		}
	}
	if got := idx.Of("prompts"); got != nil {
		t.Errorf("Of(prompts) = %v, want nil", got)
	}
}

func TestMarshalOmitsEmptyOptionalFields(t *testing.T) {
	idx := New()
	idx.Skills["a"] = Artifact{Version: "1.0.0", Description: "d"}
	data, err := idx.Marshal()
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	for _, field := range []string{"tags", "models", "tessl_score", "published"} {
		if strings.Contains(string(data), `"`+field+`"`) {
			t.Errorf("Marshal() wrote %q for an artifact without it:\n%s", field, data)
		}
	}
	if !bytes.HasSuffix(data, []byte("}\n")) {
		t.Errorf("Marshal() = %q, want a trailing newline", data)
	}
}
