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
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/tovrleaf/openkata/cmd/openkata-web/templates"
	"github.com/tovrleaf/openkata/internal/analytics"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

// isExcludedFile returns true if the file should be excluded from the Files tab.
func isExcludedFile(path string) bool {
	switch path {
	case "tile.json", ".tesslignore", "CHANGELOG.md", "RATIONALE.md", "references/ACKNOWLEDGMENTS.md":
		return true
	}
	if strings.HasPrefix(path, "evals/") || strings.HasPrefix(path, ".tessl-plugin/") {
		return true
	}
	return false
}

// isExcludedFromArchive returns true if the file should not be included in archive downloads.
func isExcludedFromArchive(path string) bool {
	switch path {
	case "tile.json", ".tesslignore", "references/ACKNOWLEDGMENTS.md":
		return true
	}
	if strings.HasPrefix(path, "evals/") || strings.HasPrefix(path, ".tessl-plugin/") {
		return true
	}
	return false
}

// artifactRedirects maps old artifact names to their current names.
// Used for permanent redirects after renames.
var artifactRedirects = map[string]map[string]string{
	"skills": {
		"create-pr": "github-create-pr",
	},
}

func handleRedirect(w http.ResponseWriter, r *http.Request, artifactType, name string) bool {
	if redirects, ok := artifactRedirects[artifactType]; ok {
		if newName, ok := redirects[name]; ok {
			http.Redirect(w, r, "/"+artifactType+"/"+newName+"/", http.StatusMovedPermanently)
			return true
		}
	}
	return false
}

var (
	artifactLinks     map[string]string // name -> URL path
	artifactLinksOnce sync.Once
)

// fileArtifactMap maps known file paths to the artifact that owns them.
var fileArtifactMap = map[string]string{
	"docs/context/GLOSSARY.md": "/skills/grill-with-docs/",
	"docs/context/CODEBASE.md": "/skills/grill-with-docs/",
	"docs/adr/":                "/skills/create-adr/",
	"specs/":                   "/skills/spec-workflow/",
	"spec.md":                  "/skills/spec-workflow/",
	"tasks.md":                 "/skills/spec-workflow/",
	"Makefile":                 "/skills/makefile-conventions/",
	"mk/":                      "/skills/makefile-conventions/",
	"profiles/":               "/profiles/",
	"validation-report.md":     "/profiles/spec-validator/",
}

// tabLinkMap maps file paths to same-page tab anchors.
var tabLinkMap = map[string]string{
	"references/ACKNOWLEDGMENTS.md": "#acknowledgments",
}

func buildArtifactLinks() map[string]string {
	links := make(map[string]string)
	data, err := os.ReadFile("web/static/versions.json")
	if err != nil {
		return links
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return links
	}
	for _, artifactType := range []string{"skills", "rules", "profiles"} {
		artifactData, ok := raw[artifactType]
		if !ok {
			continue
		}
		var artifacts map[string]struct{ Version string }
		if err := json.Unmarshal(artifactData, &artifacts); err != nil {
			continue
		}
		for name := range artifacts {
			links[name] = "/" + artifactType + "/" + name + "/"
		}
	}
	return links
}

func getArtifactLinks() map[string]string {
	artifactLinksOnce.Do(func() {
		artifactLinks = buildArtifactLinks()
	})
	return artifactLinks
}

// gitVersions returns released versions for an artifact by reading git tags.
func gitVersions(artifactType, name string) []string {
	pattern := artifactType + "/" + name + "/v*"
	out, err := exec.Command("git", "tag", "-l", pattern).Output()
	if err != nil {
		return nil
	}
	prefix := artifactType + "/" + name + "/v"
	var versions []string
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, prefix) {
			versions = append(versions, strings.TrimPrefix(line, prefix))
		}
	}
	sort.Sort(sort.Reverse(sort.StringSlice(versions)))
	return versions
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	templates.Home().Render(r.Context(), w)
}

func handleCatalog(w http.ResponseWriter, r *http.Request) {
	templates.Catalog().Render(r.Context(), w)
}

