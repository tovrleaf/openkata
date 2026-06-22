package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type modelInfo struct {
	Label         string         `json:"label"`
	Effectiveness int            `json:"effectiveness"`
	Scenarios     []scenarioInfo `json:"scenarios"`
}

type scenarioInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Pass        bool   `json:"pass"`
	Score       int    `json:"score"`
}

type artifactInfo struct {
	Version     string               `json:"version"`
	Description string               `json:"description"`
	Tags        string               `json:"tags,omitempty"`
	Models      map[string]modelInfo `json:"models,omitempty"`
	TesslScore  int                  `json:"tessl_score,omitempty"`
	Published   bool                 `json:"published,omitempty"`
}

type versionsFile struct {
	Skills   map[string]artifactInfo `json:"skills"`
	Rules    map[string]artifactInfo `json:"rules"`
	Profiles map[string]artifactInfo `json:"profiles"`
}

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

	versions := versionsFile{
		Skills:   make(map[string]artifactInfo),
		Rules:    make(map[string]artifactInfo),
		Profiles: make(map[string]artifactInfo),
	}

	// Scan skills
	scanLocal("skills", "SKILL.md", func(name string, info artifactInfo) {
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
	scanLocal("rules", "RULE.md", func(name string, info artifactInfo) {
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
		versions.Profiles[name] = artifactInfo{
			Version:     latestTag("profiles/" + name),
			Description: desc,
			Tags:        tags,
		}
	}

	data, _ := json.MarshalIndent(versions, "", "  ")
	data = append(data, '\n')

	os.MkdirAll(filepath.Dir(outPath), 0755)
	if err := os.WriteFile(outPath, data, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "write: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Wrote %s\n", outPath)
	fmt.Println(string(data))
}

func scanLocal(dir, mdFile string, add func(string, artifactInfo)) {
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
		add(name, artifactInfo{
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
	versions := versionsFile{
		Skills:   make(map[string]artifactInfo),
		Rules:    make(map[string]artifactInfo),
		Profiles: make(map[string]artifactInfo),
	}

	for _, artifactType := range []string{"skills", "rules", "profiles"} {
		names, err := listPrefixes(ctx, client, bucket, artifactType+"/")
		if err != nil {
			fmt.Fprintf(os.Stderr, "listing %s: %v\n", artifactType, err)
			continue
		}

		for _, name := range names {
			versionList, err := listPrefixes(ctx, client, bucket, artifactType+"/"+name+"/")
			if err != nil || len(versionList) == 0 {
				continue
			}

			sortVersions(versionList)
			latest := versionList[len(versionList)-1]

			mdFile := "SKILL.md"
			if artifactType == "rules" {
				mdFile = "RULE.md"
			} else if artifactType == "profiles" {
				mdFile = name + ".md"
			}

			key := artifactType + "/" + name + "/" + latest + "/" + mdFile
			description, tags := readFrontmatterS3(ctx, client, bucket, key)

			info := artifactInfo{
				Version:     latest,
				Description: description,
				Tags:        tags,
			}

			switch artifactType {
			case "skills":
				// Load tessl score from S3
				pluginKey := artifactType + "/" + name + "/" + latest + "/.tessl-plugin/plugin.json"
				if pResp, err := client.GetObject(ctx, &s3.GetObjectInput{Bucket: &bucket, Key: &pluginKey}); err == nil {
					var plugin struct {
						Score   int  `json:"score"`
						Private bool `json:"private"`
					}
					pBody, _ := io.ReadAll(pResp.Body)
					pResp.Body.Close()
					if json.Unmarshal(pBody, &plugin) == nil && plugin.Score > 0 {
						info.TesslScore = plugin.Score
						info.Published = !plugin.Private
					}
				}
				versions.Skills[name] = info
			case "rules":
				versions.Rules[name] = info
			case "profiles":
				versions.Profiles[name] = info
			}
		}
	}

	data, _ := json.MarshalIndent(versions, "", "  ")
	data = append(data, '\n')

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

func listPrefixes(ctx context.Context, client *s3.Client, bucket, prefix string) ([]string, error) {
	resp, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket:    &bucket,
		Prefix:    &prefix,
		Delimiter: strPtr("/"),
	})
	if err != nil {
		return nil, err
	}

	var names []string
	for _, p := range resp.CommonPrefixes {
		name := strings.TrimPrefix(*p.Prefix, prefix)
		name = strings.TrimSuffix(name, "/")
		if name != "" {
			names = append(names, name)
		}
	}
	return names, nil
}

func readFrontmatterS3(ctx context.Context, client *s3.Client, bucket, key string) (description, tags string) {
	resp, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return "", ""
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", ""
	}

	return parseFrontmatter(string(data))
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

// loadEvalResults reads eval results from skills/<name>/evals/results/
// and returns model effectiveness data for the given version.
func loadEvalResults(skillDir, currentVersion string) map[string]modelInfo {
	resultsDir := filepath.Join(skillDir, "evals", "results")
	modelDirs, err := os.ReadDir(resultsDir)
	if err != nil {
		return nil
	}

	models := make(map[string]modelInfo)
	for _, md := range modelDirs {
		if !md.IsDir() {
			continue
		}
		modelID := md.Name()
		modelDir := filepath.Join(resultsDir, modelID)

		files, err := os.ReadDir(modelDir)
		if err != nil {
			continue
		}
		var latest string
		for _, f := range files {
			if !f.IsDir() && strings.HasSuffix(f.Name(), ".json") {
				if f.Name() > latest {
					latest = f.Name()
				}
			}
		}
		if latest == "" {
			continue
		}

		data, err := os.ReadFile(filepath.Join(modelDir, latest))
		if err != nil {
			continue
		}

		var result struct {
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
		if err := json.Unmarshal(data, &result); err != nil {
			continue
		}

		if result.SkillVersion != currentVersion {
			continue
		}

		mi := modelInfo{
			Label:         result.Config.ModelLabel,
			Effectiveness: int(result.OverallPercentage + 0.5),
		}
		for _, s := range result.Scenarios {
			mi.Scenarios = append(mi.Scenarios, scenarioInfo{
				Name:        s.Name,
				Description: s.Description,
				Pass:        s.Pass,
				Score:       int(s.Percentage + 0.5),
			})
		}
		models[modelID] = mi
	}

	if len(models) == 0 {
		return nil
	}
	return models
}

func strPtr(s string) *string { return &s }
