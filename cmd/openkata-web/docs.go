package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Errors returned by buildSite. They are wrapped with the offending
// file path, so callers match them with errors.Is.
var (
	errDocsFrontmatter = errors.New("invalid frontmatter")
	errDocsCollision   = errors.New("url collision")
	errDocsNoIndex     = errors.New("section has no index.md")
	errDocsLayout      = errors.New("unsupported layout")
)

var docsSlugRe = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*$`)

// docsPage is one documentation page. Source is the Markdown with the
// frontmatter removed.
type docsPage struct {
	Section     string
	URL         string
	Title       string
	Description string
	Order       int
	hasOrder    bool
	name        string
	Source      []byte
}

type docsSection struct {
	Slug     string
	Title    string
	Order    int
	hasOrder bool
	Index    *docsPage
	Pages    []*docsPage
}

type docsSite struct {
	Home     *docsPage
	Sections []*docsSection
	byURL    map[string]*docsPage
}

type docsFrontmatter struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Order       *int   `yaml:"order"`
}

// buildSite reads every Markdown file in fsys and returns the site. It
// fails on the first structural problem instead of skipping the page.
func buildSite(fsys fs.FS) (*docsSite, error) {
	site := &docsSite{byURL: map[string]*docsPage{}}

	home, err := loadDocsPage(fsys, "index.md", "", "/docs/")
	if err != nil {
		return nil, err
	}
	if err := site.register(home); err != nil {
		return nil, err
	}
	site.Home = home

	entries, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return nil, fmt.Errorf("read docs root: %w", err)
	}
	for _, e := range entries {
		switch {
		case e.IsDir():
			sec, err := loadDocsSection(fsys, site, e.Name())
			if err != nil {
				return nil, err
			}
			site.Sections = append(site.Sections, sec)
		case strings.HasSuffix(e.Name(), ".md") && e.Name() != "index.md":
			return nil, fmt.Errorf("%s: %w: pages must live in a section directory", e.Name(), errDocsLayout)
		}
	}

	sort.SliceStable(site.Sections, func(i, j int) bool {
		a, b := site.Sections[i], site.Sections[j]
		if a.hasOrder != b.hasOrder {
			return a.hasOrder
		}
		if a.hasOrder && a.Order != b.Order {
			return a.Order < b.Order
		}
		return a.Slug < b.Slug
	})
	return site, nil
}

func loadDocsSection(fsys fs.FS, site *docsSite, dir string) (*docsSection, error) {
	if !docsSlugRe.MatchString(dir) {
		return nil, fmt.Errorf("%s: %w: directory names are lowercase letters, digits and hyphens", dir, errDocsLayout)
	}
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", dir, err)
	}

	sec := &docsSection{Slug: dir}
	hasIndex := false
	for _, e := range entries {
		file := path.Join(dir, e.Name())
		if e.IsDir() {
			return nil, fmt.Errorf("%s: %w: sections cannot contain directories", file, errDocsLayout)
		}
		if !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		if e.Name() == "index.md" {
			hasIndex = true
			idx, err := loadDocsPage(fsys, file, dir, "/docs/"+dir+"/")
			if err != nil {
				return nil, err
			}
			sec.Index, sec.Title, sec.Order, sec.hasOrder = idx, idx.Title, idx.Order, idx.hasOrder
			if err := site.register(idx); err != nil {
				return nil, err
			}
			continue
		}
		slug := strings.TrimSuffix(e.Name(), ".md")
		if !docsSlugRe.MatchString(slug) {
			return nil, fmt.Errorf("%s: %w: file names are lowercase letters, digits and hyphens", file, errDocsLayout)
		}
		p, err := loadDocsPage(fsys, file, dir, "/docs/"+dir+"/"+slug+"/")
		if err != nil {
			return nil, err
		}
		if err := site.register(p); err != nil {
			return nil, err
		}
		sec.Pages = append(sec.Pages, p)
	}
	if !hasIndex {
		return nil, fmt.Errorf("%s: %w", dir, errDocsNoIndex)
	}

	sort.SliceStable(sec.Pages, func(i, j int) bool {
		a, b := sec.Pages[i], sec.Pages[j]
		if a.hasOrder != b.hasOrder {
			return a.hasOrder
		}
		if a.hasOrder && a.Order != b.Order {
			return a.Order < b.Order
		}
		return a.name < b.name
	})
	return sec, nil
}

func (s *docsSite) register(p *docsPage) error {
	if prev, ok := s.byURL[p.URL]; ok {
		return fmt.Errorf("%s: %w with %s at %s", p.name, errDocsCollision, prev.name, p.URL)
	}
	s.byURL[p.URL] = p
	return nil
}

func loadDocsPage(fsys fs.FS, file, section, url string) (*docsPage, error) {
	raw, err := fs.ReadFile(fsys, file)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", file, err)
	}
	fm, err := parseDocsFrontmatter(raw)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", file, err)
	}
	body := stripFrontmatter(raw)

	title := fm.Title
	if title == "" {
		title = firstHeading(body)
	}
	if title == "" {
		return nil, fmt.Errorf("%s: %w: no title in frontmatter or first heading", file, errDocsFrontmatter)
	}

	p := &docsPage{
		Section:     section,
		URL:         url,
		Title:       title,
		Description: fm.Description,
		name:        file,
		Source:      body,
	}
	if fm.Order != nil {
		p.Order, p.hasOrder = *fm.Order, true
	}
	return p, nil
}

// parseDocsFrontmatter reads the YAML block at the top of src. A file
// without a block has empty frontmatter; an unclosed block is an error.
func parseDocsFrontmatter(src []byte) (docsFrontmatter, error) {
	var fm docsFrontmatter
	if !bytes.HasPrefix(src, []byte("---\n")) {
		return fm, nil
	}
	rest := src[4:]
	if bytes.HasPrefix(rest, []byte("---")) {
		return fm, nil
	}
	end := bytes.Index(rest, []byte("\n---"))
	if end == -1 {
		return fm, fmt.Errorf("%w: missing closing ---", errDocsFrontmatter)
	}
	dec := yaml.NewDecoder(bytes.NewReader(rest[:end]))
	dec.KnownFields(true)
	if err := dec.Decode(&fm); err != nil && !errors.Is(err, io.EOF) {
		return fm, fmt.Errorf("%w: %v", errDocsFrontmatter, err)
	}
	return fm, nil
}

// firstHeading returns the text of the first "# " heading outside a
// fenced code block.
func firstHeading(body []byte) string {
	inFence := false
	for _, line := range strings.Split(string(body), "\n") {
		if strings.HasPrefix(line, "```") {
			inFence = !inFence
			continue
		}
		if !inFence && strings.HasPrefix(line, "# ") {
			return strings.TrimSpace(line[2:])
		}
	}
	return ""
}
