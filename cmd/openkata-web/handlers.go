package main

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/tovrleaf/openkata/cmd/openkata-web/templates"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	templates.Home().Render(r.Context(), w)
}

func handleDesignSystem(w http.ResponseWriter, r *http.Request) {
	templates.DesignSystem().Render(r.Context(), w)
}

func handleSkills(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/skills/")

	// /skills/ — listing
	if path == "" {
		skills := loadSkillsList(r.Context())
		templates.Skills(skills).Render(r.Context(), w)
		return
	}

	// /skills/:name/archive or /skills/:name/archive/:version
	parts := strings.Split(strings.TrimSuffix(path, "/"), "/")
	if len(parts) >= 2 && parts[1] == "archive" {
		version := ""
		if len(parts) >= 3 {
			version = parts[2]
		}
		handleArchive(w, r, "skills", parts[0], version)
		return
	}

	http.NotFound(w, r)
}

func loadSkillsList(ctx context.Context) []templates.SkillEntry {
	if s3Client == nil {
		return nil
	}

	// Read versions.json from S3
	key := "versions.json"
	resp, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	var versions struct {
		Skills map[string]struct {
			Version     string `json:"version"`
			Description string `json:"description"`
			Tags        string `json:"tags"`
		} `json:"skills"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		return nil
	}

	// Read download counts
	counts := make(map[string]int)
	if dbClient != nil {
		scanResp, err := dbClient.Scan(ctx, &dynamodb.ScanInput{
			TableName: &table,
		})
		if err == nil {
			for _, item := range scanResp.Items {
				artAttr, ok := item["artifact"].(*types.AttributeValueMemberS)
				if !ok {
					continue
				}
				dlAttr, ok := item["downloads"].(*types.AttributeValueMemberN)
				if !ok {
					continue
				}
				n, _ := strconv.Atoi(dlAttr.Value)
				counts[artAttr.Value] = n
			}
		}
	}

	var skills []templates.SkillEntry
	for name, info := range versions.Skills {
		skills = append(skills, templates.SkillEntry{
			Name:        name,
			Version:     info.Version,
			Description: info.Description,
			Tags:        info.Tags,
			Downloads:   counts["skills/"+name],
		})
	}
	sort.Slice(skills, func(i, j int) bool { return skills[i].Name < skills[j].Name })
	return skills
}

func handleRules(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/rules/")

	// /rules/ — listing
	if path == "" {
		rules := loadRulesList(r.Context())
		templates.Rules(rules).Render(r.Context(), w)
		return
	}

	// /rules/:name/archive or /rules/:name/archive/:version
	parts := strings.Split(strings.TrimSuffix(path, "/"), "/")
	if len(parts) >= 2 && parts[1] == "archive" {
		version := ""
		if len(parts) >= 3 {
			version = parts[2]
		}
		handleArchive(w, r, "rules", parts[0], version)
		return
	}

	http.NotFound(w, r)
}

func loadRulesList(ctx context.Context) []templates.SkillEntry {
	if s3Client == nil {
		return nil
	}

	key := "versions.json"
	resp, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	var versions struct {
		Rules map[string]struct {
			Version     string `json:"version"`
			Description string `json:"description"`
			Tags        string `json:"tags"`
		} `json:"rules"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		return nil
	}

	counts := make(map[string]int)
	if dbClient != nil {
		scanResp, err := dbClient.Scan(ctx, &dynamodb.ScanInput{
			TableName: &table,
		})
		if err == nil {
			for _, item := range scanResp.Items {
				artAttr, ok := item["artifact"].(*types.AttributeValueMemberS)
				if !ok {
					continue
				}
				dlAttr, ok := item["downloads"].(*types.AttributeValueMemberN)
				if !ok {
					continue
				}
				n, _ := strconv.Atoi(dlAttr.Value)
				counts[artAttr.Value] = n
			}
		}
	}

	var rules []templates.SkillEntry
	for name, info := range versions.Rules {
		rules = append(rules, templates.SkillEntry{
			Name:        name,
			Version:     info.Version,
			Description: info.Description,
			Tags:        info.Tags,
			Downloads:   counts["rules/"+name],
		})
	}
	sort.Slice(rules, func(i, j int) bool { return rules[i].Name < rules[j].Name })
	return rules
}

