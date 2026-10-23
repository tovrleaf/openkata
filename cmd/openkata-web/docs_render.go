package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"path"
	"sort"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
)

var errDocsLink = errors.New("invalid link")

var docsMarkdown = goldmark.New(
	goldmark.WithExtensions(extension.Table),
	goldmark.WithParserOptions(parser.WithAutoHeadingID()),
	goldmark.WithRendererOptions(html.WithUnsafe()),
)

// renderDocs fills Body and Headings of every page. It runs after all
// pages are registered so links can be checked against the whole site.
func (s *docsSite) renderDocs() error {
	files := make([]string, 0, len(s.byFile))
	for f := range s.byFile {
		files = append(files, f)
	}
	sort.Strings(files)
	for _, f := range files {
		if err := s.renderPage(s.byFile[f]); err != nil {
			return err
		}
	}
	return nil
}

func (s *docsSite) renderPage(p *docsPage) error {
	reader := text.NewReader(p.Source)
	doc := docsMarkdown.Parser().Parse(reader)

	var linkErr error
	err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch n := n.(type) {
		case *ast.Link:
			dest, err := s.rewriteLink(p, string(n.Destination))
			if err != nil {
				linkErr = err
				return ast.WalkStop, nil
			}
			n.Destination = []byte(dest)
		case *ast.Heading:
			if n.Level == 2 || n.Level == 3 {
				p.Headings = append(p.Headings, nodeText(n, p.Source))
			}
		}
		return ast.WalkContinue, nil
	})
	if err != nil {
		return fmt.Errorf("%s: %w", p.name, err)
	}
	if linkErr != nil {
		return fmt.Errorf("%s: %w", p.name, linkErr)
	}

	var buf bytes.Buffer
	if err := docsMarkdown.Renderer().Render(&buf, p.Source, doc); err != nil {
		return fmt.Errorf("%s: render: %w", p.name, err)
	}
	out := addTargetBlankToExternalLinks(buf.String())
	out = stripFirstH1(out)
	out = strings.ReplaceAll(out,
		`<pre><code class="language-mermaid">`,
		`<pre class="mermaid-source"><code class="language-mermaid">`)
	p.Body = out
	return nil
}

// rewriteLink turns a relative link to a docs page into its /docs/ URL.
// Absolute URLs, site-absolute paths and in-page anchors pass through.
func (s *docsSite) rewriteLink(from *docsPage, dest string) (string, error) {
	if dest == "" || strings.HasPrefix(dest, "#") || strings.HasPrefix(dest, "/") {
		return dest, nil
	}
	u, err := url.Parse(dest)
	if err != nil {
		return "", fmt.Errorf("%w: %q: %v", errDocsLink, dest, err)
	}
	if u.Scheme != "" || u.Host != "" {
		return dest, nil
	}
	if !strings.HasSuffix(u.Path, ".md") {
		return "", fmt.Errorf("%w: %q is relative but not a .md page; link files outside web/docs with an absolute GitHub URL", errDocsLink, dest)
	}

	target := path.Join(path.Dir(from.name), u.Path)
	if target == ".." || strings.HasPrefix(target, "../") {
		return "", fmt.Errorf("%w: %q leaves web/docs; use an absolute GitHub URL for files outside it", errDocsLink, dest)
	}
	page, ok := s.byFile[target]
	if !ok {
		return "", fmt.Errorf("%w: %q does not match a docs page (looked for %s)", errDocsLink, dest, target)
	}
	if u.Fragment != "" {
		return page.URL + "#" + u.Fragment, nil
	}
	return page.URL, nil
}

// nodeText returns the plain text of a node and its children.
func nodeText(n ast.Node, src []byte) string {
	var b strings.Builder
	_ = ast.Walk(n, func(c ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch c := c.(type) {
		case *ast.Text:
			b.Write(c.Segment.Value(src))
		case *ast.CodeSpan:
			for ch := c.FirstChild(); ch != nil; ch = ch.NextSibling() {
				if t, ok := ch.(*ast.Text); ok {
					b.Write(t.Segment.Value(src))
				}
			}
			return ast.WalkSkipChildren, nil
		}
		return ast.WalkContinue, nil
	})
	return strings.TrimSpace(b.String())
}
