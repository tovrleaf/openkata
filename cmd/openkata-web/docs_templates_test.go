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
			if n := strings.Count(body, "is-current"); n != 1 {
				t.Errorf("GET %s has %d current sidebar links, want exactly 1", tt.target, n)
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

// breadcrumbHTML returns the breadcrumb nav element of a page, or "".
func breadcrumbHTML(body string) string {
	start := strings.Index(body, `<nav class="docs-breadcrumb"`)
	if start == -1 {
		return ""
	}
	end := strings.Index(body[start:], "</nav>")
	if end == -1 {
		return ""
	}
	return body[start : start+end+len("</nav>")]
}

func TestDocsBreadcrumb(t *testing.T) {
	tests := []struct {
		name      string
		target    string
		wantLinks []string // hrefs of the linked crumbs, in order
		wantLast  string   // text of the current crumb
	}{
		{"page", "/docs/guides/intro/", []string{"/docs/", "/docs/guides/"}, "Intro &lt;b&gt;"},
		{"section index", "/docs/guides/", []string{"/docs/"}, "Guides"},
		{"section without pages", "/docs/reference/", []string{"/docs/"}, "Reference"},
		{"last page of a section", "/docs/guides/other/", []string{"/docs/", "/docs/guides/"}, "Other"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			crumb := breadcrumbHTML(docsBody(t, tt.target))
			if crumb == "" {
				t.Fatalf("GET %s has no breadcrumb", tt.target)
			}
			if !strings.Contains(crumb, `aria-label="Breadcrumb"`) {
				t.Errorf("breadcrumb of %s = %q, want aria-label=\"Breadcrumb\"", tt.target, crumb)
			}
			if n := strings.Count(crumb, "<a "); n != len(tt.wantLinks) {
				t.Errorf("breadcrumb of %s has %d links, want %d", tt.target, n, len(tt.wantLinks))
			}
			pos := -1
			for _, href := range tt.wantLinks {
				i := strings.Index(crumb, `<a href="`+href+`">`)
				if i == -1 {
					t.Errorf("breadcrumb of %s is missing a link to %s in %q", tt.target, href, crumb)
					continue
				}
				if i < pos {
					t.Errorf("breadcrumb of %s lists %s out of order in %q", tt.target, href, crumb)
				}
				pos = i
			}
			wantCurrent := `<span aria-current="page">` + tt.wantLast + `</span>`
			last := strings.Index(crumb, wantCurrent)
			if last == -1 {
				t.Fatalf("breadcrumb of %s = %q, want current crumb %s", tt.target, crumb, wantCurrent)
			}
			if last < pos {
				t.Errorf("breadcrumb of %s puts the current page before its parents in %q", tt.target, crumb)
			}
			if n := strings.Count(crumb, "aria-current"); n != 1 {
				t.Errorf("breadcrumb of %s has %d current items, want 1", tt.target, n)
			}
		})
	}

	t.Run("home has no breadcrumb", func(t *testing.T) {
		if got := breadcrumbHTML(docsBody(t, "/docs/")); got != "" {
			t.Errorf("home page has breadcrumb %q, want none", got)
		}
	})

	t.Run("trail sits above the title and the sidebar", func(t *testing.T) {
		body := docsBody(t, "/docs/guides/intro/")
		crumb := strings.Index(body, `class="docs-breadcrumb"`)
		sidebar := strings.Index(body, `class="docs-sidebar"`)
		title := strings.Index(body, "<h1>")
		if crumb == -1 || !(crumb < sidebar && crumb < title) {
			t.Errorf("positions breadcrumb=%d sidebar=%d h1=%d, want the breadcrumb first", crumb, sidebar, title)
		}
	})
}

func TestDocsSidebarScript(t *testing.T) {
	const script = `src="/static/js/docs.js"`
	tests := []struct {
		target string
		want   bool
	}{
		{"/docs/guides/intro/", true},
		{"/docs/guides/", true},
		{"/docs/", false},
	}
	for _, tt := range tests {
		t.Run(tt.target, func(t *testing.T) {
			if got := strings.Contains(docsBody(t, tt.target), script); got != tt.want {
				t.Errorf("GET %s loads docs.js = %v, want %v", tt.target, got, tt.want)
			}
		})
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
