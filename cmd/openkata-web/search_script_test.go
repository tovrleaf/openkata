package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"testing"
)

// runSearch runs web/static/js/search.js in node against a fake page, see
// testdata/search_harness.js, and returns the HTML it rendered.
func runSearch(t *testing.T, scenario map[string]any) string {
	t.Helper()
	if _, err := exec.LookPath("node"); err != nil {
		t.Skip("node is not installed")
	}
	scenario["searchJS"] = "../../web/static/js/search.js"
	in, err := json.Marshal(scenario)
	if err != nil {
		t.Fatalf("json.Marshal(scenario) error = %v", err)
	}
	cmd := exec.Command("node", "testdata/search_harness.js")
	cmd.Stdin = bytes.NewReader(in)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("node search_harness.js error = %v\n%s", err, stderr.String())
	}
	var res struct {
		HTML string `json:"html"`
	}
	if err := json.Unmarshal(out, &res); err != nil {
		t.Fatalf("harness output %q is not JSON: %v", out, err)
	}
	return res.HTML
}

func searchVersions(extra int) map[string]any {
	skills := map[string]any{
		"alpha": map[string]any{"version": "1.0.0", "description": "Review code", "tags": "category:review"},
		"beta":  map[string]any{"version": "1.0.0", "description": "Write notes", "tags": "category:docs"},
	}
	for i := 0; i < extra; i++ {
		skills[fmt.Sprintf("alpha-extra-%d", i)] = map[string]any{"version": "1.0.0", "description": "More", "tags": ""}
	}
	return map[string]any{
		"skills":   skills,
		"rules":    map[string]any{},
		"profiles": map[string]any{},
	}
}

func searchDocs() []any {
	return []any{
		map[string]any{
			"title": "Alpha guide", "path": "/docs/guides/alpha/", "section": "guides",
			"description": "How alpha works", "headings": []string{"Install alpha", "Usage"},
		},
		map[string]any{
			"title": "Release notes", "path": "/docs/contributing/release/", "section": "contributing",
			"description": "Making a release", "headings": []string{"Tagging the head"},
		},
	}
}

func TestSearchOverlayIncludesDocs(t *testing.T) {
	tests := []struct {
		name     string
		scenario map[string]any
		want     []string
		notWant  []string
	}{
		{
			name:     "a docs page next to a skill, labelled as docs",
			scenario: map[string]any{"mode": "overlay", "query": "alpha", "versions": searchVersions(0), "docs": searchDocs()},
			want: []string{
				`href="/skills/alpha/"`,
				`href="/docs/guides/alpha/"`,
				`search-overlay__item-type">docs<`,
				`search-overlay__item-type">skill<`,
			},
		},
		{
			name:     "a heading is searched",
			scenario: map[string]any{"mode": "overlay", "query": "tagging", "versions": searchVersions(0), "docs": searchDocs()},
			want:     []string{`href="/docs/contributing/release/"`},
			notWant:  []string{`/skills/`},
		},
		{
			name:     "a tag search leaves the docs out",
			scenario: map[string]any{"mode": "overlay", "query": "category:review", "versions": searchVersions(0), "docs": searchDocs()},
			want:     []string{`href="/skills/alpha/"`},
			notWant:  []string{`/docs/`},
		},
		{
			name:     "no words lists the catalog only",
			scenario: map[string]any{"mode": "overlay", "query": "", "versions": searchVersions(0), "docs": searchDocs()},
			want:     []string{`href="/skills/alpha/"`, `href="/skills/beta/"`},
			notWant:  []string{`/docs/`},
		},
		{
			name:     "a missing docs index leaves the catalog working",
			scenario: map[string]any{"mode": "overlay", "query": "alpha", "versions": searchVersions(0), "docs": nil},
			want:     []string{`href="/skills/alpha/"`},
			notWant:  []string{`/docs/`, "No results"},
		},
		{
			name:     "an unreachable docs index leaves the catalog working",
			scenario: map[string]any{"mode": "overlay", "query": "alpha", "versions": searchVersions(0), "docs": searchDocs(), "docsNetworkError": true},
			want:     []string{`href="/skills/alpha/"`},
			notWant:  []string{`/docs/`},
		},
		{
			name:     "no match at all",
			scenario: map[string]any{"mode": "overlay", "query": "zzzz", "versions": searchVersions(0), "docs": searchDocs()},
			want:     []string{"No results found."},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html := runSearch(t, tt.scenario)
			for _, want := range tt.want {
				if !strings.Contains(html, want) {
					t.Errorf("overlay for %q = %s\nwant it to contain %q", tt.scenario["query"], html, want)
				}
			}
			for _, not := range tt.notWant {
				if strings.Contains(html, not) {
					t.Errorf("overlay for %q = %s\nwant it not to contain %q", tt.scenario["query"], html, not)
				}
			}
		})
	}
}

