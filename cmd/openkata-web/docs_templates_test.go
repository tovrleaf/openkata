package main

import (
	"net/http"
	"strings"
	"testing"
)

func docsBody(t *testing.T, target string) string {
	t.Helper()
	rec := serveDocs(fixtureHandler(t), http.MethodGet, target)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET %s status = %d, want 200", target, rec.Code)
	}
	return rec.Body.String()
}

func TestDocsPageSidebar(t *testing.T) {
	body := docsBody(t, "/docs/guides/intro/")

	t.Run("lists every section and page", func(t *testing.T) {
		for _, href := range []string{
			`href="/docs/guides/"`, `href="/docs/guides/intro/"`,
			`href="/docs/guides/other/"`, `href="/docs/reference/"`,
		} {
			if !strings.Contains(body, href) {
				t.Errorf("sidebar of /docs/guides/intro/ is missing %s", href)
			}
		}
	})

	t.Run("sidebar is a details element", func(t *testing.T) {
		if !strings.Contains(body, `<details class="docs-sidebar"`) {
			t.Errorf("page = %q, want a <details class=\"docs-sidebar\"> sidebar", body)
		}
	})

	t.Run("sections come before their pages", func(t *testing.T) {
		guides := strings.Index(body, `href="/docs/guides/"`)
		intro := strings.Index(body, `href="/docs/guides/intro/"`)
		reference := strings.Index(body, `href="/docs/reference/"`)
		if !(guides < intro && intro < reference) {
			t.Errorf("link positions guides=%d intro=%d reference=%d, want section, its pages, next section", guides, intro, reference)
		}
	})

	t.Run("body is rendered HTML, not escaped", func(t *testing.T) {
		if !strings.Contains(body, `<h2 id="setup">Setup</h2>`) {
			t.Errorf("page = %q, want the rendered h2 from the Markdown body", body)
		}
	})
}

func TestDocsPageCurrentHighlight(t *testing.T) {
	tests := []struct {
		target      string
		wantCurrent string
	}{
		{"/docs/guides/", `href="/docs/guides/" class="docs-nav__section-link is-current" aria-current="page"`},
		{"/docs/guides/intro/", `href="/docs/guides/intro/" class="docs-nav__link is-current" aria-current="page"`},
		{"/docs/reference/", `href="/docs/reference/" class="docs-nav__section-link is-current" aria-current="page"`},
	}
	for _, tt := range tests {
		t.Run(tt.target, func(t *testing.T) {
			body := docsBody(t, tt.target)
			if !strings.Contains(body, tt.wantCurrent) {
				t.Errorf("GET %s: want current link %s in %q", tt.target, tt.wantCurrent, body)
			}
			if n := strings.Count(body, `aria-current="page"`); n != 1 {
				t.Errorf("GET %s has %d current links, want exactly 1", tt.target, n)
			}
		})
	}
}

func TestDocsPagePager(t *testing.T) {
	tests := []struct {
		name      string
		target    string
		wantPrev  string
		wantNext  string
		wantPager bool
	}{
		{"first page follows the section index", "/docs/guides/intro/", "/docs/guides/", "/docs/guides/other/", true},
		{"section index has only next", "/docs/guides/", "", "/docs/guides/intro/", true},
		{"last page has only prev", "/docs/guides/other/", "/docs/guides/intro/", "", true},
		{"section without pages has no pager", "/docs/reference/", "", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := docsBody(t, tt.target)
			if got := strings.Contains(body, `class="docs-pager"`); got != tt.wantPager {
				t.Errorf("GET %s has pager = %v, want %v", tt.target, got, tt.wantPager)
			}
			for _, c := range []struct{ rel, want string }{{"prev", tt.wantPrev}, {"next", tt.wantNext}} {
				marker := `rel="` + c.rel + `"`
				if c.want == "" {
					if strings.Contains(body, marker) {
						t.Errorf("GET %s has a %s link, want none", tt.target, c.rel)
					}
					continue
				}
				want := `href="` + c.want + `" class="docs-pager__link docs-pager__link--` + c.rel + `" ` + marker
				if !strings.Contains(body, want) {
					t.Errorf("GET %s: want %s link to %s in %q", tt.target, c.rel, c.want, body)
				}
			}
		})
	}

	t.Run("pager does not cross sections", func(t *testing.T) {
		body := docsBody(t, "/docs/guides/other/")
		if strings.Contains(body, `rel="next"`) {
			t.Error("last page of guides links to the next section, want no next link")
		}
	})
}

func TestDocsHomePage(t *testing.T) {
	body := docsBody(t, "/docs/")

	for _, want := range []string{
		`href="/docs/guides/"`, `href="/docs/reference/"`,
		"Learn it.", "Home text", `class="docs docs--home"`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("home page is missing %q", want)
		}
	}
	if strings.Contains(body, "docs-sidebar") {
		t.Error("home page has a sidebar, want the section list instead")
	}
	if strings.Contains(body, `aria-current`) {
		t.Error("home page marks a current link, want none")
	}
}

func TestDocsPageUsesSiteLayout(t *testing.T) {
	for _, target := range []string{"/docs/", "/docs/guides/intro/"} {
		t.Run(target, func(t *testing.T) {
			body := docsBody(t, target)
			for _, want := range []string{"<!doctype html>", `class="nav"`, "/static/css/style.css", "— Open Kata</title>"} {
				if !strings.Contains(strings.ToLower(body), strings.ToLower(want)) {
					t.Errorf("GET %s is missing %q from the site layout", target, want)
				}
			}
		})
	}

	t.Run("nav has no Docs link yet", func(t *testing.T) {
		body := docsBody(t, "/docs/guides/")
		if strings.Contains(body, `class="nav-link">Docs`) || strings.Contains(body, `class="nav-menu-link">Docs`) {
			t.Error("nav has a Docs link, want none until go-live")
		}
	})
}
