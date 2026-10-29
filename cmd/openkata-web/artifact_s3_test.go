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
