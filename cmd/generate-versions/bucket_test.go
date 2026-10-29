package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// fakeBucket is an in-memory bucket.
type fakeBucket struct {
	objects  map[string]string
	pageSize int              // keys per page for listings without a delimiter, 0 for one page
	listErr  map[string]error // prefix to the error its listing returns
	getErr   map[string]error // key to the error reading it returns
	gets     []string         // every key that was read
}

func (f *fakeBucket) GetObject(ctx context.Context, in *s3.GetObjectInput, opts ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	key := *in.Key
	f.gets = append(f.gets, key)
	if err, ok := f.getErr[key]; ok {
		return nil, err
	}
	body, ok := f.objects[key]
	if !ok {
		return nil, &types.NoSuchKey{}
	}
	return &s3.GetObjectOutput{Body: io.NopCloser(strings.NewReader(body))}, nil
}

func (f *fakeBucket) ListObjectsV2(ctx context.Context, in *s3.ListObjectsV2Input, opts ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
	prefix := *in.Prefix
	if err, ok := f.listErr[prefix]; ok {
		return nil, err
	}
	var keys []string
	for k := range f.objects {
		if strings.HasPrefix(k, prefix) {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	if in.Delimiter != nil && *in.Delimiter != "" {
		seen := map[string]bool{}
		out := &s3.ListObjectsV2Output{}
		for _, k := range keys {
			rest := strings.TrimPrefix(k, prefix)
			i := strings.Index(rest, "/")
			if i < 0 {
				continue
			}
			p := prefix + rest[:i+1]
			if !seen[p] {
				seen[p] = true
				out.CommonPrefixes = append(out.CommonPrefixes, types.CommonPrefix{Prefix: aws.String(p)})
			}
		}
		return out, nil
	}

	start := 0
	if in.ContinuationToken != nil {
		start, _ = strconv.Atoi(*in.ContinuationToken)
	}
	end := len(keys)
	if f.pageSize > 0 && start+f.pageSize < end {
		end = start + f.pageSize
	}
	out := &s3.ListObjectsV2Output{}
	for _, k := range keys[start:end] {
		out.Contents = append(out.Contents, types.Object{Key: aws.String(k)})
	}
	if end < len(keys) {
		out.IsTruncated = aws.Bool(true)
		out.NextContinuationToken = aws.String(strconv.Itoa(end))
	}
	return out, nil
}

func report(version, label string, overall float64) string {
	return fmt.Sprintf(
		`{"skill_version":%q,"config":{"model_label":%q},"overall_percentage":%v,`+
			`"scenarios":[{"name":"scenario-1","description":"first","percentage":%v,"pass":true}]}`,
		version, label, overall, overall)
}

func skillMD(name, desc, tags string) string {
	return fmt.Sprintf("---\nname: %s\ndescription: >\n  %s\nmetadata:\n  tags: %q\n---\n\n# %s\n", name, desc, tags, name)
}

// testBucket holds skills, a rule and a profile, with several versions.
func testBucket() *fakeBucket {
	return &fakeBucket{objects: map[string]string{
		// alpha: 1.10.0 is newer than 1.9.0 although it sorts first as text.
		"skills/alpha/1.9.0/SKILL.md":                                 skillMD("alpha", "Old alpha.", "category:old"),
		"skills/alpha/1.9.0/evals/results/m1/2026-05-01T100000.json":  report("1.9.0", "Model one", 50),
		"skills/alpha/1.10.0/SKILL.md":                                skillMD("alpha", "Alpha does things.", "category:review, tool:git"),
		"skills/alpha/1.10.0/.tessl-plugin/plugin.json":               `{"score": 97, "private": false}`,
		"skills/alpha/1.10.0/evals/results/m1/2026-06-01T100000.json": report("1.10.0", "Model one", 90),
		"skills/alpha/1.10.0/evals/results/m1/2026-06-09T100000.json": report("1.10.0", "Model one", 92.4),
		"skills/alpha/1.10.0/evals/results/m2/2026-06-02T100000.json": report("1.9.0", "Model two", 80),
		"skills/alpha/1.10.0/evals/results/README.md":                 "not a report",
		// beta has no tessl file and no evals.
		"skills/beta/1.0.0/SKILL.md": "# beta\n",
		// gamma has a zero score, delta a private plugin.
		"skills/gamma/1.0.0/SKILL.md":                  skillMD("gamma", "Gamma.", "category:x"),
		"skills/gamma/1.0.0/.tessl-plugin/plugin.json": `{"score": 0, "private": false}`,
		"skills/delta/2.0.0/SKILL.md":                  skillMD("delta", "Delta.", "category:y"),
		"skills/delta/2.0.0/.tessl-plugin/plugin.json": `{"score": 88, "private": true}`,
		// A rule never gets scores or models, even if it holds reports.
		"rules/rule-a/2.0.0/RULE.md":                   "---\nname: rule-a\ndescription: A rule.\ntags: category:conventions\n---\n",
		"rules/rule-a/2.0.0/.tessl-plugin/plugin.json": `{"score": 99, "private": false}`,
		"rules/rule-a/2.0.0/evals/results/m1/x.json":   report("2.0.0", "Model one", 70),
		"profiles/prof-a/1.0.0/prof-a.md":              "---\ndescription: A profile.\ntags: category:review\n---\n",
	}}
}

func TestBuildIndexFromBucket(t *testing.T) {
	idx := buildIndexFromBucket(context.Background(), testBucket(), "bucket")

	t.Run("the newest version by number is used", func(t *testing.T) {
		if got := idx.Skills["alpha"].Version; got != "1.10.0" {
			t.Errorf("Skills[alpha].Version = %q, want 1.10.0 (1.10.0 is newer than 1.9.0)", got)
		}
	})

	t.Run("description and tags come from that version", func(t *testing.T) {
		a := idx.Skills["alpha"]
		if a.Description != "Alpha does things." || a.Tags != "category:review, tool:git" {
			t.Errorf("Skills[alpha] description %q tags %q, want those of 1.10.0", a.Description, a.Tags)
		}
	})

	t.Run("tessl score and published flag", func(t *testing.T) {
		tests := []struct {
			name      string
			score     int
			published bool
		}{
			{"alpha", 97, true},
			{"beta", 0, false},
			{"gamma", 0, false},
			{"delta", 88, false},
		}
		for _, tt := range tests {
			a := idx.Skills[tt.name]
			if a.TesslScore != tt.score || a.Published != tt.published {
				t.Errorf("Skills[%s] score %d published %v, want %d and %v", tt.name, a.TesslScore, a.Published, tt.score, tt.published)
			}
		}
	})

	t.Run("effectiveness comes from the newest report of the version", func(t *testing.T) {
		models := idx.Skills["alpha"].Models
		if len(models) != 1 {
			t.Fatalf("Skills[alpha].Models = %v, want only m1: m2 was run against another version", models)
		}
		m := models["m1"]
		if m.Label != "Model one" || m.Effectiveness != 92 {
			t.Errorf("Models[m1] = %+v, want label Model one and 92, from the newest report", m)
		}
		if len(m.Scenarios) != 1 || m.Scenarios[0].Name != "scenario-1" || !m.Scenarios[0].Pass || m.Scenarios[0].Score != 92 {
			t.Errorf("Models[m1].Scenarios = %+v, want the scenario of the newest report", m.Scenarios)
		}
	})

	t.Run("a skill without evals has no models", func(t *testing.T) {
		for _, name := range []string{"beta", "gamma", "delta"} {
			if got := idx.Skills[name].Models; got != nil {
				t.Errorf("Skills[%s].Models = %v, want none", name, got)
			}
		}
	})

	t.Run("rules and profiles have no scores or models", func(t *testing.T) {
		r := idx.Rules["rule-a"]
		if r.Version != "2.0.0" || r.Models != nil || r.TesslScore != 0 || r.Published {
			t.Errorf("Rules[rule-a] = %+v, want version 2.0.0 and no scores or models", r)
		}
		p := idx.Profiles["prof-a"]
		if p.Version != "1.0.0" || p.Description != "A profile." || p.Models != nil {
			t.Errorf("Profiles[prof-a] = %+v, want version 1.0.0 and its description", p)
		}
	})
}

func TestBuildIndexFromBucketReadsOnlyWhatItNeeds(t *testing.T) {
	b := testBucket()
	buildIndexFromBucket(context.Background(), b, "bucket")

	reports := map[string]bool{}
	for _, k := range b.gets {
		if strings.Contains(k, "/evals/results/") {
			reports[k] = true
		}
	}
	want := map[string]bool{
		"skills/alpha/1.10.0/evals/results/m1/2026-06-09T100000.json": true,
		"skills/alpha/1.10.0/evals/results/m2/2026-06-02T100000.json": true,
	}
	if len(reports) != len(want) {
		t.Errorf("reports read = %v, want only the newest of each model in the newest version: %v", reports, want)
	}
	for k := range want {
		if !reports[k] {
			t.Errorf("report %s was not read", k)
		}
	}
}

func TestBuildIndexFromBucketFollowsPages(t *testing.T) {
	b := testBucket()
	b.pageSize = 1
	idx := buildIndexFromBucket(context.Background(), b, "bucket")
	if got := idx.Skills["alpha"].Models["m1"].Effectiveness; got != 92 {
		t.Errorf("Skills[alpha].Models[m1].Effectiveness = %d with one key per page, want 92", got)
	}
}

func TestBuildIndexFromBucketSurvivesFailures(t *testing.T) {
	tests := []struct {
		name  string
		setup func(*fakeBucket)
		check func(*testing.T, *fakeBucket)
	}{
		{
			name: "listing the reports fails",
			setup: func(b *fakeBucket) {
				b.listErr = map[string]error{"skills/alpha/1.10.0/evals/results/": errors.New("boom")}
			},
		},
		{
			name: "reading the newest report fails",
			setup: func(b *fakeBucket) {
				b.getErr = map[string]error{"skills/alpha/1.10.0/evals/results/m1/2026-06-09T100000.json": errors.New("boom")}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := testBucket()
			tt.setup(b)
			idx := buildIndexFromBucket(context.Background(), b, "bucket")

			a, ok := idx.Skills["alpha"]
			if !ok || a.Version != "1.10.0" || a.TesslScore != 97 {
				t.Fatalf("Skills[alpha] = %+v, want it indexed with its version and score", a)
			}
			if _, has := a.Models["m1"]; has {
				t.Errorf("Skills[alpha].Models = %v, want no m1 when its report cannot be read", a.Models)
			}
			if idx.Skills["delta"].TesslScore != 88 || idx.Rules["rule-a"].Version != "2.0.0" {
				t.Errorf("other artifacts were affected: delta %+v, rule %+v", idx.Skills["delta"], idx.Rules["rule-a"])
			}
		})
	}
}

func TestBuildIndexFromEmptyBucket(t *testing.T) {
	idx := buildIndexFromBucket(context.Background(), &fakeBucket{objects: map[string]string{}}, "bucket")
	if idx.Skills == nil || idx.Rules == nil || idx.Profiles == nil || len(idx.Skills)+len(idx.Rules)+len(idx.Profiles) != 0 {
		t.Errorf("index of an empty bucket = %+v, want empty maps", idx)
	}
}

func TestBuildIndexOutputShape(t *testing.T) {
	data, err := buildIndexFromBucket(context.Background(), testBucket(), "bucket").Marshal()
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	out := string(data)
	for _, want := range []string{
		`"tessl_score": 97`, `"published": true`, `"models": {`, `"effectiveness": 92`, `"label": "Model one"`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("index = %s\nwant it to contain %s", out, want)
		}
	}
}
