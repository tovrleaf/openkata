package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStripFrontmatter(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "no frontmatter",
			in:   "# Hello\n\nWorld",
			want: "# Hello\n\nWorld",
		},
		{
			name: "with frontmatter",
			in:   "---\ntitle: Test\ntags: foo\n---\n# Hello\n\nWorld",
			want: "# Hello\n\nWorld",
		},
		{
			name: "unclosed frontmatter",
			in:   "---\ntitle: Test\n# Hello",
			want: "---\ntitle: Test\n# Hello",
		},
		{
			name: "empty content after frontmatter",
			in:   "---\nfoo: bar\n---\n",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := string(stripFrontmatter([]byte(tt.in)))
			if got != tt.want {
				t.Errorf("stripFrontmatter(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestStripFirstH1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "removes first h1",
			in:   "<h1>Title</h1>\n<p>Body</p>\n",
			want: "<p>Body</p>\n",
		},
		{
			name: "no h1 present",
			in:   "<h2>Sub</h2>\n<p>Body</p>\n",
			want: "<h2>Sub</h2>\n<p>Body</p>\n",
		},
		{
			name: "h1 with attributes",
			in:   "<h1 id=\"title\">Title</h1>\n<p>Rest</p>\n",
			want: "<p>Rest</p>\n",
		},
		{
			name: "only removes first h1",
			in:   "<h1>First</h1>\n<h1>Second</h1>\n",
			want: "<h1>Second</h1>\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripFirstH1(tt.in)
			if got != tt.want {
				t.Errorf("stripFirstH1(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestRenderMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		in       string
		contains []string
		excludes []string
	}{
		{
			name:     "heading",
			in:       "## Section\n\nParagraph.",
			contains: []string{"<h2>Section</h2>", "<p>Paragraph.</p>"},
		},
		{
			name:     "bold and italic",
			in:       "This is **bold** and *italic*.",
			contains: []string{"<strong>bold</strong>", "<em>italic</em>"},
		},
		{
			name:     "unordered list",
			in:       "- one\n- two\n",
			contains: []string{"<ul>", "<li>one</li>", "<li>two</li>"},
		},
		{
			name:     "code block",
			in:       "```go\nfmt.Println()\n```\n",
			contains: []string{"<pre>", "<code", "fmt.Println()"},
		},
		{
			name:     "internal link no target blank",
			in:       "[link](/page)",
			contains: []string{`href="/page"`},
			excludes: []string{`target="_blank"`},
		},
		{
			name:     "external link gets target blank",
			in:       "[ext](https://example.com)",
			contains: []string{`href="https://example.com"`, `target="_blank"`},
		},
		{
			name:     "strips frontmatter",
			in:       "---\ntitle: Test\n---\n# Title\n\nBody text.",
			contains: []string{"<p>Body text.</p>"},
			excludes: []string{"title: Test", "markdown-raw"},
		},
		{
			name:     "strips first h1",
			in:       "# Main Title\n\nContent here.",
			contains: []string{"<p>Content here.</p>"},
			excludes: []string{"<h1>Main Title</h1>"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderMarkdown([]byte(tt.in))
			for _, want := range tt.contains {
				if !strings.Contains(got, want) {
					t.Errorf("renderMarkdown(%q) missing %q\ngot: %s", tt.in, want, got)
				}
			}
			for _, exc := range tt.excludes {
				if strings.Contains(got, exc) {
					t.Errorf("renderMarkdown(%q) should not contain %q\ngot: %s", tt.in, exc, got)
				}
			}
		})
	}
}

func TestAddTargetBlankToExternalLinks(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "external https link",
			in:   `<a href="https://example.com">link</a>`,
			want: `<a href="https://example.com" target="_blank">link</a>`,
		},
		{
			name: "external http link",
			in:   `<a href="http://example.com">link</a>`,
			want: `<a href="http://example.com" target="_blank">link</a>`,
		},
		{
			name: "internal link unchanged",
			in:   `<a href="/page">link</a>`,
			want: `<a href="/page">link</a>`,
		},
		{
			name: "mixed links",
			in:   `<a href="/local">a</a> and <a href="https://ext.com">b</a>`,
			want: `<a href="/local">a</a> and <a href="https://ext.com" target="_blank">b</a>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := addTargetBlankToExternalLinks(tt.in)
			if got != tt.want {
				t.Errorf("addTargetBlankToExternalLinks(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestIsExcludedFile(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"tile.json", true},
		{".tesslignore", true},
		{"CHANGELOG.md", true},
		{"references/ACKNOWLEDGMENTS.md", true},
		{"evals/test.yaml", true},
		{"evals/nested/deep.yaml", true},
		{"SKILL.md", false},
		{"assets/diagram.png", false},
		{"references/other.md", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := isExcludedFile(tt.path)
			if got != tt.want {
				t.Errorf("isExcludedFile(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestGitVersions(t *testing.T) {
	// gitVersions shells out to git, so we test parsing logic indirectly.
	// This test verifies it returns nil for a non-existent skill (no tags).
	versions := gitVersions("nonexistent-skill-xyz-12345")
	if versions != nil {
		t.Errorf("gitVersions(nonexistent) = %v, want nil", versions)
	}
}

func setupTestSkillDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	// Create skill directory with SKILL.md
	skillDir := filepath.Join(dir, "skills", "test-skill")
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("# Test Skill\n\nA test."), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(skillDir, "helper.sh"), []byte("#!/bin/bash\necho hello"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create versions.json
	versionsDir := filepath.Join(dir, "web", "static")
	if err := os.MkdirAll(versionsDir, 0755); err != nil {
		t.Fatal(err)
	}
	versionsJSON := `{
  "skills": {
    "test-skill": {
      "version": "1.2.0",
      "description": "A test skill",
      "tags": "category:test"
    }
  },
  "rules": {}
}`
	if err := os.WriteFile(filepath.Join(versionsDir, "versions.json"), []byte(versionsJSON), 0644); err != nil {
		t.Fatal(err)
	}

	return dir
}

func TestHandleSkillsRouting(t *testing.T) {
	// Save and restore working directory
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(origDir) })

	testDir := setupTestSkillDir(t)
	if err := os.Chdir(testDir); err != nil {
		t.Fatal(err)
	}

	// Create a fake git tag so gitVersions returns something
	// Since we can't easily mock git, we test with the versions.json fallback
	tests := []struct {
		name       string
		path       string
		wantCode   int
		wantHeader string // for redirects
		wantBody   string // substring match
	}{
		{
			name:     "listing",
			path:     "/skills/",
			wantCode: 200,
		},
		{
			name:     "name renders latest version",
			path:     "/skills/test-skill",
			wantCode: 200,
		},
		{
			name:     "specific version renders detail",
			path:     "/skills/test-skill/1.2.0",
			wantCode: 200,
		},
		{
			name:     "raw file serves content",
			path:     "/skills/test-skill/1.2.0/raw/helper.sh",
			wantCode: 200,
			wantBody: "#!/bin/bash",
		},
		{
			name:     "raw file not found",
			path:     "/skills/test-skill/1.2.0/raw/nonexistent.txt",
			wantCode: 404,
		},
		{
			name:     "unknown skill 404",
			path:     "/skills/no-such-skill",
			wantCode: 404,
		},
		{
			name:     "unknown version 404",
			path:     "/skills/test-skill/9.9.9",
			wantCode: 404,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			rec := httptest.NewRecorder()
			handleSkills(rec, req)

			if rec.Code != tt.wantCode {
				t.Errorf("handleSkills(%s) status = %d, want %d", tt.path, rec.Code, tt.wantCode)
			}
			if tt.wantHeader != "" {
				got := rec.Header().Get("Location")
				if got != tt.wantHeader {
					t.Errorf("handleSkills(%s) Location = %q, want %q", tt.path, got, tt.wantHeader)
				}
			}
			if tt.wantBody != "" {
				body := rec.Body.String()
				if !strings.Contains(body, tt.wantBody) {
					t.Errorf("handleSkills(%s) body missing %q, got %q", tt.path, tt.wantBody, body)
				}
			}
		})
	}
}

func TestHandleSkillsRawContentType(t *testing.T) {
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(origDir) })

	testDir := setupTestSkillDir(t)
	if err := os.Chdir(testDir); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("GET", "/skills/test-skill/1.2.0/raw/helper.sh", nil)
	rec := httptest.NewRecorder()
	handleSkills(rec, req)

	ct := rec.Header().Get("Content-Type")
	if ct != "text/plain; charset=utf-8" {
		t.Errorf("handleSkills raw Content-Type = %q, want %q", ct, "text/plain; charset=utf-8")
	}
}

func TestHandleSkillsRawLatestRedirect(t *testing.T) {
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(origDir) })

	testDir := setupTestSkillDir(t)
	if err := os.Chdir(testDir); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("GET", "/skills/test-skill/raw/helper.sh", nil)
	rec := httptest.NewRecorder()
	handleSkills(rec, req)

	if rec.Code != http.StatusFound {
		t.Errorf("handleSkills(/skills/test-skill/raw/helper.sh) status = %d, want %d", rec.Code, http.StatusFound)
	}
	loc := rec.Header().Get("Location")
	want := "/skills/test-skill/1.2.0/raw/helper.sh"
	if loc != want {
		t.Errorf("handleSkills(/skills/test-skill/raw/helper.sh) Location = %q, want %q", loc, want)
	}
}

func TestFilterChangelogByVersion(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		version string
		want    string
	}{
		{
			name: "keeps only matching and older versions",
			input: `# Changelog

## [1.3.0] - 2025-06-01

- New feature

## [1.2.0] - 2025-05-15

- Bug fix

## [1.1.0] - 2025-04-01

- Initial release
`,
			version: "1.2.0",
			want: `# Changelog

## [1.2.0] - 2025-05-15

- Bug fix

## [1.1.0] - 2025-04-01

- Initial release
`,
		},
		{
			name: "keeps all when viewing latest",
			input: `# Changelog

## [1.3.0] - 2025-06-01

- New feature

## [1.2.0] - 2025-05-15

- Bug fix
`,
			version: "1.3.0",
			want: `# Changelog

## [1.3.0] - 2025-06-01

- New feature

## [1.2.0] - 2025-05-15

- Bug fix
`,
		},
		{
			name:    "empty changelog",
			input:   "",
			version: "1.0.0",
			want:    "",
		},
		{
			name: "no matching version keeps nothing after header",
			input: `# Changelog

## [2.0.0] - 2025-07-01

- Breaking change
`,
			version: "1.0.0",
			want: `# Changelog
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := string(filterChangelogByVersion([]byte(tt.input), tt.version))
			if got != tt.want {
				t.Errorf("filterChangelogByVersion(%q, %q) =\n%q\nwant:\n%q", tt.input, tt.version, got, tt.want)
			}
		})
	}
}

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{"1.0.0", "1.0.0", 0},
		{"1.0.0", "1.0.1", -1},
		{"1.1.0", "1.0.9", 1},
		{"2.0.0", "1.9.9", 1},
		{"0.1.0", "0.2.0", -1},
	}

	for _, tt := range tests {
		t.Run(tt.a+"_vs_"+tt.b, func(t *testing.T) {
			got := compareVersions(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("compareVersions(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestRenderMarkdownFull(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "keeps first H1",
			in:   "# Title\n\nBody text",
			want: "<h1>",
		},
		{
			name: "strips frontmatter but keeps H1",
			in:   "---\nname: test\n---\n\n# Title\n\nBody",
			want: "<h1>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderMarkdownFull([]byte(tt.in))
			if !strings.Contains(got, tt.want) {
				t.Errorf("renderMarkdownFull(%q) = %q, want to contain %q", tt.in, got, tt.want)
			}
		})
	}
}
