package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/tovrleaf/openkata/internal/evalresults"
	"github.com/tovrleaf/openkata/internal/versionindex"
)

var (
	local  = flag.Bool("local", false, "Read from local filesystem instead of S3")
	output = flag.String("out", "", "Output path (default: /tmp/versions.json, or web/static/versions.json with --local)")
)

func main() {
	flag.Parse()

	if *local {
		runLocal()
	} else {
		runS3()
	}
}

func runLocal() {
	outPath := *output
	if outPath == "" {
		outPath = "web/static/versions.json"
	}

	versions := versionindex.New()

	// Scan skills
	scanLocal("skills", "SKILL.md", func(name string, info versionindex.Artifact) {
		versions.Skills[name] = info
	})

	// Load eval results for skills
	for name, info := range versions.Skills {
		models := loadEvalResults(filepath.Join("skills", name), info.Version)
		if models != nil {
			info.Models = models
			versions.Skills[name] = info
		}
	}

	// Load tessl scores for skills
	for name, info := range versions.Skills {
		pData, err := os.ReadFile(filepath.Join("skills", name, ".tessl-plugin", "plugin.json"))
		if err == nil {
			var plugin struct {
				Score   int  `json:"score"`
				Private bool `json:"private"`
			}
			if json.Unmarshal(pData, &plugin) == nil && plugin.Score > 0 {
				info.TesslScore = plugin.Score
				info.Published = !plugin.Private
				versions.Skills[name] = info
			}
		}
	}

	// Scan rules
	scanLocal("rules", "RULE.md", func(name string, info versionindex.Artifact) {
		versions.Rules[name] = info
	})

	// Scan profiles
	entries, _ := os.ReadDir("profiles")
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		data, err := os.ReadFile(filepath.Join("profiles", name, name+".md"))
		if err != nil {
			continue
		}
		desc, tags := parseFrontmatter(string(data))
		versions.Profiles[name] = versionindex.Artifact{
			Version:     latestTag("profiles/" + name),
			Description: desc,
			Tags:        tags,
		}
	}

	data, _ := versions.Marshal()

	os.MkdirAll(filepath.Dir(outPath), 0755)
	if err := os.WriteFile(outPath, data, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "write: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Wrote %s\n", outPath)
	fmt.Println(string(data))
}

func scanLocal(dir, mdFile string, add func(string, versionindex.Artifact)) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		path := filepath.Join(dir, name, mdFile)
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		desc, tags := parseFrontmatter(string(data))
		add(name, versionindex.Artifact{
			Version:     latestTag(dir + "/" + name),
			Description: desc,
			Tags:        tags,
		})
	}
}

func latestTag(prefix string) string {
	// Try to find latest git tag for this artifact
	// Fall back to "0.0.0" if none found
	cmd := fmt.Sprintf("git tag -l '%s/v*' --sort=version:refname | tail -1", prefix)
	out, err := execCmd(cmd)
	if err != nil || out == "" {
		return "0.0.0"
	}
	// Extract version from tag like "skills/create-adr/v1.0.0"
	parts := strings.Split(strings.TrimSpace(out), "/")
	v := parts[len(parts)-1]
	return strings.TrimPrefix(v, "v")
}

func execCmd(cmd string) (string, error) {
	out, err := execShell(cmd)
	return strings.TrimSpace(string(out)), err
}

func runS3() {
	bucket := os.Getenv("OPENKATA_BUCKET")
	if bucket == "" {
		bucket = "openkata-artifacts"
	}

	outPath := *output
	if outPath == "" {
		outPath = "/tmp/versions.json"
	}

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load AWS config: %v\n", err)
		os.Exit(1)
	}

	client := s3.NewFromConfig(cfg)
	versions := buildIndexFromBucket(ctx, client, bucket)

	data, _ := versions.Marshal()

	if err := os.WriteFile(outPath, data, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "write: %v\n", err)
		os.Exit(1)
	}

	// Upload to S3
	f, _ := os.Open(outPath)
	defer f.Close()
	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      &bucket,
		Key:         strPtr("versions.json"),
		Body:        f,
		ContentType: strPtr("application/json"),
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "upload: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Uploaded versions.json to s3://%s/versions.json\n", bucket)
	fmt.Println(string(data))
}

func parseFrontmatter(content string) (description, tags string) {
	if !strings.HasPrefix(content, "---") {
		return "", ""
	}

	end := strings.Index(content[3:], "---")
	if end < 0 {
		return "", ""
	}
	frontmatter := content[3 : end+3]

	lines := strings.Split(frontmatter, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "description:") {
			val := strings.TrimSpace(strings.TrimPrefix(trimmed, "description:"))
			val = strings.Trim(val, `"'`)
			if val != "" && val != ">" && val != "|" {
				description = val
			} else {
				var parts []string
				for j := i + 1; j < len(lines); j++ {
					l := lines[j]
					if len(l) > 0 && (l[0] == ' ' || l[0] == '\t') {
						parts = append(parts, strings.TrimSpace(l))
					} else {
						break
					}
				}
				description = strings.Join(parts, " ")
			}
		}

		if strings.HasPrefix(trimmed, "tags:") {
			tags = strings.TrimSpace(strings.TrimPrefix(trimmed, "tags:"))
			tags = strings.Trim(tags, `"'`)
		}

		if trimmed == "metadata:" {
			for j := i + 1; j < len(lines); j++ {
				ml := strings.TrimSpace(lines[j])
				if strings.HasPrefix(ml, "tags:") {
					tags = strings.TrimSpace(strings.TrimPrefix(ml, "tags:"))
					tags = strings.Trim(tags, `"'`)
					break
				}
				if len(lines[j]) > 0 && lines[j][0] != ' ' && lines[j][0] != '\t' {
					break
				}
			}
		}
	}
	return
}

func sortVersions(versions []string) {
	sort.Slice(versions, func(i, j int) bool {
		a := parseVersion(versions[i])
		b := parseVersion(versions[j])
		for k := 0; k < 3; k++ {
			if a[k] != b[k] {
				return a[k] < b[k]
			}
		}
		return false
	})
}

func parseVersion(v string) [3]int {
	parts := strings.Split(v, ".")
	var result [3]int
	for i := 0; i < 3 && i < len(parts); i++ {
		result[i], _ = strconv.Atoi(parts[i])
	}
	return result
}

// loadEvalResults reads the eval reports in skills/<name>/evals/results/ and
// returns the effectiveness by model for the given version. The rules for
// which reports count are in internal/evalresults.
func loadEvalResults(skillDir, currentVersion string) map[string]versionindex.Model {
	resultsDir := filepath.Join(skillDir, "evals", "results")
	var files []evalresults.File
	filepath.WalkDir(resultsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(resultsDir, path)
		if err != nil {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		files = append(files, evalresults.File{Path: filepath.ToSlash(rel), Data: data})
		return nil
	})
	return evalresults.ForVersion(files, currentVersion)
}

func strPtr(s string) *string { return &s }
