package main

import (
	"errors"
	"net/http"
	"strings"
	"testing"
	"testing/fstest"
)

func TestBuildSiteDates(t *testing.T) {
	t.Run("dates are applied by file", func(t *testing.T) {
		site, err := buildSite(handlerFixture(), withDates(map[string]string{
			"index.md":        "2026-01-31",
			"guides/index.md": "2026-02-01",
			"guides/intro.md": "2026-10-05",
		}))
		if err != nil {
			t.Fatalf("buildSite() error = %v", err)
		}
		for url, want := range map[string]string{
			"/docs/":              "2026-01-31",
			"/docs/guides/":       "2026-02-01",
			"/docs/guides/intro/": "2026-10-05",
			"/docs/guides/other/": "",
			"/docs/reference/":    "",
		} {
			if got := site.byURL[url].Updated; got != want {
				t.Errorf("Updated of %s = %q, want %q", url, got, want)
			}
		}
	})

	t.Run("no dates is not an error", func(t *testing.T) {
		for name, opts := range map[string][]siteOption{
			"no option": nil,
			"nil map":   {withDates(nil)},
			"empty map": {withDates(map[string]string{})},
		} {
			site, err := buildSite(handlerFixture(), opts...)
			if err != nil {
				t.Fatalf("buildSite() with %s error = %v", name, err)
			}
			if got := site.byURL["/docs/guides/intro/"].Updated; got != "" {
				t.Errorf("Updated with %s = %q, want empty", name, got)
			}
		}
	})

	t.Run("dates for missing files are ignored", func(t *testing.T) {
		_, err := buildSite(handlerFixture(), withDates(map[string]string{"guides/gone.md": "2026-10-05"}))
		if err != nil {
			t.Errorf("buildSite() with a date for a missing file error = %v, want nil", err)
		}
	})

	invalid := []string{"", "yesterday", "2026-13-01", "2026-02-30", "25.10.2026", "2026-10-5", "20261005", "2026-10-05T10:00:00Z"}
	for _, d := range invalid {
		t.Run("invalid date "+d, func(t *testing.T) {
			_, err := buildSite(handlerFixture(), withDates(map[string]string{"guides/intro.md": d}))
			if !errors.Is(err, errDocsDate) {
				t.Fatalf("buildSite() with date %q error = %v, want %v", d, err, errDocsDate)
			}
			if !strings.Contains(err.Error(), "guides/intro.md") {
				t.Errorf("buildSite() with date %q error = %q, want it to name guides/intro.md", d, err)
			}
		})
	}
}

func TestGeneratedDirIsNotASection(t *testing.T) {
	fsys := handlerFixture()
	fsys["generated/.gitkeep"] = docsFile("")
	fsys["generated/dates.json"] = docsFile(`{"guides/intro.md": "2026-10-05"}`)

	site, err := buildSite(fsys)
	if err != nil {
		t.Fatalf("buildSite() with a generated directory error = %v, want it skipped", err)
	}
	for _, s := range site.Sections {
		if s.Slug == "generated" {
			t.Error("buildSite() made a section of the generated directory")
		}
	}
	for url := range site.byURL {
		if strings.Contains(url, "generated") {
			t.Errorf("buildSite() has a page at %s from the generated directory", url)
		}
	}
}

func TestSiteFromFSDates(t *testing.T) {
	with := func(dates string) fstest.MapFS {
		fsys := handlerFixture()
		if dates != "" {
			fsys["generated/dates.json"] = docsFile(dates)
		}
		return fsys
	}

	t.Run("dates come from generated/dates.json", func(t *testing.T) {
		site, err := siteFromFS(with(`{"index.md": "2026-01-31", "guides/intro.md": "2026-10-05"}`))
		if err != nil {
			t.Fatalf("siteFromFS() error = %v", err)
		}
		if got := site.byURL["/docs/guides/intro/"].Updated; got != "2026-10-05" {
			t.Errorf("Updated of intro = %q, want 2026-10-05", got)
		}
		if got := site.byURL["/docs/"].Updated; got != "2026-01-31" {
			t.Errorf("Updated of home = %q, want 2026-01-31", got)
		}
		if got := site.byURL["/docs/guides/other/"].Updated; got != "" {
			t.Errorf("Updated of a page without an entry = %q, want empty", got)
		}
	})

	t.Run("a tree without the file has no dates", func(t *testing.T) {
		site, err := siteFromFS(with(""))
		if err != nil {
			t.Fatalf("siteFromFS() without dates.json error = %v, want nil", err)
		}
		if got := site.byURL["/docs/guides/intro/"].Updated; got != "" {
			t.Errorf("Updated = %q, want empty", got)
		}
	})

	t.Run("an empty object is fine", func(t *testing.T) {
		if _, err := siteFromFS(with(`{}`)); err != nil {
			t.Errorf("siteFromFS() with {} error = %v, want nil", err)
		}
	})

	for name, content := range map[string]string{
		"not JSON":              "not json",
		"truncated":             `{"guides/intro.md": "2026-10`,
		"wrong type":            `["2026-10-05"]`,
		"value is not a string": `{"guides/intro.md": 20261005}`,
	} {
		t.Run("malformed file: "+name, func(t *testing.T) {
			_, err := siteFromFS(with(content))
			if !errors.Is(err, errDocsDate) {
				t.Errorf("siteFromFS(%q) error = %v, want %v", content, err, errDocsDate)
			}
		})
	}

	t.Run("a bad date inside the file is an error", func(t *testing.T) {
		_, err := siteFromFS(with(`{"guides/intro.md": "05/10/2026"}`))
		if !errors.Is(err, errDocsDate) {
			t.Errorf("siteFromFS() with a bad date error = %v, want %v", err, errDocsDate)
		}
	})
}

