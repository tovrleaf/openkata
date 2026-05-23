package main

import (
	"strings"
	"testing"
)

func TestStripFrontmatter(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "no frontmatter",
			in:   "# Hello\n\nWorld",
			want: "# Hello\n\nWorld",
		},
		{
			name: "with frontmatter",
			in:   "---\ntitle: Test\ntags: foo\n---\n# Hello\n\nWorld",
			want: "# Hello\n\nWorld",
		},
		{
			name: "unclosed frontmatter",
			in:   "---\ntitle: Test\n# Hello",
			want: "---\ntitle: Test\n# Hello",
		},
		{
			name: "empty content after frontmatter",
			in:   "---\nfoo: bar\n---\n",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := string(stripFrontmatter([]byte(tt.in)))
			if got != tt.want {
				t.Errorf("stripFrontmatter(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestStripFirstH1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "removes first h1",
			in:   "<h1>Title</h1>\n<p>Body</p>\n",
			want: "<p>Body</p>\n",
		},
		{
			name: "no h1 present",
			in:   "<h2>Sub</h2>\n<p>Body</p>\n",
			want: "<h2>Sub</h2>\n<p>Body</p>\n",
		},
		{
			name: "h1 with attributes",
			in:   "<h1 id=\"title\">Title</h1>\n<p>Rest</p>\n",
			want: "<p>Rest</p>\n",
		},
		{
			name: "only removes first h1",
			in:   "<h1>First</h1>\n<h1>Second</h1>\n",
			want: "<h1>Second</h1>\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripFirstH1(tt.in)
			if got != tt.want {
				t.Errorf("stripFirstH1(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestRenderMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		in       string
		contains []string
		excludes []string
	}{
		{
			name:     "heading",
			in:       "## Section\n\nParagraph.",
			contains: []string{"<h2>Section</h2>", "<p>Paragraph.</p>"},
		},
		{
			name:     "bold and italic",
			in:       "This is **bold** and *italic*.",
			contains: []string{"<strong>bold</strong>", "<em>italic</em>"},
		},
		{
			name:     "unordered list",
			in:       "- one\n- two\n",
			contains: []string{"<ul>", "<li>one</li>", "<li>two</li>"},
		},
		{
			name:     "code block",
			in:       "```go\nfmt.Println()\n```\n",
			contains: []string{"<pre>", "<code", "fmt.Println()"},
		},
		{
			name:     "internal link no target blank",
			in:       "[link](/page)",
			contains: []string{`href="/page"`},
			excludes: []string{`target="_blank"`},
		},
		{
			name:     "external link gets target blank",
			in:       "[ext](https://example.com)",
			contains: []string{`href="https://example.com"`, `target="_blank"`},
		},
		{
			name:     "strips frontmatter",
			in:       "---\ntitle: Test\n---\n# Title\n\nBody text.",
			contains: []string{"<p>Body text.</p>"},
			excludes: []string{"title: Test", "markdown-raw"},
		},
		{
			name:     "strips first h1",
			in:       "# Main Title\n\nContent here.",
			contains: []string{"<p>Content here.</p>"},
			excludes: []string{"<h1>Main Title</h1>"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderMarkdown([]byte(tt.in))
			for _, want := range tt.contains {
				if !strings.Contains(got, want) {
					t.Errorf("renderMarkdown(%q) missing %q\ngot: %s", tt.in, want, got)
				}
			}
			for _, exc := range tt.excludes {
				if strings.Contains(got, exc) {
					t.Errorf("renderMarkdown(%q) should not contain %q\ngot: %s", tt.in, exc, got)
				}
			}
		})
	}
}

func TestAddTargetBlankToExternalLinks(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "external https link",
			in:   `<a href="https://example.com">link</a>`,
			want: `<a href="https://example.com" target="_blank">link</a>`,
		},
		{
			name: "external http link",
			in:   `<a href="http://example.com">link</a>`,
			want: `<a href="http://example.com" target="_blank">link</a>`,
		},
		{
			name: "internal link unchanged",
			in:   `<a href="/page">link</a>`,
			want: `<a href="/page">link</a>`,
		},
		{
			name: "mixed links",
			in:   `<a href="/local">a</a> and <a href="https://ext.com">b</a>`,
			want: `<a href="/local">a</a> and <a href="https://ext.com" target="_blank">b</a>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := addTargetBlankToExternalLinks(tt.in)
			if got != tt.want {
				t.Errorf("addTargetBlankToExternalLinks(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