func handleGettingStarted(w http.ResponseWriter, r *http.Request) {
	templates.GettingStarted().Render(r.Context(), w)
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

	parts := strings.Split(strings.TrimSuffix(path, "/"), "/")

	if len(parts) >= 1 && handleRedirect(w, r, "skills", parts[0]) {
		return
	}

	// /skills/:name/archive or /skills/:name/archive/:version
	if len(parts) >= 2 && parts[1] == "archive" {
		version := ""
		if len(parts) >= 3 {
			version = parts[2]
		}
		handleArchive(w, r, "skills", parts[0], version)
		return
	}

	// /skills/:name/raw/:filepath... — latest version, redirect
	if len(parts) >= 3 && parts[1] == "raw" {
		name := parts[0]
		skill := loadSkillDetailVersion(r.Context(), name, "")
		if skill == nil {
			http.NotFound(w, r)
			return
		}
		filePath := strings.Join(parts[2:], "/")
		http.Redirect(w, r, "/skills/"+name+"/"+skill.Version+"/raw/"+filePath, http.StatusFound)
		return
	}

	// /skills/:name/:version/raw/:filepath...
	if len(parts) >= 4 && parts[2] == "raw" {
		name := parts[0]
		version := parts[1]
		filePath := strings.Join(parts[3:], "/")
		skill := loadSkillDetailVersion(r.Context(), name, version)
		if skill == nil {
			http.NotFound(w, r)
			return
		}
		content, ok := skill.FileContents[filePath]
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(content))
		return
	}

	// /skills/:name/:version — detail page for specific version
	if len(parts) == 2 {
		skill := loadSkillDetailVersion(r.Context(), parts[0], parts[1])
		if skill == nil {
			http.NotFound(w, r)
			return
		}
		setPrevNext(r.Context(), skill)
		templates.SkillDetailPage(*skill).Render(r.Context(), w)
		return
	}

	// /skills/:name — show latest version
	if len(parts) == 1 {
		skill := loadSkillDetailVersion(r.Context(), parts[0], "")
		if skill == nil {
			http.NotFound(w, r)
			return
		}
		setPrevNext(r.Context(), skill)
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
	return loadArtifactList(ctx, "skills")
}

func loadDownloadCounts(ctx context.Context) map[string]int {
	counts := make(map[string]int)
	if dbClient == nil {
		return counts
	}
	scanResp, err := dbClient.Scan(ctx, &dynamodb.ScanInput{
		TableName: &table,
	})
	if err != nil {
		return counts
	}
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
	return counts
}

func loadArtifactList(ctx context.Context, artifactType string) []templates.SkillEntry {
	data := loadVersionsJSON(ctx)
	if data == nil {
		return nil
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil
	}

	artifactData, ok := raw[artifactType]
	if !ok {
		return nil
	}

	var artifacts map[string]struct {
		Version     string `json:"version"`
		Description string `json:"description"`
		Tags        string `json:"tags"`
	}
	if err := json.Unmarshal(artifactData, &artifacts); err != nil {
		return nil
	}

	counts := loadDownloadCounts(ctx)

	var entries []templates.SkillEntry
	for name, info := range artifacts {
		if info.Version == "0.0.0" {
			continue
		}
		entries = append(entries, templates.SkillEntry{
			Name:        name,
			Version:     info.Version,
			Description: info.Description,
			Tags:        info.Tags,
			Downloads:   counts[artifactType+"/"+name],
		})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name < entries[j].Name })
	return entries
}

func setPrevNext(ctx context.Context, detail *templates.ArtifactDetail) {
	if detail == nil {
		return
	}
	list := loadArtifactList(ctx, detail.Type)
	for i, entry := range list {
		if entry.Name == detail.Name {
			if i > 0 {
				detail.Prev = list[i-1].Name
			}
			if i < len(list)-1 {
				detail.Next = list[i+1].Name
			}
			return
		}
	}
}

func loadSkillDetailVersion(ctx context.Context, name, version string) *templates.ArtifactDetail {
	// Dev mode: read from local filesystem
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") == "" {
		return loadArtifactDetailLocal("skills", name, version, "SKILL.md")
	}
	return loadArtifactDetailS3(ctx, "skills", name, version, "SKILL.md")
}