func datedHandler(t *testing.T, dates map[string]string) http.HandlerFunc {
	t.Helper()
	site, err := buildSite(handlerFixture(), withDates(dates))
	if err != nil {
		t.Fatalf("buildSite() error = %v", err)
	}
	return docsHandler(func() (*docsSite, error) { return site, nil })
}

func TestDocsFooterKnownDate(t *testing.T) {
	tests := []struct {
		name   string
		file   string
		target string
		date   string
		want   string
	}{
		{"page", "guides/intro.md", "/docs/guides/intro/", "2026-10-05", `Last updated <time datetime="2026-10-05">5 October 2026</time>`},
		{"section index", "guides/index.md", "/docs/guides/", "2026-12-25", `Last updated <time datetime="2026-12-25">25 December 2026</time>`},
		{"home", "index.md", "/docs/", "2026-01-31", `Last updated <time datetime="2026-01-31">31 January 2026</time>`},
		{"day without leading zero", "guides/other.md", "/docs/guides/other/", "2026-03-09", `<time datetime="2026-03-09">9 March 2026</time>`},
		{"two digit day", "guides/other.md", "/docs/guides/other/", "2026-11-30", `<time datetime="2026-11-30">30 November 2026</time>`},
		{"february", "reference/index.md", "/docs/reference/", "2026-02-28", `<time datetime="2026-02-28">28 February 2026</time>`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := datedHandler(t, map[string]string{tt.file: tt.date})
			rec := serveDocs(h, http.MethodGet, tt.target)
			body := rec.Body.String()
			if !strings.Contains(body, `<footer class="docs-updated">`) {
				t.Fatalf("GET %s has no docs-updated footer", tt.target)
			}
			if !strings.Contains(body, tt.want) {
				t.Errorf("GET %s footer = %q, want it to contain %q", tt.target, footerHTML(body), tt.want)
			}
			if strings.Contains(body, "not committed yet") {
				t.Errorf("GET %s shows the placeholder although the date is known", tt.target)
			}
		})
	}
}

func TestDocsFooterUnknownDate(t *testing.T) {
	const placeholder = `Last updated: <span class="docs-updated__pending">not committed yet</span>`
	h := datedHandler(t, nil)
	for _, target := range []string{"/docs/", "/docs/guides/", "/docs/guides/intro/", "/docs/reference/"} {
		t.Run(target, func(t *testing.T) {
			body := serveDocs(h, http.MethodGet, target).Body.String()
			if !strings.Contains(footerHTML(body), placeholder) {
				t.Errorf("GET %s footer = %q, want the placeholder %q", target, footerHTML(body), placeholder)
			}
			if strings.Contains(footerHTML(body), "<time") {
				t.Errorf("GET %s footer = %q, want no <time> element without a date", target, footerHTML(body))
			}
		})
	}
}

func TestDocsFooterIsOnEveryPage(t *testing.T) {
	dated := datedHandler(t, map[string]string{"guides/intro.md": "2026-10-05"})
	undated := datedHandler(t, nil)

	for name, h := range map[string]http.HandlerFunc{"with a date": dated, "without a date": undated} {
		for _, target := range []string{"/docs/", "/docs/guides/", "/docs/guides/intro/", "/docs/guides/other/", "/docs/reference/"} {
			t.Run(name+" "+target, func(t *testing.T) {
				body := serveDocs(h, http.MethodGet, target).Body.String()
				if n := strings.Count(body, `<footer class="docs-updated">`); n != 1 {
					t.Errorf("GET %s has %d docs-updated footers, want exactly 1", target, n)
				}
			})
		}
	}

	t.Run("footer is below the previous and next links", func(t *testing.T) {
		body := serveDocs(dated, http.MethodGet, "/docs/guides/intro/").Body.String()
		pager := strings.Index(body, `class="docs-pager"`)
		footer := strings.Index(body, `class="docs-updated"`)
		end := strings.Index(body, "</article>")
		if pager == -1 || !(pager < footer && footer < end) {
			t.Errorf("positions pager=%d footer=%d article end=%d, want the footer after the pager and inside the article", pager, footer, end)
		}
	})

	t.Run("footer is below the section list on the home page", func(t *testing.T) {
		body := serveDocs(dated, http.MethodGet, "/docs/").Body.String()
		list := strings.Index(body, `class="docs-sections"`)
		footer := strings.Index(body, `class="docs-updated"`)
		if list == -1 || footer < list {
			t.Errorf("home positions section list=%d footer=%d, want the footer after the list", list, footer)
		}
	})

	t.Run("known and unknown dates use the same footer markup", func(t *testing.T) {
		withDate := footerHTML(serveDocs(dated, http.MethodGet, "/docs/guides/intro/").Body.String())
		without := footerHTML(serveDocs(undated, http.MethodGet, "/docs/guides/intro/").Body.String())
		const open = `<footer class="docs-updated">`
		if !strings.HasPrefix(withDate, open) || !strings.HasPrefix(without, open) || !strings.HasSuffix(withDate, "</footer>") || !strings.HasSuffix(without, "</footer>") {
			t.Errorf("footers differ in their wrapper: %q and %q", withDate, without)
		}
	})
}

// footerHTML returns the docs-updated footer element of a page, or "".
func footerHTML(body string) string {
	start := strings.Index(body, `<footer class="docs-updated">`)
	if start == -1 {
		return ""
	}
	end := strings.Index(body[start:], "</footer>")
	if end == -1 {
		return ""
	}
	return body[start : start+end+len("</footer>")]
}
