// Package evalresults turns the reports written by the eval runner into the
// effectiveness figures of a skill. It is given the report files, wherever
// they were read from, so the same rules apply to a working tree and to a
// released version in the bucket.
package evalresults

import (
	"encoding/json"
	"path"
	"strings"

	"github.com/tovrleaf/openkata/internal/versionindex"
)

// File is one report. Path is relative to the skill's evals/results folder,
// in the form "<model>/<time>.json".
type File struct {
	Path string
	Data []byte
}

// report is the part of a runner report that the index needs.
type report struct {
	SkillVersion      string  `json:"skill_version"`
	OverallPercentage float64 `json:"overall_percentage"`
	Config            struct {
		ModelLabel string `json:"model_label"`
	} `json:"config"`
	Scenarios []struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Pass        bool    `json:"pass"`
		Percentage  float64 `json:"percentage"`
	} `json:"scenarios"`
}

// ForVersion returns the effectiveness of a skill by model, keyed by the
// model folder name, or nil when no report counts.
//
// For each model only the newest report is considered, the one whose file
// name sorts last. It counts only if it parses and was run against
// version. A newest report that does not count leaves the model out, even if
// an older one would have matched, so a stale result is never shown. Files
// that are not "<model>/<name>.json" are ignored.
func ForVersion(files []File, version string) map[string]versionindex.Model {
	newest := map[string]File{}
	for _, f := range files {
		model, name, ok := split(f.Path)
		if !ok {
			continue
		}
		if cur, seen := newest[model]; !seen || name > fileName(cur.Path) {
			newest[model] = f
		}
	}

	models := map[string]versionindex.Model{}
	for model, f := range newest {
		var r report
		if err := json.Unmarshal(f.Data, &r); err != nil {
			continue
		}
		if r.SkillVersion != version {
			continue
		}
		m := versionindex.Model{
			Label:         r.Config.ModelLabel,
			Effectiveness: round(r.OverallPercentage),
		}
		for _, s := range r.Scenarios {
			m.Scenarios = append(m.Scenarios, versionindex.Scenario{
				Name:        s.Name,
				Description: s.Description,
				Pass:        s.Pass,
				Score:       round(s.Percentage),
			})
		}
		models[model] = m
	}
	if len(models) == 0 {
		return nil
	}
	return models
}

// split separates "<model>/<name>.json" into its model and file name.
func split(p string) (model, name string, ok bool) {
	p = path.Clean(strings.ReplaceAll(p, "\\", "/"))
	parts := strings.Split(p, "/")
	if len(parts) != 2 || parts[0] == "" || parts[0] == "." || !strings.HasSuffix(parts[1], ".json") {
		return "", "", false
	}
	return parts[0], parts[1], true
}

func fileName(p string) string {
	_, name, _ := split(p)
	return name
}

func round(f float64) int { return int(f + 0.5) }
