package main

import (
	"errors"
	"io/fs"
	"os"
	"path"
	"regexp"
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/tovrleaf/openkata/web/docs"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

func embeddedSite(t *testing.T) *docsSite {
	t.Helper()
	site, err := loadEmbeddedSite()
	if err != nil {
		t.Fatalf("loadEmbeddedSite() error = %v", err)
	}
	return site
}

// TestEmbeddedGeneratedDir checks the directory the deploy writes the dates
// into is part of the binary even when nothing was generated. Without
// generated/.gitkeep the go:embed pattern would not compile on a clean
// checkout.
func TestEmbeddedGeneratedDir(t *testing.T) {
	if _, err := fs.Stat(docs.FS, "generated/.gitkeep"); err != nil {
		t.Errorf("fs.Stat(docs.FS, generated/.gitkeep) error = %v, want the file embedded", err)
	}
}

// TestManifestoMatchesRoot checks the manifesto in the docs is the same as
// MANIFESTO.md, after its frontmatter, so the two copies cannot drift.
func TestManifestoMatchesRoot(t *testing.T) {
	root, err := os.ReadFile("../../MANIFESTO.md")
	if err != nil {
		t.Fatalf("os.ReadFile(MANIFESTO.md) error = %v", err)
	}
	page, err := fs.ReadFile(docs.FS, "guides/manifesto.md")
	if err != nil {
		t.Fatalf("fs.ReadFile(guides/manifesto.md) error = %v", err)
	}
	_, body, found := strings.Cut(string(page), "\n---\n\n")
	if !found {
		t.Fatal("guides/manifesto.md has no frontmatter followed by a blank line")
	}
	if body != string(root) {
		t.Error("guides/manifesto.md differs from MANIFESTO.md after its frontmatter, want the same text: copy MANIFESTO.md into it")
	}
}

func TestEmbeddedDocsTree(t *testing.T) {
	site := embeddedSite(t)

	t.Run("has the four sections in order", func(t *testing.T) {
		var got []string
		for _, s := range site.Sections {
			got = append(got, s.Slug)
		}
		want := []string{"guides", "reference", "architecture", "contributing"}
		if !slices.Equal(got, want) {
			t.Errorf("embedded sections = %v, want %v", got, want)
		}
	})

	t.Run("every page has a title and a rendered body", func(t *testing.T) {
		for url, p := range site.byURL {
			if p.Title == "" {
				t.Errorf("page %s has an empty title", url)
			}
			if p.Body == "" {
				t.Errorf("page %s has an empty body", url)
			}
		}
	})

	t.Run("every section has an index and a title", func(t *testing.T) {
		for _, s := range site.Sections {
			if s.Index == nil || s.Title == "" {
				t.Errorf("section %q has index %v and title %q, want both", s.Slug, s.Index, s.Title)
			}
		}
	})

	t.Run("architecture holds the moved documents", func(t *testing.T) {
		for _, url := range []string{
			"/docs/architecture/agent-workflow/",
			"/docs/architecture/infrastructure-diagram/",
		} {
			if site.byURL[url] == nil {
				t.Errorf("embedded docs have no page at %s", url)
			}
		}
	})

	t.Run("search index lists every page except home", func(t *testing.T) {
		index := site.searchIndex()
		if want := len(site.byURL) - 1; len(index) != want {
			t.Errorf("search index has %d entries, want %d (every page but the home page)", len(index), want)
		}
		seen := map[string]bool{}
		for _, e := range index {
			if seen[e.Path] {
				t.Errorf("search index lists %s twice", e.Path)
			}
			seen[e.Path] = true
		}
	})
}

// TestEmbeddedMatchesDisk catches a page or section directory that exists
// in web/docs but is missing from the go:embed patterns in embed.go, which
// would otherwise ship without it.
func TestEmbeddedMatchesDisk(t *testing.T) {
	disk, err := buildSite(os.DirFS("../../web/docs"))
	if err != nil {
		t.Fatalf("buildSite(web/docs on disk) error = %v", err)
	}
	embedded := embeddedSite(t)

	urls := func(s *docsSite) []string {
		var out []string
		for url := range s.byURL {
			out = append(out, url)
		}
		sort.Strings(out)
		return out
	}
	if got, want := urls(embedded), urls(disk); !slices.Equal(got, want) {
		t.Errorf("embedded pages = %v, want the pages on disk %v; add the new files to the go:embed line in web/docs/embed.go", got, want)
	}
}

// TestArchitectureDocsMoved checks the move was a move: the old copies under
// docs/ must not come back, because web/docs is the only source.
func TestArchitectureDocsMoved(t *testing.T) {
	for _, old := range []string{
		"../../docs/agent-workflow.md",
		"../../docs/infrastructure-diagram.md",
	} {
		if _, err := os.Stat(old); !errors.Is(err, fs.ErrNotExist) {
			t.Errorf("os.Stat(%s) error = %v, want the file gone (moved to web/docs/architecture)", old, err)
		}
	}
}

var svgTitle = regexp.MustCompile(`(?s)<title[^>]*>(.*?)</title>`)

// TestDiagramTitlesMatchAltText checks that every diagram's <title> is the
// alt text of the image link that shows it. On the site, screen readers
// read the title, because the inlined SVG replaces the image and its alt
// text. On GitHub they read the alt text.
func TestDiagramTitlesMatchAltText(t *testing.T) {
	found := 0
	err := fs.WalkDir(docs.FS, ".", func(file string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(file, ".md") {
			return err
		}
		src, err := fs.ReadFile(docs.FS, file)
		if err != nil {
			return err
		}
		// Parse as the renderer does, so links in code blocks are skipped.
		doc := docsMarkdown.Parser().Parse(text.NewReader(src))
		var images []*ast.Image
		_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
			if img, ok := n.(*ast.Image); ok && entering && isLocalSVG(string(img.Destination)) {
				images = append(images, img)
			}
			return ast.WalkContinue, nil
		})
		for _, img := range images {
			found++
			alt := strings.Join(strings.Fields(altText(img, src)), " ")
			svgFile := path.Join(path.Dir(file), string(img.Destination))
			svg, err := fs.ReadFile(docs.FS, svgFile)
			if err != nil {
				t.Errorf("%s: read %s: %v", file, svgFile, err)
				continue
			}
			tm := svgTitle.FindSubmatch(svg)
			if tm == nil {
				t.Errorf("%s has no <title>, want one that matches the alt text in %s", svgFile, file)
				continue
			}
			if title := strings.Join(strings.Fields(string(tm[1])), " "); title != alt {
				t.Errorf("%s <title> = %q, want the alt text in %s: %q", svgFile, title, file, alt)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("WalkDir(docs.FS) error = %v", err)
	}
	if found == 0 {
		t.Error("found no diagram image links in the docs, want at least the agent workflow diagrams")
	}
}

// altText returns an image's alt text. A line break inside it counts as
// a space, as it does when the page is shown.
func altText(img *ast.Image, src []byte) string {
	var b strings.Builder
	for c := img.FirstChild(); c != nil; c = c.NextSibling() {
		if t, ok := c.(*ast.Text); ok {
			b.Write(t.Segment.Value(src))
			if t.SoftLineBreak() || t.HardLineBreak() {
				b.WriteByte(' ')
			}
		}
	}
	return b.String()
}
