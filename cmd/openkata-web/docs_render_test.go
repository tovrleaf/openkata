package main

import (
	"errors"
	"slices"
	"strings"
	"testing"
	"testing/fstest"
)

// renderFixture builds a site where architecture/page.md holds body and
// the other pages are link targets.
func renderFixture(body string) fstest.MapFS {
	return fstest.MapFS{
		"index.md":                docsFile("---\ntitle: Docs\n---\n"),
		"guides/index.md":         docsFile("---\ntitle: Guides\n---\n"),
		"guides/intro.md":         docsFile("---\ntitle: Intro\n---\n"),
		"architecture/index.md":   docsFile("---\ntitle: Architecture\n---\n"),
		"architecture/sibling.md": docsFile("---\ntitle: Sibling\n---\n"),
		"architecture/page.md":    docsFile("---\ntitle: Page\n---\n" + body),
	}
}

func renderedBody(t *testing.T, body string) *docsPage {
	t.Helper()
	site, err := buildSite(renderFixture(body))
	if err != nil {
		t.Fatalf("buildSite() error = %v", err)
	}
	return site.byURL["/docs/architecture/page/"]
}

func TestRenderLinks(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"sibling page", "[s](sibling.md)", `href="/docs/architecture/sibling/"`},
		{"other section", "[i](../guides/intro.md)", `href="/docs/guides/intro/"`},
		{"section index", "[g](../guides/index.md)", `href="/docs/guides/"`},
		{"docs home", "[h](../index.md)", `href="/docs/"`},
		{"dot slash", "[s](./sibling.md)", `href="/docs/architecture/sibling/"`},
		{"fragment kept", "[s](sibling.md#usage)", `href="/docs/architecture/sibling/#usage"`},
		{"in-page anchor", "[a](#top)", `href="#top"`},
		{"site absolute", "[a](/stats/)", `href="/stats/"`},
		{"mailto", "[m](mailto:a@b.c)", `href="mailto:a@b.c"`},
		{"external opens in new tab", "[e](https://example.com/x)", `href="https://example.com/x" target="_blank"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderedBody(t, tt.in).Body
			if !strings.Contains(got, tt.want) {
				t.Errorf("Body of %q = %q, want it to contain %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestRenderLinkErrors(t *testing.T) {
	tests := []struct {
		name string
		in   string
	}{
		{"leaves docs tree", "[a](../../docs/adr/0001-x.md)"},
		{"leaves tree sibling dir", "[a](../../README.md)"},
		{"missing page", "[a](nope.md)"},
		{"missing in other section", "[a](../guides/nope.md)"},
		{"relative non-markdown", "[a](diagram.png)"},
		{"relative directory", "[a](../guides/)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := buildSite(renderFixture(tt.in))
			if !errors.Is(err, errDocsLink) {
				t.Fatalf("buildSite() with %q error = %v, want %v", tt.in, err, errDocsLink)
			}
			if !strings.Contains(err.Error(), "architecture/page.md") {
				t.Errorf("buildSite() with %q error = %q, want it to name architecture/page.md", tt.in, err)
			}
		})
	}

	t.Run("leaving the tree points to GitHub URLs", func(t *testing.T) {
		_, err := buildSite(renderFixture("[a](../../docs/adr/0001-x.md)"))
		if err == nil || !strings.Contains(err.Error(), "absolute GitHub URL") {
			t.Errorf("buildSite() error = %v, want it to mention an absolute GitHub URL", err)
		}
	})
}

func TestRenderDocsMarkdown(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"table", "| a | b |\n|---|---|\n| 1 | 2 |\n", "<table>"},
		{"heading anchor", "## Install it\n", `<h2 id="install-it">`},
		{"unsafe html passes", "<div class=\"note\">x</div>\n", `<div class="note">x</div>`},
		{"mermaid block marked", "```mermaid\ngraph TD\n  A-->B\n```\n", `<pre class="mermaid-source"><code class="language-mermaid">graph TD`},
		{"mermaid source escaped", "```mermaid\nA-->B\n```\n", "A--&gt;B"},
		{"other code highlighted", "```go\nx := 1\n```\n", `<pre class="chroma"><code>`},
		{"highlight uses classes", "```go\nx := 1\n```\n", `<span class="mi">1</span>`},
		{"plain fence not marked", "```\nbox\n```\n", "<pre><code>box"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderedBody(t, tt.in).Body
			if !strings.Contains(got, tt.want) {
				t.Errorf("Body of %q = %q, want it to contain %q", tt.in, got, tt.want)
			}
		})
	}

	t.Run("plain fence is not marked as mermaid", func(t *testing.T) {
		got := renderedBody(t, "```\nbox\n```\n").Body
		if strings.Contains(got, "mermaid-source") {
			t.Errorf("Body = %q, want no mermaid-source class", got)
		}
	})

	t.Run("first h1 is stripped", func(t *testing.T) {
		got := renderedBody(t, "# Page title\n\ntext\n\n## Sub\n").Body
		if strings.Contains(got, "<h1") {
			t.Errorf("Body = %q, want the first h1 removed", got)
		}
		if !strings.Contains(got, "<h2") {
			t.Errorf("Body = %q, want h2 kept", got)
		}
	})

	t.Run("frontmatter is not rendered", func(t *testing.T) {
		got := renderedBody(t, "text\n").Body
		if strings.Contains(got, "title:") {
			t.Errorf("Body = %q, want no frontmatter", got)
		}
	})
}

const testSVG = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 10 10" role="img"><title>Flow</title></svg>`

// diagramFixture is renderFixture with SVG files added under
// architecture/page/, the page's own file directory.
func diagramFixture(body string, files map[string]string) fstest.MapFS {
	fsys := renderFixture(body)
	for name, data := range files {
		fsys["architecture/page/"+name] = docsFile(data)
	}
	return fsys
}

func TestRenderDiagrams(t *testing.T) {
	render := func(t *testing.T, body string, files map[string]string) string {
		t.Helper()
		site, err := buildSite(diagramFixture(body, files))
		if err != nil {
			t.Fatalf("buildSite() error = %v", err)
		}
		return site.byURL["/docs/architecture/page/"].Body
	}

	t.Run("local svg is inlined in a figure", func(t *testing.T) {
		got := render(t, "Before.\n\n![A flow](page/flow.svg)\n\nAfter.\n", map[string]string{"flow.svg": testSVG})
		want := "<figure class=\"docs-diagram\">\n" + testSVG + "\n</figure>"
		if !strings.Contains(got, want) {
			t.Errorf("Body = %q, want it to contain %q", got, want)
		}
		if strings.Contains(got, "<img") || strings.Contains(got, "<p><figure") {
			t.Errorf("Body = %q, want no <img> and no figure inside <p>", got)
		}
		if !strings.Contains(got, "<p>Before.</p>") || !strings.Contains(got, "<p>After.</p>") {
			t.Errorf("Body = %q, want the paragraphs around the diagram kept", got)
		}
	})

	t.Run("xml declaration is dropped", func(t *testing.T) {
		got := render(t, "![A flow](page/flow.svg)\n", map[string]string{
			"flow.svg": "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" + testSVG + "\n",
		})
		if strings.Contains(got, "<?xml") || !strings.Contains(got, testSVG) {
			t.Errorf("Body = %q, want the <svg> element without the XML declaration", got)
		}
	})

	t.Run("two diagrams on one page", func(t *testing.T) {
		got := render(t, "![a](page/a.svg)\n\n![b](./page/b.svg)\n", map[string]string{"a.svg": testSVG, "b.svg": testSVG})
		if n := strings.Count(got, `<figure class="docs-diagram">`); n != 2 {
			t.Errorf("Body has %d diagrams, want 2: %q", n, got)
		}
	})

	unchanged := []struct {
		name string
		in   string
		want string
	}{
		{"external svg", "![x](https://example.com/a.svg)\n", `<img src="https://example.com/a.svg" alt="x">`},
		{"site-absolute svg", "![x](/static/img/a.svg)\n", `<img src="/static/img/a.svg" alt="x">`},
		{"relative png", "![x](pic.png)\n", `<img src="pic.png" alt="x">`},
	}
	for _, tt := range unchanged {
		t.Run("unchanged/"+tt.name, func(t *testing.T) {
			got := render(t, tt.in, nil)
			if !strings.Contains(got, tt.want) || strings.Contains(got, "docs-diagram") {
				t.Errorf("Body of %q = %q, want the plain image %q", tt.in, got, tt.want)
			}
		})
	}

	errTests := []struct {
		name  string
		in    string
		files map[string]string
	}{
		{"missing file", "![x](page/nope.svg)\n", nil},
		{"leaves docs tree", "![x](../../static/img/a.svg)\n", nil},
		{"inside a sentence", "See ![x](page/flow.svg) here.\n", map[string]string{"flow.svg": testSVG}},
		{"not an svg element", "![x](page/flow.svg)\n", map[string]string{"flow.svg": "<html></html>"}},
	}
	for _, tt := range errTests {
		t.Run("error/"+tt.name, func(t *testing.T) {
			_, err := buildSite(diagramFixture(tt.in, tt.files))
			if !errors.Is(err, errDocsLink) {
				t.Fatalf("buildSite() with %q error = %v, want %v", tt.in, err, errDocsLink)
			}
			if !strings.Contains(err.Error(), "architecture/page.md") {
				t.Errorf("buildSite() with %q error = %q, want it to name architecture/page.md", tt.in, err)
			}
		})
	}
}

func TestRenderHeadingAnchors(t *testing.T) {
	got := renderedBody(t, "# Title\n\n## Install `it`\n\n### Then\n\n#### Deep\n").Body

	for _, want := range []string{
		`<h2 id="install-it">Install <code>it</code> <a class="heading-anchor" href="#install-it" aria-label="Copy link to this section">#</a></h2>`,
		`<h3 id="then">Then <a class="heading-anchor" href="#then" aria-label="Copy link to this section">#</a></h3>`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("Body = %q, want it to contain %q", got, want)
		}
	}
	if n := strings.Count(got, `class="heading-anchor"`); n != 2 {
		t.Errorf("Body has %d heading anchors, want 2: h2 and h3 only", n)
	}

	t.Run("headings text has no anchor", func(t *testing.T) {
		if h := renderedBody(t, "## Install it\n").Headings; !slices.Equal(h, []string{"Install it"}) {
			t.Errorf("Headings = %v, want [Install it]", h)
		}
	})
}

func TestRenderHeadings(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want []string
	}{
		{"h2 and h3", "# T\n## One\n### Two\n", []string{"One", "Two"}},
		{"h1 and h4 skipped", "# T\n#### Deep\n## Kept\n", []string{"Kept"}},
		{"inline code and emphasis", "## Use `go test` *now*\n", []string{"Use go test now"}},
		{"fenced hashes ignored", "```\n## not a heading\n```\n## Real\n", []string{"Real"}},
		{"none", "just text\n", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderedBody(t, tt.in).Headings
			if !slices.Equal(got, tt.want) {
				t.Errorf("Headings of %q = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}
