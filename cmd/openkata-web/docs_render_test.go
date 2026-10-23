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
		{"other code not marked", "```go\nx := 1\n```\n", `<pre><code class="language-go">`},
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