func handleArchive(w http.ResponseWriter, r *http.Request, artifactType, name, version string) {
	if s3Client == nil {
		http.Error(w, "service unavailable", http.StatusServiceUnavailable)
		return
	}

	ctx := r.Context()

	// Resolve latest version if not specified
	if version == "" {
		version = resolveLatestVersion(ctx, artifactType, name)
		if version == "" {
			http.NotFound(w, r)
			return
		}
	}

	// List all files in the S3 prefix
	prefix := artifactType + "/" + name + "/" + version + "/"
	resp, err := s3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: &bucket,
		Prefix: &prefix,
	})
	if err != nil || len(resp.Contents) == 0 {
		http.NotFound(w, r)
		return
	}

	// Collect files
	type fileEntry struct {
		path string
		data []byte
	}
	var files []fileEntry

	for _, obj := range resp.Contents {
		key := *obj.Key
		relPath := strings.TrimPrefix(key, prefix)
		if relPath == "" || strings.HasSuffix(relPath, "/") {
			continue
		}
		if relPath == "tile.json" || relPath == "references/ACKNOWLEDGMENTS.md" {
			continue
		}

		getResp, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: &bucket,
			Key:    &key,
		})
		if err != nil {
			continue
		}
		data, err := io.ReadAll(getResp.Body)
		getResp.Body.Close()
		if err != nil {
			continue
		}
		files = append(files, fileEntry{path: relPath, data: data})
	}

	if len(files) == 0 {
		http.NotFound(w, r)
		return
	}

	// Generate checksums
	checksums := make(map[string]string)
	for _, f := range files {
		h := sha256Sum(f.data)
		checksums[f.path] = "sha256:" + h
	}
	checksumJSON, _ := json.MarshalIndent(checksums, "    ", "  ")

	manifest := fmt.Sprintf(`{
  "name": %q,
  "version": %q,
  "source": "github.com/tovrleaf/openkata",
  "installedAt": "",
  "checksums": %s
}
`, name, version, string(checksumJSON))

	// Stream tar.gz response
	filename := fmt.Sprintf("%s-%s.tar.gz", name, version)
	w.Header().Set("Content-Type", "application/gzip")
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)

	gw := gzip.NewWriter(w)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	for _, f := range files {
		tw.WriteHeader(&tar.Header{
			Name: name + "/" + f.path,
			Size: int64(len(f.data)),
			Mode: 0644,
		})
		tw.Write(f.data)
	}

	// Write .manifest.json
	tw.WriteHeader(&tar.Header{
		Name: name + "/.manifest.json",
		Size: int64(len(manifest)),
		Mode: 0644,
	})
	tw.Write([]byte(manifest))

	// Increment download counter
	incrementDownload(ctx, artifactType+"/"+name)
}

func resolveLatestVersion(ctx context.Context, artifactType, name string) string {
	key := "versions.json"
	resp, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	var v struct {
		Skills map[string]struct{ Version string } `json:"skills"`
		Rules  map[string]struct{ Version string } `json:"rules"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return ""
	}

	if artifactType == "skills" {
		if s, ok := v.Skills[name]; ok {
			return s.Version
		}
	} else {
		if r, ok := v.Rules[name]; ok {
			return r.Version
		}
	}
	return ""
}

func incrementDownload(ctx context.Context, artifact string) {
	if dbClient == nil {
		return
	}
	expr := "ADD downloads :inc"
	dbClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: &table,
		Key: map[string]types.AttributeValue{
			"artifact": &types.AttributeValueMemberS{Value: artifact},
		},
		UpdateExpression: &expr,
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":inc": &types.AttributeValueMemberN{Value: "1"},
		},
	})
}

func sha256Sum(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}
