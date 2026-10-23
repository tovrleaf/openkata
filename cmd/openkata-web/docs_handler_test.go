package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"
	"testing/fstest"
)

func handlerFixture() fstest.MapFS {
	return fstest.MapFS{
		"index.md":           docsFile("---\ntitle: Docs\n---\nHome text\n"),
		"guides/index.md":    docsFile("---\ntitle: Guides\ndescription: Learn it.\norder: 1\n---\nGuides text\n"),
		"guides/intro.md":    docsFile("---\ntitle: Intro <b>\ndescription: Start here.\norder: 1\n---\n## Setup\n\n### Install\n\nIntro text\n"),
		"guides/other.md":    docsFile("---\ntitle: Other\n---\nOther text\n"),
		"reference/index.md": docsFile("---\ntitle: Reference\norder: 2\n---\nReference text\n"),
	}
}

func fixtureHandler(t *testing.T) http.HandlerFunc {
	t.Helper()
	site, err := buildSite(handlerFixture())
	if err != nil {
		t.Fatalf("buildSite() error = %v", err)
	}
	return docsHandler(func() (*docsSite, error) { return site, nil })
}

func serveDocs(h http.HandlerFunc, method, target string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(method, target, nil))
	return rec
}

func TestDocsHandlerRouting(t *testing.T) {
	tests := []struct {
		name         string
		target       string
		wantStatus   int
		wantLocation string
		wantBody     string
	}{
		{"home", "/docs/", 200, "", "Home text"},
		{"section index", "/docs/guides/", 200, "", "Guides text"},
		{"page", "/docs/guides/intro/", 200, "", "Intro text"},
		{"page title is escaped", "/docs/guides/intro/", 200, "", "Intro &lt;b&gt;"},
		{"section without slash redirects", "/docs/guides", 301, "/docs/guides/", ""},
		{"page without slash redirects", "/docs/guides/intro", 301, "/docs/guides/intro/", ""},
		{"redirect keeps query", "/docs/guides/intro?x=1", 301, "/docs/guides/intro/?x=1", ""},
		{"unknown page", "/docs/guides/nope/", 404, "", ""},
		{"unknown section", "/docs/nope/", 404, "", ""},
		{"unknown without slash is not redirected", "/docs/nope", 404, "", ""},
		{"source file name is not a page", "/docs/guides/intro.md", 404, "", ""},
		{"index file name is not a page", "/docs/index.md", 404, "", ""},
		{"deeper than a page", "/docs/guides/intro/extra/", 404, "", ""},
		{"search index with slash", "/docs/search-index.json/", 404, "", ""},
		{"traversal dot dot", "/docs/../go.mod", 404, "", ""},
		{"traversal inside section", "/docs/guides/../../go.mod", 404, "", ""},
		{"traversal encoded", "/docs/%2e%2e/%2e%2e/go.mod", 404, "", ""},
		{"traversal encoded slash", "/docs/guides%2f..%2f..%2fgo.mod", 404, "", ""},
		{"double slash", "/docs//guides/", 404, "", ""},
		{"backslash", "/docs/guides\\intro/", 404, "", ""},
		{"null byte", "/docs/guides/intro%00/", 404, "", ""},
	}
	h := fixtureHandler(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := serveDocs(h, http.MethodGet, tt.target)
			if rec.Code != tt.wantStatus {
				t.Fatalf("GET %s status = %d, want %d", tt.target, rec.Code, tt.wantStatus)
			}
			if got := rec.Header().Get("Location"); got != tt.wantLocation {
				t.Errorf("GET %s Location = %q, want %q", tt.target, got, tt.wantLocation)
			}
			if tt.wantBody != "" && !strings.Contains(rec.Body.String(), tt.wantBody) {
				t.Errorf("GET %s body = %q, want it to contain %q", tt.target, rec.Body.String(), tt.wantBody)
			}
		})
	}
}

