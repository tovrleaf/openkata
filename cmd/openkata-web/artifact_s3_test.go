package main

import (
	"context"
	"io"
	"sort"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// fakeBucket is an in-memory bucket.
type fakeBucket map[string]string

func (f fakeBucket) GetObject(ctx context.Context, in *s3.GetObjectInput, opts ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	body, ok := f[*in.Key]
	if !ok {
		return nil, &types.NoSuchKey{}
	}
	return &s3.GetObjectOutput{Body: io.NopCloser(strings.NewReader(body))}, nil
}

func (f fakeBucket) ListObjectsV2(ctx context.Context, in *s3.ListObjectsV2Input, opts ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
	var keys []string
	for k := range f {
		if strings.HasPrefix(k, *in.Prefix) {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	out := &s3.ListObjectsV2Output{}
	if in.Delimiter != nil && *in.Delimiter != "" {
		seen := map[string]bool{}
		for _, k := range keys {
			rest := strings.TrimPrefix(k, *in.Prefix)
			if i := strings.Index(rest, *in.Delimiter); i >= 0 {
				p := *in.Prefix + rest[:i+1]
				if !seen[p] {
					seen[p] = true
					out.CommonPrefixes = append(out.CommonPrefixes, types.CommonPrefix{Prefix: aws.String(p)})
				}
			}
		}
		return out, nil
	}
	for _, k := range keys {
		out.Contents = append(out.Contents, types.Object{Key: aws.String(k)})
	}
	return out, nil
}

const testIndex = `{
  "skills": {
    "with-bench": {"version": "1.2.0", "description": "Has results.", "tags": "category:x",
      "tessl_score": 91, "published": true,
      "models": {
        "m-b": {"label": "Model B", "effectiveness": 70, "scenarios": [{"name": "s1", "description": "d", "pass": false, "score": 70}]},
        "m-a": {"label": "Model A", "effectiveness": 95, "scenarios": [{"name": "s1", "description": "d", "pass": true, "score": 95}]}
      }},
    "no-bench": {"version": "1.0.0", "description": "No results."}
  },
  "rules": {}, "profiles": {}
}`

func useFakeBucket(t *testing.T, objects fakeBucket) {
	t.Helper()
	oldClient, oldBucket := s3Client, bucket
	s3Client, bucket = objects, "test-bucket"
	t.Cleanup(func() { s3Client, bucket = oldClient, oldBucket })
}

func testObjects() fakeBucket {
	return fakeBucket{
		"versions.json":                    testIndex,
		"skills/with-bench/1.2.0/SKILL.md": "# With bench\n",
		"skills/with-bench/1.1.0/SKILL.md": "# With bench\n",
		"skills/no-bench/1.0.0/SKILL.md":   "# No bench\n",
	}
}

func TestLoadArtifactDetailS3Benchmarks(t *testing.T) {
	useFakeBucket(t, testObjects())
	ctx := context.Background()

	t.Run("the indexed version shows score and models", func(t *testing.T) {
		d := loadArtifactDetailS3(ctx, "skills", "with-bench", "1.2.0", "SKILL.md")
		if d == nil {
			t.Fatal("loadArtifactDetailS3() = nil, want a detail")
		}
		if d.TesslScore != 91 || !d.Published {
			t.Errorf("score %d published %v, want 91 and true", d.TesslScore, d.Published)
		}
		if len(d.Models) != 2 || d.Models[0].ID != "m-a" || d.Models[1].ID != "m-b" {
			t.Fatalf("Models = %+v, want m-a then m-b, sorted by id", d.Models)
		}
		if d.Models[0].Label != "Model A" || d.Models[0].Effectiveness != 95 ||
			len(d.Models[0].Scenarios) != 1 || !d.Models[0].Scenarios[0].Pass {
			t.Errorf("Models[0] = %+v, want the scenarios and scores of the index", d.Models[0])
		}
	})

	t.Run("no version means the indexed one", func(t *testing.T) {
		d := loadArtifactDetailS3(ctx, "skills", "with-bench", "", "SKILL.md")
		if d == nil || d.Version != "1.2.0" || len(d.Models) != 2 {
			t.Errorf("detail = %+v, want version 1.2.0 with models", d)
		}
	})

	t.Run("an older version shows neither", func(t *testing.T) {
		d := loadArtifactDetailS3(ctx, "skills", "with-bench", "1.1.0", "SKILL.md")
		if d == nil {
			t.Fatal("loadArtifactDetailS3() = nil, want a detail")
		}
		if d.TesslScore != 0 || d.Published || len(d.Models) != 0 {
			t.Errorf("score %d published %v models %v, want none for a version the index does not hold",
				d.TesslScore, d.Published, d.Models)
		}
	})

	t.Run("a skill without results has no tab", func(t *testing.T) {
		d := loadArtifactDetailS3(ctx, "skills", "no-bench", "1.0.0", "SKILL.md")
		if d == nil || d.TesslScore != 0 || len(d.Models) != 0 {
			t.Errorf("detail = %+v, want no score and no models", d)
		}
	})
}

func TestLoadArtifactDetailS3Versions(t *testing.T) {
	objects := testObjects()
	for _, v := range []string{"1.9.0", "1.10.0"} {
		objects["skills/many/"+v+"/SKILL.md"] = "# Many\n"
	}
	objects["skills/many/2.0.0/SKILL.md"] = "# Many\n"
	objects["skills/many/2.0.0/evals/results/m/x.json"] = "{}"
	objects["versions.json"] = `{"skills": {
	  "many": {"version": "2.0.0"}, "with-bench": {"version": "1.2.0"}, "no-bench": {"version": "1.0.0"}},
	  "rules": {"r": {"version": "1.0.0"}}, "profiles": {}}`
	objects["rules/r/1.0.0/RULE.md"] = "# R\n"
	objects["rules/r/0.9.0/RULE.md"] = "# R\n"
	useFakeBucket(t, objects)
	ctx := context.Background()

	tests := []struct {
		name, typ, art, doc string
		want                []string
	}{
		{"several versions, numeric order", "skills", "many", "SKILL.md", []string{"2.0.0", "1.10.0", "1.9.0"}},
		{"two versions", "skills", "with-bench", "SKILL.md", []string{"1.2.0", "1.1.0"}},
		{"one version", "skills", "no-bench", "SKILL.md", []string{"1.0.0"}},
		{"rules get the list too", "rules", "r", "RULE.md", []string{"1.0.0", "0.9.0"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := loadArtifactDetailS3(ctx, tt.typ, tt.art, "", tt.doc)
			if d == nil {
				t.Fatal("loadArtifactDetailS3() = nil")
			}
			if strings.Join(d.Versions, " ") != strings.Join(tt.want, " ") {
				t.Errorf("Versions = %v, want %v, newest first", d.Versions, tt.want)
			}
		})
	}

	t.Run("a page for an older version lists all of them", func(t *testing.T) {
		d := loadArtifactDetailS3(ctx, "skills", "many", "1.9.0", "SKILL.md")
		if d == nil || len(d.Versions) != 3 || d.Versions[0] != "2.0.0" || d.Version != "1.9.0" {
			t.Errorf("detail = %+v, want version 1.9.0 among three listed", d)
		}
	})
}

func TestVersionsWith(t *testing.T) {
	if got := versionsWith(nil, "1.0.0"); len(got) != 1 || got[0] != "1.0.0" {
		t.Errorf("versionsWith(nil) = %v, want the shown version", got)
	}
	got := versionsWith([]string{"1.10.0", "1.2.0"}, "1.9.0")
	if strings.Join(got, " ") != "1.10.0 1.9.0 1.2.0" {
		t.Errorf("versionsWith() = %v, want the shown version added in order", got)
	}
}
