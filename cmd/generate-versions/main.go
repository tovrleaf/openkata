package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type artifactInfo struct {
	Version     string `json:"version"`
	Description string `json:"description"`
	Tags        string `json:"tags,omitempty"`
}

type versionsFile struct {
	Skills   map[string]artifactInfo `json:"skills"`
	Rules    map[string]artifactInfo `json:"rules"`
	Profiles map[string]artifactInfo `json:"profiles"`
}

func main() {
	bucket := os.Getenv("OPENKATA_BUCKET")
	if bucket == "" {
		bucket = "openkata-artifacts"
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
			description, tags := readFrontmatter(ctx, client, bucket, key)

			info := artifactInfo{
				Version:     latest,
				Description: description,
				Tags:        tags,
			}

			switch artifactType {
			case "skills":
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

	// Write locally
	if err := os.WriteFile("/tmp/versions.json", data, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "write: %v\n", err)
		os.Exit(1)
	}

	// Upload to S3
	f, _ := os.Open("/tmp/versions.json")
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

func readFrontmatter(ctx context.Context, client *s3.Client, bucket, key string) (description, tags string) {
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

	content := string(data)
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
				// Multi-line
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

		// Read tags from metadata.tags
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

func strPtr(s string) *string { return &s }
