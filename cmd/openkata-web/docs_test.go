package main

import (
	"errors"
	"slices"
	"strings"
	"testing"
	"testing/fstest"
)

func docsFile(s string) *fstest.MapFile { return &fstest.MapFile{Data: []byte(s)} }

func TestParseDocsFrontmatter(t *testing.T) {
	tests := []struct {
		name      string
		in        string
		wantTitle string
		wantOrder *int
		wantErr   error
	}{
		{"no frontmatter", "# Hi\n", "", nil, nil},
		{"title and order", "---\ntitle: A\norder: 3\n---\nbody", "A", intPtr(3), nil},
		{"order zero is set", "---\ntitle: A\norder: 0\n---\n", "A", intPtr(0), nil},
		{"empty block", "---\n---\nbody", "", nil, nil},
		{"unclosed block", "---\ntitle: A\n", "", nil, errDocsFrontmatter},
		{"bad yaml", "---\ntitle: [\n---\n", "", nil, errDocsFrontmatter},
		{"unknown key", "---\ntitel: A\n---\n", "", nil, errDocsFrontmatter},
		{"order not a number", "---\norder: first\n---\n", "", nil, errDocsFrontmatter},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fm, err := parseDocsFrontmatter([]byte(tt.in))
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("parseDocsFrontmatter(%q) error = %v, want %v", tt.in, err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if fm.Title != tt.wantTitle {
				t.Errorf("parseDocsFrontmatter(%q).Title = %q, want %q", tt.in, fm.Title, tt.wantTitle)
			}
			switch {
			case (fm.Order == nil) != (tt.wantOrder == nil):
				t.Errorf("parseDocsFrontmatter(%q).Order = %v, want %v", tt.in, fm.Order, tt.wantOrder)
			case fm.Order != nil && *fm.Order != *tt.wantOrder:
				t.Errorf("parseDocsFrontmatter(%q).Order = %d, want %d", tt.in, *fm.Order, *tt.wantOrder)
			}
		})
	}
}

func intPtr(n int) *int { return &n }