func TestSearchOverlayListsDocsLast(t *testing.T) {
	// The docs page matches in its title, its description and a heading, so
	// it scores far above a skill that matches only in its description.
	versions := map[string]any{
		"skills": map[string]any{
			"zeta": map[string]any{"version": "1.0.0", "description": "mentions alpha once", "tags": ""},
		},
		"rules": map[string]any{
			"alpha-rule": map[string]any{"version": "1.0.0", "description": "", "tags": ""},
		},
		"profiles": map[string]any{},
	}
	docs := []any{
		map[string]any{"title": "alpha", "path": "/docs/guides/alpha/", "description": "alpha", "headings": []string{"alpha"}},
		map[string]any{"title": "Notes on alpha", "path": "/docs/guides/notes/", "description": "", "headings": []string{}},
	}
	html := runSearch(t, map[string]any{"mode": "overlay", "query": "alpha", "versions": versions, "docs": docs})

	order := []string{`/rules/alpha-rule/`, `/skills/zeta/`, `/docs/guides/alpha/`, `/docs/guides/notes/`}
	last := -1
	for _, link := range order {
		at := strings.Index(html, link)
		if at < 0 {
			t.Fatalf("overlay = %s\nwant a result for %s", html, link)
		}
		if at < last {
			t.Errorf("overlay = %s\nwant the results in the order %v: the catalog first, then the docs, each by match", html, order)
		}
		last = at
	}
}

func TestSearchOverlayDocsDoNotPushOutTheCatalog(t *testing.T) {
	// Eight skills match, so the seven places go to them even though the docs
	// page scores higher than any of them.
	html := runSearch(t, map[string]any{"mode": "overlay", "query": "alpha", "versions": searchVersions(7), "docs": searchDocs()})
	if strings.Contains(html, `/docs/`) {
		t.Errorf("overlay = %s\nwant no docs page while skills, rules and profiles fill the list", html)
	}
	if got := strings.Count(html, `search-overlay__item-type">skill<`); got != 7 {
		t.Errorf("overlay shows %d skills, want 7", got)
	}
}

func TestSearchOverlayViewAllCountsTheCatalogOnly(t *testing.T) {
	// 9 skills match "alpha" and so does a docs page that ranks first, so
	// 10 results are shown 7 at a time. The catalog page lists only the 9.
	html := runSearch(t, map[string]any{"mode": "overlay", "query": "alpha", "versions": searchVersions(8), "docs": searchDocs()})
	if !strings.Contains(html, "View all 9 results") {
		t.Errorf("overlay = %s\nwant a link to View all 9 results, the catalog matches only", html)
	}
	if strings.Contains(html, "View all 10") {
		t.Errorf("overlay = %s\nwant the docs page not counted in View all", html)
	}
}

func TestSearchViewAllNeedsMoreCatalogResults(t *testing.T) {
	// Two skills and a docs page: everything fits, and the catalog has nothing
	// more to show.
	html := runSearch(t, map[string]any{"mode": "overlay", "query": "alpha", "versions": searchVersions(0), "docs": searchDocs()})
	if strings.Contains(html, "View all") {
		t.Errorf("overlay = %s\nwant no View all link when the catalog has nothing more", html)
	}
}

func TestCatalogPageLeavesTheDocsOut(t *testing.T) {
	html := runSearch(t, map[string]any{"mode": "catalog", "query": "alpha", "versions": searchVersions(0), "docs": searchDocs()})
	if !strings.Contains(html, `href="/skills/alpha/"`) {
		t.Errorf("catalog for alpha = %s\nwant the skill", html)
	}
	for _, not := range []string{"/docs/", "Alpha guide", "catalog-item__type\">docs"} {
		if strings.Contains(html, not) {
			t.Errorf("catalog for alpha = %s\nwant no docs result (%q)", html, not)
		}
	}
}
