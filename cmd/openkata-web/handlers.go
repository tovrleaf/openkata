package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/tovrleaf/openkata/cmd/openkata-web/templates"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
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

	// /skills/:name/ — detail page
	if len(parts) == 1 {
		skill := loadSkillDetail(r.Context(), parts[0])
		if skill == nil {
			http.NotFound(w, r)
			return
		}
		templates.SkillDetailPage(*skill).Render(r.Context(), w)
		return
	}

	http.NotFound(w, r)
}

func loadVersionsJSON(ctx context.Context) []byte {
	// In dev mode, try local file first
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") == "" {
		data, err := os.ReadFile("web/static/versions.json")
		if err == nil {
			return data
		}
	}

	// Fall back to S3
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
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	return data
}

func loadSkillsList(ctx context.Context) []templates.SkillEntry {
	data := loadVersionsJSON(ctx)
	if data == nil {
		return nil
	}

	var versions struct {
		Skills map[string]struct {
			Version     string `json:"version"`
			Description string `json:"description"`
			Tags        string `json:"tags"`
		} `json:"skills"`
	}
	if err := json.Unmarshal(data, &versions); err != nil {
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
		if info.Version == "0.0.0" {
			continue
		}
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

func loadSkillDetail(ctx context.Context, name string) *templates.SkillDetail {
	// Dev mode: read from local filesystem
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") == "" {
		return loadSkillDetailLocal(name)
	}

	if s3Client == nil {
		return nil
	}

	version := resolveLatestVersion(ctx, "skills", name)
	if version == "" {
		return nil
	}

	// Get metadata from versions.json
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
		} `json:"skills"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		return nil
	}

	info, ok := versions.Skills[name]
	if !ok {
		return nil
	}

	// List files in the skill prefix
	prefix := "skills/" + name + "/" + version + "/"
	listResp, err := s3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: &bucket,
		Prefix: &prefix,
	})
	if err != nil {
		return nil
	}

	detail := &templates.SkillDetail{
		Name:        name,
		Version:     version,
		Description: info.Description,
	}

	// Get download count
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
				if artAttr.Value == "skills/"+name {
					dlAttr, ok := item["downloads"].(*types.AttributeValueMemberN)
					if ok {
						n, _ := strconv.Atoi(dlAttr.Value)
						detail.Downloads = n
					}
				}
			}
		}
	}

	// Fetch key files and build file list
	for _, obj := range listResp.Contents {
		relPath := strings.TrimPrefix(*obj.Key, prefix)
		if relPath == "" || strings.HasSuffix(relPath, "/") || relPath == "tile.json" {
			continue
		}
		detail.Files = append(detail.Files, relPath)

		// Fetch content for known markdown files
		var target *string
		switch relPath {
		case "SKILL.md":
			target = &detail.Docs
		case "CHANGELOG.md":
			target = &detail.Changelog
		case "references/ACKNOWLEDGMENTS.md":
			target = &detail.Acknowledgments
		default:
			continue
		}

		getResp, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: &bucket,
			Key:    obj.Key,
		})
		if err != nil {
			continue
		}
		data, err := io.ReadAll(getResp.Body)
		getResp.Body.Close()
		if err != nil {
			continue
		}
		*target = renderMarkdown(data)
	}

	sort.Strings(detail.Files)
	return detail
}

func loadSkillDetailLocal(name string) *templates.SkillDetail {
	dir := "skills/" + name
	mdPath := dir + "/SKILL.md"
	data, err := os.ReadFile(mdPath)
	if err != nil {
		return nil
	}

	// Get version and description from versions.json
	var version, description string
	vjData, err := os.ReadFile("web/static/versions.json")
	if err == nil {
		var vj struct {
			Skills map[string]struct {
				Version     string `json:"version"`
				Description string `json:"description"`
			} `json:"skills"`
		}
		if json.Unmarshal(vjData, &vj) == nil {
			if info, ok := vj.Skills[name]; ok {
				version = info.Version
				description = info.Description
			}
		}
	}

	detail := &templates.SkillDetail{
		Name:        name,
		Version:     version,
		Description: description,
		Docs:        renderMarkdown(data),
	}

	// Changelog
	if cl, err := os.ReadFile(dir + "/CHANGELOG.md"); err == nil {
		detail.Changelog = renderMarkdown(cl)
	}

	// Acknowledgments
	if ack, err := os.ReadFile(dir + "/references/ACKNOWLEDGMENTS.md"); err == nil {
		detail.Acknowledgments = renderMarkdown(ack)
	}

	// File list
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		if e.Name() == "tile.json" {
			continue
		}
		if e.IsDir() {
			sub, _ := os.ReadDir(dir + "/" + e.Name())
			for _, s := range sub {
				detail.Files = append(detail.Files, e.Name()+"/"+s.Name())
			}
		} else {
			detail.Files = append(detail.Files, e.Name())
		}
	}
	sort.Strings(detail.Files)
	return detail
}

func stripFrontmatter(src []byte) []byte {
	s := string(src)
	if !strings.HasPrefix(s, "---") {
		return src
	}
	end := strings.Index(s[3:], "\n---")
	if end == -1 {
		return src
	}
	return []byte(strings.TrimLeft(s[end+7:], "\n"))
}

func stripFirstH1(html string) string {
	// Remove the first <h1>...</h1> block from rendered output
	start := strings.Index(html, "<h1")
	if start == -1 {
		return html
	}
	end := strings.Index(html[start:], "</h1>")
	if end == -1 {
		return html
	}
	// Include the closing tag and any trailing newline
	cut := start + end + len("</h1>")
	if cut < len(html) && html[cut] == '\n' {
		cut++
	}
	return html[:start] + html[cut:]
}

func renderMarkdown(src []byte) string {
	src = stripFrontmatter(src)

	md := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)

	var buf bytes.Buffer
	if err := md.Convert(src, &buf); err != nil {
		return ""
	}

	output := addTargetBlankToExternalLinks(buf.String())
	return stripFirstH1(output)
}

func addTargetBlankToExternalLinks(s string) string {
	// Add target="_blank" to links starting with http:// or https://
	result := strings.Builder{}
	for {
		idx := strings.Index(s, "<a ")
		if idx == -1 {
			result.WriteString(s)
			break
		}
		result.WriteString(s[:idx])
		s = s[idx:]
		end := strings.Index(s, ">")
		if end == -1 {
			result.WriteString(s)
			break
		}
		tag := s[:end+1]
		if strings.Contains(tag, "href=\"http://") || strings.Contains(tag, "href=\"https://") {
			// Insert target="_blank" before closing >
			tag = s[:end] + " target=\"_blank\">"
		}
		result.WriteString(tag)
		s = s[end+1:]
	}
	return result.String()
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
	data := loadVersionsJSON(ctx)
	if data == nil {
		return nil
	}

	var versions struct {
		Rules map[string]struct {
			Version     string `json:"version"`
			Description string `json:"description"`
			Tags        string `json:"tags"`
		} `json:"rules"`
	}
	if err := json.Unmarshal(data, &versions); err != nil {
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

func handleProfiles(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/profiles/")

	if path == "" {
		profiles := loadProfilesList(r.Context())
		templates.Profiles(profiles).Render(r.Context(), w)
		return
	}

	parts := strings.Split(strings.TrimSuffix(path, "/"), "/")
	if len(parts) >= 2 && parts[1] == "archive" {
		version := ""
		if len(parts) >= 3 {
			version = parts[2]
		}
		handleArchive(w, r, "profiles", parts[0], version)
		return
	}

	http.NotFound(w, r)
}

func loadProfilesList(ctx context.Context) []templates.SkillEntry {
	data := loadVersionsJSON(ctx)
	if data == nil {
		return nil
	}

	var versions struct {
		Profiles map[string]struct {
			Version     string `json:"version"`
			Description string `json:"description"`
			Tags        string `json:"tags"`
		} `json:"profiles"`
	}
	if err := json.Unmarshal(data, &versions); err != nil {
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

	var profiles []templates.SkillEntry
	for name, info := range versions.Profiles {
		profiles = append(profiles, templates.SkillEntry{
			Name:        name,
			Version:     info.Version,
			Description: info.Description,
			Tags:        info.Tags,
			Downloads:   counts["profiles/"+name],
		})
	}
	sort.Slice(profiles, func(i, j int) bool { return profiles[i].Name < profiles[j].Name })
	return profiles
}