func loadArtifactDetailS3(ctx context.Context, artifactType, name, version, docFile string) *templates.ArtifactDetail {
	if s3Client == nil {
		return nil
	}

	if version == "" {
		version = resolveLatestVersion(ctx, artifactType, name)
		if version == "" {
			return nil
		}
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

	var raw map[string]json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil
	}

	artifactData, ok := raw[artifactType]
	if !ok {
		return nil
	}

	var artifacts map[string]struct {
		Version     string `json:"version"`
		Description string `json:"description"`
		Tags        string `json:"tags"`
	}
	if err := json.Unmarshal(artifactData, &artifacts); err != nil {
		return nil
	}

	info, ok := artifacts[name]
	if !ok {
		return nil
	}

	// List files in the artifact prefix
	prefix := artifactType + "/" + name + "/" + version + "/"
	listResp, err := s3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: &bucket,
		Prefix: &prefix,
	})
	if err != nil {
		return nil
	}

	detail := &templates.ArtifactDetail{
		Type:         artifactType,
		Name:         name,
		Version:      version,
		Description:  info.Description,
		Tags:         info.Tags,
		Versions:     []string{version},
		FileContents: make(map[string]string),
	}

	// Get download count
	counts := loadDownloadCounts(ctx)
	detail.Downloads = counts[artifactType+"/"+name]

	// Fetch key files and build file list
	for _, obj := range listResp.Contents {
		relPath := strings.TrimPrefix(*obj.Key, prefix)
		if relPath == "" || strings.HasSuffix(relPath, "/") {
			continue
		}

		// Fetch content for special files regardless of exclusion
		var target *string
		switch relPath {
		case docFile:
			target = &detail.Docs
		case "RATIONALE.md":
			target = &detail.Rationale
		case "CHANGELOG.md":
			target = &detail.Changelog
		case "references/ACKNOWLEDGMENTS.md":
			target = &detail.Acknowledgments
		}

		if target != nil {
			getResp, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
				Bucket: &bucket,
				Key:    obj.Key,
			})
			if err == nil {
				data, err := io.ReadAll(getResp.Body)
				getResp.Body.Close()
				if err == nil {
					if relPath == "CHANGELOG.md" {
						*target = renderMarkdown(filterChangelogByVersion(data, version), name)
					} else {
						*target = renderMarkdown(data, name)
					}
					if !isExcludedFile(relPath) {
						detail.FileContents[relPath] = string(data)
						if strings.HasSuffix(relPath, ".md") {
							detail.FileContents["__rendered__"+relPath] = renderMarkdownFull(data, name)
						}
					}
				}
			}
		}

		// Only add non-excluded files to the Files tab and FileContents
		if !isExcludedFile(relPath) {
			detail.Files = append(detail.Files, relPath)
			if target == nil {
				getResp, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
					Bucket: &bucket,
					Key:    obj.Key,
				})
				if err == nil {
					data, err := io.ReadAll(getResp.Body)
					getResp.Body.Close()
					if err == nil {
						detail.FileContents[relPath] = string(data)
						if strings.HasSuffix(relPath, ".md") {
							detail.FileContents["__rendered__"+relPath] = renderMarkdownFull(data, name)
						}
					}
				}
			}
		}
	}

	sort.Strings(detail.Files)
	return detail
}

