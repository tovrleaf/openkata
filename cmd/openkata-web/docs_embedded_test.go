package main

import (
	"errors"
	"io/fs"
	"os"
	"slices"
	"sort"
	"testing"

	"github.com/tovrleaf/openkata/web/docs"
)

func embeddedSite(t *testing.T) *docsSite {
	t.Helper()
	site, err := buildSite(docs.FS)
	if err != nil {
		t.Fatalf("buildSite(docs.FS) error = %v", err)
	}
	return site
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