func TestDocsHandlerMethods(t *testing.T) {
	tests := []struct {
		method     string
		wantStatus int
	}{
		{http.MethodGet, 200},
		{http.MethodHead, 200},
		{http.MethodPost, 405},
		{http.MethodPut, 405},
		{http.MethodDelete, 405},
	}
	h := fixtureHandler(t)
	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			rec := serveDocs(h, tt.method, "/docs/guides/")
			if rec.Code != tt.wantStatus {
				t.Errorf("%s /docs/guides/ status = %d, want %d", tt.method, rec.Code, tt.wantStatus)
			}
			if tt.wantStatus == 405 && rec.Header().Get("Allow") != "GET, HEAD" {
				t.Errorf("%s /docs/guides/ Allow = %q, want %q", tt.method, rec.Header().Get("Allow"), "GET, HEAD")
			}
		})
	}
}

func TestDocsHandlerLoadError(t *testing.T) {
	loadErr := errors.New("guides/broken.md: invalid frontmatter")
	h := docsHandler(func() (*docsSite, error) { return nil, loadErr })

	for _, target := range []string{"/docs/", "/docs/guides/", "/docs/search-index.json", "/docs/nope/"} {
		t.Run(target, func(t *testing.T) {
			rec := serveDocs(h, http.MethodGet, target)
			if rec.Code != http.StatusInternalServerError {
				t.Fatalf("GET %s status = %d, want 500", target, rec.Code)
			}
			if !strings.Contains(rec.Body.String(), "guides/broken.md") {
				t.Errorf("GET %s body = %q, want it to name guides/broken.md", target, rec.Body.String())
			}
		})
	}
}

func TestDocsHandlerRebuildsPerRequest(t *testing.T) {
	calls := 0
	h := docsHandler(func() (*docsSite, error) {
		calls++
		return buildSite(handlerFixture())
	})
	serveDocs(h, http.MethodGet, "/docs/")
	serveDocs(h, http.MethodGet, "/docs/guides/")
	if calls != 2 {
		t.Errorf("loader calls after 2 requests = %d, want 2", calls)
	}
}

func TestDocsSearchIndex(t *testing.T) {
	rec := serveDocs(fixtureHandler(t), http.MethodGet, "/docs/search-index.json")
	if rec.Code != 200 {
		t.Fatalf("GET search-index.json status = %d, want 200", rec.Code)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("search-index.json Content-Type = %q, want application/json", got)
	}

	var got []docsSearchEntry
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("search-index.json is not valid JSON: %v\n%s", err, rec.Body.String())
	}

	var paths []string
	for _, e := range got {
		paths = append(paths, e.Path)
	}
	wantPaths := []string{
		"/docs/guides/", "/docs/guides/intro/", "/docs/guides/other/", "/docs/reference/",
	}
	if !slices.Equal(paths, wantPaths) {
		t.Fatalf("search index paths = %v, want %v (home excluded, navigation order)", paths, wantPaths)
	}

	intro := got[1]
	want := docsSearchEntry{
		Title:       "Intro <b>",
		Path:        "/docs/guides/intro/",
		Section:     "guides",
		Description: "Start here.",
		Headings:    []string{"Setup", "Install"},
	}
	if intro.Title != want.Title || intro.Path != want.Path || intro.Section != want.Section ||
		intro.Description != want.Description || !slices.Equal(intro.Headings, want.Headings) {
		t.Errorf("search index entry = %+v, want %+v", intro, want)
	}
}

func TestDocsSearchIndexEmptyHeadings(t *testing.T) {
	rec := serveDocs(fixtureHandler(t), http.MethodGet, "/docs/search-index.json")
	if !strings.Contains(rec.Body.String(), `"headings":[]`) {
		t.Errorf("search-index.json = %s, want \"headings\":[] for pages without headings, not null", rec.Body.String())
	}
	if strings.Contains(rec.Body.String(), "null") {
		t.Errorf("search-index.json = %s, want no null values", rec.Body.String())
	}
}