func loadArtifactDetailLocal(artifactType, name, version, docFile string) *templates.ArtifactDetail {
	dir := artifactType + "/" + name
	mdPath := dir + "/" + docFile
	data, err := os.ReadFile(mdPath)
	if err != nil {
		return nil
	}

	// Get version, description, and tags from versions.json
	var latestVersion, description, tags string
	vjData, err := os.ReadFile("web/static/versions.json")
	if err == nil {
		var raw map[string]json.RawMessage
		if json.Unmarshal(vjData, &raw) == nil {
			if artifactData, ok := raw[artifactType]; ok {
				var artifacts map[string]struct {
					Version     string `json:"version"`
					Description string `json:"description"`
					Tags        string `json:"tags"`
				}
				if json.Unmarshal(artifactData, &artifacts) == nil {
					if info, ok := artifacts[name]; ok {
						latestVersion = info.Version
						description = info.Description
						tags = info.Tags
					}
				}
			}
		}
	}

	allVersions := gitVersions(artifactType, name)

	// If no version specified, use latest
	if version == "" {
		if len(allVersions) > 0 {
			version = allVersions[0]
		} else {
			version = latestVersion
		}
	} else {
		// Validate that the requested version exists
		found := false
		for _, v := range allVersions {
			if v == version {
				found = true
				break
			}
		}
		if !found && version != latestVersion {
			return nil
		}
	}

	detail := &templates.ArtifactDetail{
		Type:         artifactType,
		Name:         name,
		Version:      version,
		Description:  description,
		Tags:         tags,
		Versions:     allVersions,
		Docs:         renderMarkdown(data, name),
		FileContents: make(map[string]string),
	}

	// Changelog
	if cl, err := os.ReadFile(dir + "/CHANGELOG.md"); err == nil {
		detail.Changelog = renderMarkdown(filterChangelogByVersion(cl, version), name)
	}

	// Acknowledgments
	if ack, err := os.ReadFile(dir + "/references/ACKNOWLEDGMENTS.md"); err == nil {
		detail.Acknowledgments = renderMarkdown(ack, name)
	}

	// Rationale
	if rat, err := os.ReadFile(dir + "/RATIONALE.md"); err == nil {
		detail.Rationale = renderMarkdown(rat, name)
	}

	// Walk directory for file list and contents
	filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		relPath, _ := filepath.Rel(dir, path)
		relPath = filepath.ToSlash(relPath)
		if isExcludedFile(relPath) {
			return nil
		}
		detail.Files = append(detail.Files, relPath)
		if content, err := os.ReadFile(path); err == nil {
			detail.FileContents[relPath] = string(content)
			if strings.HasSuffix(relPath, ".md") {
				detail.FileContents["__rendered__"+relPath] = renderMarkdownFull(content, name)
			}
		}
		return nil
	})
	sort.Strings(detail.Files)
	return detail
}

// filterChangelogByVersion keeps only changelog sections for versions <= the given version.
// It splits raw markdown by "## " lines and compares version strings.
func filterChangelogByVersion(raw []byte, version string) []byte {
	lines := strings.Split(string(raw), "\n")
	var result []string
	inSection := false
	keep := true

	for _, line := range lines {
		if strings.HasPrefix(line, "## ") {
			inSection = true
			v := extractVersionFromHeading(line)
			keep = v != "" && v != "Unreleased" && compareVersions(v, version) <= 0
			if keep {
				result = append(result, line)
			}
		} else if inSection && keep {
			result = append(result, line)
		}
	}

	return []byte(strings.Join(result, "\n"))
}

// extractVersionFromHeading extracts a version from a changelog heading like "## [1.2.3] - 2025-05-19".
func extractVersionFromHeading(line string) string {
	// Look for [version] pattern
	start := strings.Index(line, "[")
	end := strings.Index(line, "]")
	if start >= 0 && end > start {
		return line[start+1 : end]
	}
	// Fallback: try bare version after "## "
	parts := strings.Fields(line)
	if len(parts) >= 2 {
		return strings.TrimPrefix(parts[1], "v")
	}
	return ""
}

