package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"path"
	"regexp"
	"sort"
	"strings"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

var errDocsLink = errors.New("invalid link")

// Code blocks are highlighted with CSS classes instead of inline
// colors, so docs.css can color them from the theme tokens.
var docsMarkdown = goldmark.New(
	goldmark.WithExtensions(
		extension.Table,
		highlighting.NewHighlighting(
			highlighting.WithFormatOptions(chromahtml.WithClasses(true)),
		),
	),
	goldmark.WithParserOptions(parser.WithAutoHeadingID()),
	goldmark.WithRendererOptions(
		html.WithUnsafe(),
		renderer.WithNodeRenderers(util.Prioritized(docsDiagramRenderer{}, 500)),
	),
)

// docsDiagram is an SVG file's contents, inlined in place of the
// paragraph that held only its image link. Inline SVG follows the site
// theme through the design tokens, which an <img> cannot.
type docsDiagram struct {
	ast.BaseBlock
	SVG []byte
}

var kindDocsDiagram = ast.NewNodeKind("DocsDiagram")

func (n *docsDiagram) Kind() ast.NodeKind { return kindDocsDiagram }

func (n *docsDiagram) Dump(src []byte, level int) { ast.DumpHelper(n, src, level, nil, nil) }

type docsDiagramRenderer struct{}

func (docsDiagramRenderer) RegisterFuncs(r renderer.NodeRendererFuncRegisterer) {
	r.Register(kindDocsDiagram, func(w util.BufWriter, _ []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			_, _ = w.WriteString("<figure class=\"docs-diagram\">\n")
			_, _ = w.Write(n.(*docsDiagram).SVG)
			_, _ = w.WriteString("\n</figure>\n")
		}
		return ast.WalkSkipChildren, nil
	})
}

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
	var diagrams []*ast.Image
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
		case *ast.Image:
			if isLocalSVG(string(n.Destination)) {
				diagrams = append(diagrams, n)
			}
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
	// The tree is changed after the walk, so the walk never sees a
	// node it has replaced.
	for _, img := range diagrams {
		if err := s.inlineDiagram(p, img); err != nil {
			return fmt.Errorf("%s: %w", p.name, err)
		}
	}

	var buf bytes.Buffer
	if err := docsMarkdown.Renderer().Render(&buf, p.Source, doc); err != nil {
		return fmt.Errorf("%s: render: %w", p.name, err)
	}
	out := addTargetBlankToExternalLinks(buf.String())
	out = stripFirstH1(out)
	out = addHeadingAnchors(out)
	p.Body = out
	return nil
}

// docsHeadingRe matches a rendered h2 or h3 with its id. A heading's
// content cannot hold a closing heading tag, so the first one ends it.
var docsHeadingRe = regexp.MustCompile(`(?s)<(h[23]) id="([^"]+)">(.*?)</h[23]>`)

// addHeadingAnchors ends every h2 and h3 with a "#" link to itself.
// docs.js copies the link when it is clicked; without the script it is
// an ordinary link to the section.
func addHeadingAnchors(html string) string {
	return docsHeadingRe.ReplaceAllString(html,
		`<$1 id="$2">$3 <a class="heading-anchor" href="#$2" aria-label="Copy link to this section">#</a></$1>`)
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

// isLocalSVG reports whether an image destination is a relative path to
// an .svg file, which the renderer inlines. Absolute paths and URLs are
// left as images.
func isLocalSVG(dest string) bool {
	if dest == "" || strings.HasPrefix(dest, "/") || strings.HasPrefix(dest, "#") {
		return false
	}
	u, err := url.Parse(dest)
	if err != nil || u.Scheme != "" || u.Host != "" {
		return false
	}
	return strings.HasSuffix(u.Path, ".svg")
}

// inlineDiagram replaces the paragraph that holds img with the SVG file
// it links to. The image must be alone in its paragraph, because a
// figure cannot sit inside a <p>.
func (s *docsSite) inlineDiagram(from *docsPage, img *ast.Image) error {
	dest := string(img.Destination)
	para, ok := img.Parent().(*ast.Paragraph)
	if !ok || para.ChildCount() != 1 {
		return fmt.Errorf("%w: diagram %q must be alone in its paragraph", errDocsLink, dest)
	}
	u, _ := url.Parse(dest) // isLocalSVG has parsed it
	target := path.Join(path.Dir(from.name), u.Path)
	if target == ".." || strings.HasPrefix(target, "../") {
		return fmt.Errorf("%w: diagram %q leaves web/docs", errDocsLink, dest)
	}
	data, err := fs.ReadFile(s.fsys, target)
	if err != nil {
		return fmt.Errorf("%w: diagram %q: no file at %s", errDocsLink, dest, target)
	}
	svg, err := svgElement(data)
	if err != nil {
		return fmt.Errorf("%w: diagram %s: %v", errDocsLink, target, err)
	}
	para.Parent().ReplaceChild(para.Parent(), para, &docsDiagram{SVG: svg})
	return nil
}

// svgElement returns the <svg> element of an SVG file, without the XML
// declaration that is not allowed inside an HTML page.
func svgElement(data []byte) ([]byte, error) {
	data = bytes.TrimSpace(data)
	if bytes.HasPrefix(data, []byte("<?xml")) {
		end := bytes.Index(data, []byte("?>"))
		if end == -1 {
			return nil, errors.New("unclosed XML declaration")
		}
		data = bytes.TrimSpace(data[end+2:])
	}
	if !bytes.HasPrefix(data, []byte("<svg")) {
		return nil, errors.New("file does not start with an <svg> element")
	}
	return data, nil
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