func TestFirstHeading(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"plain", "# Title\n\ntext", "Title"},
		{"after text", "intro\n\n# Title", "Title"},
		{"skips h2", "## Sub\n# Title", "Title"},
		{"skips fenced comment", "```\n# not a title\n```\n# Real", "Real"},
		{"none", "just text", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := firstHeading([]byte(tt.in)); got != tt.want {
				t.Errorf("firstHeading(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestBuildSite(t *testing.T) {
	base := func() fstest.MapFS {
		return fstest.MapFS{
			"index.md":        docsFile("---\ntitle: Docs\n---\n"),
			"guides/index.md": docsFile("---\ntitle: Guides\norder: 1\n---\n"),
		}
	}
	with := func(extra map[string]string) fstest.MapFS {
		fsys := base()
		for k, v := range extra {
			fsys[k] = docsFile(v)
		}
		return fsys
	}

	errTests := []struct {
		name    string
		fsys    fstest.MapFS
		wantErr error
	}{
		{"bad frontmatter in page", with(map[string]string{"guides/a.md": "---\ntitle: [\n---\n"}), errDocsFrontmatter},
		{"page without title", with(map[string]string{"guides/a.md": "text only"}), errDocsFrontmatter},
		{"section without index", with(map[string]string{"reference/a.md": "# A"}), errDocsNoIndex},
		{"top-level page", with(map[string]string{"extra.md": "# Extra"}), errDocsLayout},
		{"nested directory", with(map[string]string{"guides/sub/a.md": "# A"}), errDocsLayout},
		{"uppercase file name", with(map[string]string{"guides/Intro.md": "# Intro"}), errDocsLayout},
		{"space in file name", with(map[string]string{"guides/my page.md": "# Page"}), errDocsLayout},
		{"uppercase section", with(map[string]string{"Guides2/index.md": "# G"}), errDocsLayout},
	}
	for _, tt := range errTests {
		t.Run("error/"+tt.name, func(t *testing.T) {
			_, err := buildSite(tt.fsys)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("buildSite() error = %v, want %v", err, tt.wantErr)
			}
		})
	}

	t.Run("missing root index", func(t *testing.T) {
		fsys := base()
		delete(fsys, "index.md")
		if _, err := buildSite(fsys); err == nil {
			t.Error("buildSite() without index.md succeeded, want error")
		}
	})

	t.Run("urls", func(t *testing.T) {
		site, err := buildSite(with(map[string]string{"guides/intro.md": "# Intro"}))
		if err != nil {
			t.Fatalf("buildSite() error = %v", err)
		}
		for _, url := range []string{"/docs/", "/docs/guides/", "/docs/guides/intro/"} {
			if site.byURL[url] == nil {
				t.Errorf("buildSite() has no page at %s", url)
			}
		}
		if got := site.byURL["/docs/guides/intro/"].Title; got != "Intro" {
			t.Errorf("title of /docs/guides/intro/ = %q, want heading fallback %q", got, "Intro")
		}
	})

	t.Run("error names the file", func(t *testing.T) {
		_, err := buildSite(with(map[string]string{"guides/broken.md": "---\ntitle: [\n---\n"}))
		if err == nil || !strings.Contains(err.Error(), "guides/broken.md") {
			t.Errorf("buildSite() error = %v, want it to name guides/broken.md", err)
		}
	})
}

func TestBuildSiteOrdering(t *testing.T) {
	fsys := fstest.MapFS{
		"index.md":       docsFile("# Docs"),
		"zeta/index.md":  docsFile("---\ntitle: Zeta\norder: 1\n---\n"),
		"alpha/index.md": docsFile("---\ntitle: Alpha\norder: 2\n---\n"),
		"beta/index.md":  docsFile("# Beta"),
		"aaa/index.md":   docsFile("# Aaa"),
		"zeta/c.md":      docsFile("---\ntitle: C\n---\n"),
		"zeta/b.md":      docsFile("---\ntitle: B\n---\n"),
		"zeta/two.md":    docsFile("---\ntitle: Two\norder: 2\n---\n"),
		"zeta/one.md":    docsFile("---\ntitle: One\norder: 1\n---\n"),
		"zeta/zero.md":   docsFile("---\ntitle: Zero\norder: 0\n---\n"),
		"zeta/tie-b.md":  docsFile("---\ntitle: TieB\norder: 5\n---\n"),
		"zeta/tie-a.md":  docsFile("---\ntitle: TieA\norder: 5\n---\n"),
	}
	site, err := buildSite(fsys)
	if err != nil {
		t.Fatalf("buildSite() error = %v", err)
	}

	var sections []string
	for _, s := range site.Sections {
		sections = append(sections, s.Slug)
	}
	wantSections := []string{"zeta", "alpha", "aaa", "beta"}
	if !slices.Equal(sections, wantSections) {
		t.Errorf("section order = %v, want %v (ordered first, then by slug)", sections, wantSections)
	}

	var pages []string
	for _, p := range site.Sections[0].Pages {
		pages = append(pages, p.Title)
	}
	wantPages := []string{"Zero", "One", "Two", "TieA", "TieB", "B", "C"}
	if !slices.Equal(pages, wantPages) {
		t.Errorf("page order = %v, want %v (ordered first, then by file name)", pages, wantPages)
	}
}

func TestRegisterCollision(t *testing.T) {
	site := &docsSite{byURL: map[string]*docsPage{}, byFile: map[string]*docsPage{}}
	if err := site.register(&docsPage{URL: "/docs/a/", name: "a/one.md"}); err != nil {
		t.Fatalf("register() first page error = %v", err)
	}
	err := site.register(&docsPage{URL: "/docs/a/", name: "a/two.md"})
	if !errors.Is(err, errDocsCollision) {
		t.Errorf("register() duplicate URL error = %v, want %v", err, errDocsCollision)
	}
}