// compareVersions compares two semver-like version strings.
// Returns -1 if a < b, 0 if a == b, 1 if a > b.
func compareVersions(a, b string) int {
	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")
	for i := 0; i < len(aParts) || i < len(bParts); i++ {
		var ai, bi int
		if i < len(aParts) {
			ai, _ = strconv.Atoi(aParts[i])
		}
		if i < len(bParts) {
			bi, _ = strconv.Atoi(bParts[i])
		}
		if ai < bi {
			return -1
		}
		if ai > bi {
			return 1
		}
	}
	return 0
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

func renderMarkdown(src []byte, self ...string) string {
	src = stripFrontmatter(src)

	md := goldmark.New(
		goldmark.WithExtensions(extension.Table),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)

	var buf bytes.Buffer
	if err := md.Convert(src, &buf); err != nil {
		return ""
	}

	selfName := ""
	if len(self) > 0 {
		selfName = self[0]
	}
	output := addTargetBlankToExternalLinks(buf.String())
	output = linkArtifactReferences(output, selfName)
	return stripFirstH1(output)
}

func renderMarkdownFull(src []byte, self ...string) string {
	src = stripFrontmatter(src)
	md := goldmark.New(
		goldmark.WithExtensions(extension.Table),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)
	var buf bytes.Buffer
	if err := md.Convert(src, &buf); err != nil {
		return ""
	}
	selfName := ""
	if len(self) > 0 {
		selfName = self[0]
	}
	output := addTargetBlankToExternalLinks(buf.String())
	return linkArtifactReferences(output, selfName)
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

// linkArtifactReferences replaces known artifact names and file paths
// with links to their detail pages. Skips content inside <code>, <pre>, and <a> tags.
func linkArtifactReferences(html, self string) string {
	links := getArtifactLinks()
	if len(links) == 0 {
		return html
	}

	// Build sorted list of names (longest first to avoid partial matches)
	names := make([]string, 0, len(links)+len(fileArtifactMap)+len(tabLinkMap))
	for name := range links {
		names = append(names, name)
	}
	for path := range fileArtifactMap {
		names = append(names, path)
	}
	for path := range tabLinkMap {
		names = append(names, path)
	}
	sort.Slice(names, func(i, j int) bool {
		return len(names[i]) > len(names[j])
	})

	// Remove self from names to avoid self-linking
	if self != "" {
		filtered := names[:0]
		for _, n := range names {
			if n != self {
				filtered = append(filtered, n)
			}
		}
		names = filtered
	}

	selfURL := ""
	if self != "" {
		if u, ok := links[self]; ok {
			selfURL = u
		}
	}

	// Process HTML, skipping tags we shouldn't modify
	var result strings.Builder
	for len(html) > 0 {
		// Check if we're at a tag we should skip entirely
		if strings.HasPrefix(html, "<pre") || strings.HasPrefix(html, "<a ") ||
			strings.HasPrefix(html, "<h1") || strings.HasPrefix(html, "<h2") ||
			strings.HasPrefix(html, "<h3") || strings.HasPrefix(html, "<h4") ||
			strings.HasPrefix(html, "<h5") || strings.HasPrefix(html, "<h6") {
			var closeTag string
			switch {
			case strings.HasPrefix(html, "<pre"):
				closeTag = "</pre>"
			case strings.HasPrefix(html, "<a "):
				closeTag = "</a>"
			case strings.HasPrefix(html, "<h1"):
				closeTag = "</h1>"
			case strings.HasPrefix(html, "<h2"):
				closeTag = "</h2>"
			case strings.HasPrefix(html, "<h3"):
				closeTag = "</h3>"
			case strings.HasPrefix(html, "<h4"):
				closeTag = "</h4>"
			case strings.HasPrefix(html, "<h5"):
				closeTag = "</h5>"
			case strings.HasPrefix(html, "<h6"):
				closeTag = "</h6>"
			}
			end := strings.Index(html, closeTag)
			if end == -1 {
				result.WriteString(html)
				break
			}
			end += len(closeTag)
			result.WriteString(html[:end])
			html = html[end:]
			continue
		}
		if strings.HasPrefix(html, "<code") {
			closeTag := "</code>"
			end := strings.Index(html, closeTag)
			if end == -1 {
				result.WriteString(html)
				break
			}
			// Extract content between <code...> and </code>
			tagEnd := strings.Index(html, ">")
			if tagEnd == -1 || tagEnd >= end {
				end += len(closeTag)
				result.WriteString(html[:end])
				html = html[end:]
				continue
			}
			openTag := html[:tagEnd+1]
			content := html[tagEnd+1 : end]
			// Check if content exactly matches an artifact name (skip self)
			links := getArtifactLinks()
			if content == self {
				result.WriteString(html[:end+len(closeTag)])
				html = html[end+len(closeTag):]
				continue
			}
			if url, ok := tabLinkMap[content]; ok {
				result.WriteString(`<a href="` + url + `" class="artifact-link">` + openTag + content + closeTag + `</a>`)
				html = html[end+len(closeTag):]
				continue
			}
			if url, ok := links[content]; ok && url != selfURL {
				result.WriteString(`<a href="` + url + `" class="artifact-link">` + openTag + content + closeTag + `</a>`)
			} else if url, ok := fileArtifactMap[content]; ok && url != selfURL {
				result.WriteString(`<a href="` + url + `" class="artifact-link">` + openTag + content + closeTag + `</a>`)
			} else {
				result.WriteString(html[:end+len(closeTag)])
			}
			html = html[end+len(closeTag):]
			continue
		}

		// Check if we're inside any HTML tag (e.g., <p>, <li>)
		if html[0] == '<' {
			end := strings.Index(html, ">")
			if end == -1 {
				result.WriteString(html)
				break
			}
			result.WriteString(html[:end+1])
			html = html[end+1:]
			continue
		}

		// Find the next tag
		nextTag := strings.Index(html, "<")
		if nextTag == -1 {
			nextTag = len(html)
		}
		// Process text segment
		segment := html[:nextTag]
		segment = linkSegment(segment, names, links, selfURL)
		result.WriteString(segment)
		html = html[nextTag:]
	}
	return result.String()
}

func linkSegment(segment string, names []string, links map[string]string, selfURL string) string {
	var replacements []string
	placeholder := "\x00LINK"

	for _, name := range names {
		idx := strings.Index(segment, name)
		if idx == -1 {
			continue
		}
		// Check word boundaries
		before := idx > 0 && isWordChar(segment[idx-1])
		after := idx+len(name) < len(segment) && isWordChar(segment[idx+len(name)])
		if before || after {
			continue
		}

		url, ok := links[name]
		if !ok {
			url, ok = tabLinkMap[name]
		}
		if !ok {
			url = fileArtifactMap[name]
		}
		if url == selfURL {
			continue
		}
		link := `<a href="` + url + `" class="artifact-link">` + name + `</a>`
		replacements = append(replacements, link)
		marker := placeholder + fmt.Sprintf("%d", len(replacements)-1) + "\x00"
		segment = segment[:idx] + marker + segment[idx+len(name):]
	}

	// Restore placeholders
	for i, link := range replacements {
		marker := placeholder + fmt.Sprintf("%d", i) + "\x00"
		segment = strings.Replace(segment, marker, link, 1)
	}
	return segment
}

func isWordChar(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_'
}

func handleRules(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/rules/")

	// /rules/ — listing
	if path == "" {
		rules := loadRulesList(r.Context())
		templates.Rules(rules).Render(r.Context(), w)
		return
	}

	parts := strings.Split(strings.TrimSuffix(path, "/"), "/")

	if len(parts) >= 1 && handleRedirect(w, r, "rules", parts[0]) {
		return
	}

	// /rules/:name/archive or /rules/:name/archive/:version
	if len(parts) >= 2 && parts[1] == "archive" {
		version := ""
		if len(parts) >= 3 {
			version = parts[2]
		}
		handleArchive(w, r, "rules", parts[0], version)
		return
	}

	// /rules/:name/raw/:filepath... — latest version, redirect
	if len(parts) >= 3 && parts[1] == "raw" {
		name := parts[0]
		rule := loadRuleDetailVersion(r.Context(), name, "")
		if rule == nil {
			http.NotFound(w, r)
			return
		}
		filePath := strings.Join(parts[2:], "/")
		http.Redirect(w, r, "/rules/"+name+"/"+rule.Version+"/raw/"+filePath, http.StatusFound)
		return
	}

	// /rules/:name/:version/raw/:filepath...
	if len(parts) >= 4 && parts[2] == "raw" {
		name := parts[0]
		version := parts[1]
		filePath := strings.Join(parts[3:], "/")
		rule := loadRuleDetailVersion(r.Context(), name, version)
		if rule == nil {
			http.NotFound(w, r)
			return
		}
		content, ok := rule.FileContents[filePath]
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(content))
		return
	}

	// /rules/:name/:version — detail for specific version
	if len(parts) == 2 {
		rule := loadRuleDetailVersion(r.Context(), parts[0], parts[1])
		if rule == nil {
			http.NotFound(w, r)
			return
		}
		setPrevNext(r.Context(), rule)
		templates.RuleDetailPage(*rule).Render(r.Context(), w)
		return
	}

	// /rules/:name — show latest version
	if len(parts) == 1 {
		rule := loadRuleDetailVersion(r.Context(), parts[0], "")
		if rule == nil {
			http.NotFound(w, r)
			return
		}
		setPrevNext(r.Context(), rule)
		templates.RuleDetailPage(*rule).Render(r.Context(), w)
		return
	}

	http.NotFound(w, r)
}

func loadRuleDetailVersion(ctx context.Context, name, version string) *templates.ArtifactDetail {
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") == "" {
		return loadArtifactDetailLocal("rules", name, version, "RULE.md")
	}
	return loadArtifactDetailS3(ctx, "rules", name, version, "RULE.md")
}

func loadRulesList(ctx context.Context) []templates.SkillEntry {
	return loadArtifactList(ctx, "rules")
}

func handleArchive(w http.ResponseWriter, r *http.Request, artifactType, name, version string) {
	// Dev mode: serve from local filesystem
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") == "" {
		handleArchiveLocal(w, r, artifactType, name, version)
		return
	}

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
		if isExcludedFromArchive(relPath) {
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

	// Record analytics event
	analytics.RecordDownload(ctx, dbClient, analytics.Event{
		Artifact: artifactType + "/" + name,
		Version:  version,
		Source:   "web",
		Client:   analytics.ParseClient(r.Header.Get("User-Agent")),
		Country:  r.Header.Get("CloudFront-Viewer-Country"),
	})
}

func handleArchiveLocal(w http.ResponseWriter, r *http.Request, artifactType, name, version string) {
	dir := artifactType + "/" + name

	if _, err := os.Stat(dir); err != nil {
		http.NotFound(w, r)
		return
	}

	if version == "" {
		data := loadVersionsJSON(r.Context())
		if data == nil {
			http.NotFound(w, r)
			return
		}
		var raw map[string]json.RawMessage
		if err := json.Unmarshal(data, &raw); err != nil {
			http.NotFound(w, r)
			return
		}
		artifactData, ok := raw[artifactType]
		if !ok {
			http.NotFound(w, r)
			return
		}
		var artifacts map[string]struct {
			Version string `json:"version"`
		}
		if err := json.Unmarshal(artifactData, &artifacts); err != nil {
			http.NotFound(w, r)
			return
		}
		info, ok := artifacts[name]
		if !ok || info.Version == "0.0.0" {
			http.NotFound(w, r)
			return
		}
		version = info.Version
	}

	type fileEntry struct {
		path string
		data []byte
	}
	var files []fileEntry

	filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		relPath, _ := filepath.Rel(dir, path)
		relPath = filepath.ToSlash(relPath)
		if isExcludedFromArchive(relPath) {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		files = append(files, fileEntry{path: relPath, data: data})
		return nil
	})

	if len(files) == 0 {
		http.NotFound(w, r)
		return
	}

	checksums := make(map[string]string)
	for _, f := range files {
		checksums[f.path] = "sha256:" + sha256Sum(f.data)
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

	tw.WriteHeader(&tar.Header{
		Name: name + "/.manifest.json",
		Size: int64(len(manifest)),
		Mode: 0644,
	})
	tw.Write([]byte(manifest))
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
		Skills   map[string]struct{ Version string } `json:"skills"`
		Rules    map[string]struct{ Version string } `json:"rules"`
		Profiles map[string]struct{ Version string } `json:"profiles"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return ""
	}

	switch artifactType {
	case "skills":
		if s, ok := v.Skills[name]; ok {
			return s.Version
		}
	case "rules":
		if r, ok := v.Rules[name]; ok {
			return r.Version
		}
	case "profiles":
		if p, ok := v.Profiles[name]; ok {
			return p.Version
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

	if len(parts) >= 1 && handleRedirect(w, r, "profiles", parts[0]) {
		return
	}

	// /profiles/:name/archive or /profiles/:name/archive/:version
	if len(parts) >= 2 && parts[1] == "archive" {
		version := ""
		if len(parts) >= 3 {
			version = parts[2]
		}
		handleArchive(w, r, "profiles", parts[0], version)
		return
	}

	// /profiles/:name/raw — redirect to latest version
	if len(parts) == 2 && parts[1] == "raw" {
		name := parts[0]
		profile := loadProfileDetailVersion(r.Context(), name, "")
		if profile == nil {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, "/profiles/"+name+"/"+profile.Version+"/raw", http.StatusFound)
		return
	}

	// /profiles/:name/:version/raw — serve raw markdown
	if len(parts) == 3 && parts[2] == "raw" {
		name := parts[0]
		version := parts[1]
		profile := loadProfileDetailVersion(r.Context(), name, version)
		if profile == nil {
			http.NotFound(w, r)
			return
		}
		content, ok := profile.FileContents[name+".md"]
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(content))
		return
	}

	// /profiles/:name/:version — detail for specific version
	if len(parts) == 2 {
		profile := loadProfileDetailVersion(r.Context(), parts[0], parts[1])
		if profile == nil {
			http.NotFound(w, r)
			return
		}
		setPrevNext(r.Context(), profile)
		templates.ProfileDetailPage(*profile).Render(r.Context(), w)
		return
	}

	// /profiles/:name — show latest version
	if len(parts) == 1 {
		profile := loadProfileDetailVersion(r.Context(), parts[0], "")
		if profile == nil {
			http.NotFound(w, r)
			return
		}
		setPrevNext(r.Context(), profile)
		templates.ProfileDetailPage(*profile).Render(r.Context(), w)
		return
	}

	http.NotFound(w, r)
}

func loadProfilesList(ctx context.Context) []templates.SkillEntry {
	return loadArtifactList(ctx, "profiles")
}

func loadProfileDetailVersion(ctx context.Context, name, version string) *templates.ArtifactDetail {
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") == "" {
		return loadArtifactDetailLocal("profiles", name, version, name+".md")
	}
	return loadArtifactDetailS3(ctx, "profiles", name, version, name+".md")
}

func handleStats(w http.ResponseWriter, r *http.Request) {
	data := templates.StatsData{}

	eventsData, err := os.ReadFile(".local/stats/download-events.json")
	if err != nil {
		data.Empty = true
		templates.Stats(data).Render(r.Context(), w)
		return
	}

	var events []templates.DownloadEvent
	if err := json.Unmarshal(eventsData, &events); err != nil {
		data.Empty = true
		templates.Stats(data).Render(r.Context(), w)
		return
	}

	data.Events = events
	data.TotalDownloads = len(events)

	artifactCounts := make(map[string]int)
	typeCounts := make(map[string]int)
	clientCounts := make(map[string]int)
	countryCounts := make(map[string]int)

	for _, ev := range events {
		artifactCounts[ev.Artifact]++
		parts := strings.SplitN(ev.Artifact, "/", 2)
		if len(parts) > 0 {
			typeCounts[parts[0]]++
		}
		if ev.Client != "" {
			clientCounts[ev.Client]++
		}
		if ev.Country != "" {
			countryCounts[ev.Country]++
		}
	}

	for name, count := range artifactCounts {
		artType := ""
		parts := strings.SplitN(name, "/", 2)
		if len(parts) > 0 {
			artType = parts[0]
		}
		data.Artifacts = append(data.Artifacts, templates.ArtifactStats{Name: name, Type: artType, Downloads: count})
	}
	sort.Slice(data.Artifacts, func(i, j int) bool { return data.Artifacts[i].Downloads > data.Artifacts[j].Downloads })

	for t, count := range typeCounts {
		data.Types = append(data.Types, templates.TypeStats{Type: t, Downloads: count})
	}
	sort.Slice(data.Types, func(i, j int) bool { return data.Types[i].Downloads > data.Types[j].Downloads })

	for client, count := range clientCounts {
		data.Clients = append(data.Clients, templates.ClientStats{Client: client, Downloads: count})
	}
	sort.Slice(data.Clients, func(i, j int) bool { return data.Clients[i].Downloads > data.Clients[j].Downloads })

	for country, count := range countryCounts {
		data.Countries = append(data.Countries, templates.CountryStats{Country: country, Downloads: count})
	}
	sort.Slice(data.Countries, func(i, j int) bool { return data.Countries[i].Downloads > data.Countries[j].Downloads })

	// Page metrics
	if metricsData, err := os.ReadFile(".local/stats/page-metrics.json"); err == nil {
		var metrics []templates.DailyMetric
		if err := json.Unmarshal(metricsData, &metrics); err == nil {
			for _, m := range metrics {
				data.PageLoads += m.Invocations
			}
		}
	}

	// Page paths
	if pathsData, err := os.ReadFile(".local/stats/page-paths.json"); err == nil {
		type rawPath struct {
			Date  string `json:"date"`
			Path  string `json:"path"`
			Count int    `json:"count"`
		}
		var rawPaths []rawPath
		if err := json.Unmarshal(pathsData, &rawPaths); err == nil {
			pathAgg := make(map[string]int)
			for _, rp := range rawPaths {
				pathAgg[rp.Path] += rp.Count
			}
			for path, count := range pathAgg {
				pType := "page"
				if strings.Contains(path, "/archive") {
					pType = "download"
				}
				data.PagePaths = append(data.PagePaths, templates.PathStats{Path: path, Type: pType, Count: count})
			}
			sort.Slice(data.PagePaths, func(i, j int) bool { return data.PagePaths[i].Count > data.PagePaths[j].Count })
		}
	}

	templates.Stats(data).Render(r.Context(), w)
}
