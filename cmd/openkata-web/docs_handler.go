package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/tovrleaf/openkata/cmd/openkata-web/templates"
)

const docsSearchIndexPath = "/docs/search-index.json"

// handleDocs serves the documentation from web/docs on disk, rebuilding
// the site on every request so edits show without a rebuild. It is only
// registered in local development.
var handleDocs = docsHandler(func() (*docsSite, error) {
	return buildSite(os.DirFS("web/docs"))
})

// docsHandler serves pages from the site returned by load. Requests are
// resolved only through the site's URL index, so no file path ever
// comes from the request.
func docsHandler(load func() (*docsSite, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.Header().Set("Allow", "GET, HEAD")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		site, err := load()
		if err != nil {
			http.Error(w, "docs failed to load: "+err.Error(), http.StatusInternalServerError)
			return
		}

		switch p := r.URL.Path; {
		case p == docsSearchIndexPath:
			writeDocsSearchIndex(w, site)
		case site.byURL[p] != nil:
			writeDocsPage(w, r, site, site.byURL[p])
		case site.byURL[p+"/"] != nil:
			target := p + "/"
			if r.URL.RawQuery != "" {
				target += "?" + r.URL.RawQuery
			}
			http.Redirect(w, r, target, http.StatusMovedPermanently)
		default:
			http.NotFound(w, r)
		}
	}
}

// writeDocsPage renders the home page or a docs page through the
// templates. It renders to a buffer first so a template error becomes a
// 500 instead of a half-written page.
func writeDocsPage(w http.ResponseWriter, r *http.Request, site *docsSite, p *docsPage) {
	var buf bytes.Buffer
	var err error
	if p == site.Home {
		err = templates.DocsHome(site.homeModel()).Render(r.Context(), &buf)
	} else {
		err = templates.DocsPage(site.pageModel(p)).Render(r.Context(), &buf)
	}
	if err != nil {
		http.Error(w, "docs render failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(buf.Bytes())
}

func (s *docsSite) homeModel() templates.DocsHomeModel {
	m := templates.DocsHomeModel{Title: s.Home.Title, Body: s.Home.Body, Updated: updatedFor(s.Home)}
	for _, sec := range s.Sections {
		m.Sections = append(m.Sections, templates.DocsHomeSection{
			DocsLink:    templates.DocsLink{Title: sec.Title, URL: sec.Index.URL},
			Description: sec.Index.Description,
		})
	}
	return m
}

func (s *docsSite) pageModel(p *docsPage) templates.DocsPageModel {
	m := templates.DocsPageModel{Title: p.Title, Body: p.Body, Crumbs: s.crumbs(p), Nav: s.navFor(p), Updated: updatedFor(p)}
	m.Prev, m.Next = s.pager(p)
	return m
}

// updatedFor formats the last-updated date of p for display, for example
// "25 October 2026". A page without a date gets the zero value, which the
// footer shows as a placeholder.
func updatedFor(p *docsPage) templates.DocsUpdated {
	if p.Updated == "" {
		return templates.DocsUpdated{}
	}
	t, err := time.Parse(docsDateLayout, p.Updated)
	if err != nil {
		return templates.DocsUpdated{}
	}
	return templates.DocsUpdated{ISO: p.Updated, Label: t.Format("2 January 2006")}
}

// crumbs returns the breadcrumb trail for p: Docs, then its section,
// then the page itself. A section index ends at the section.
func (s *docsSite) crumbs(p *docsPage) []templates.DocsLink {
	trail := []templates.DocsLink{{Title: "Docs", URL: s.Home.URL}}
	for _, sec := range s.Sections {
		if sec.Slug != p.Section {
			continue
		}
		trail = append(trail, templates.DocsLink{Title: sec.Title, URL: sec.Index.URL})
		if p != sec.Index {
			trail = append(trail, templates.DocsLink{Title: p.Title, URL: p.URL})
		}
	}
	return trail
}

// navFor lists every section and its pages, marking the current page.
func (s *docsSite) navFor(current *docsPage) []templates.DocsNavSection {
	navPage := func(p *docsPage) templates.DocsNavPage {
		return templates.DocsNavPage{
			DocsLink: templates.DocsLink{Title: p.Title, URL: p.URL},
			Current:  p == current,
		}
	}
	nav := make([]templates.DocsNavSection, 0, len(s.Sections))
	for _, sec := range s.Sections {
		ns := templates.DocsNavSection{DocsNavPage: navPage(sec.Index)}
		for _, p := range sec.Pages {
			ns.Pages = append(ns.Pages, navPage(p))
		}
		nav = append(nav, ns)
	}
	return nav
}

// pager returns the neighbours of p within its section. A section's
// index page comes first, then its pages in order.
func (s *docsSite) pager(p *docsPage) (prev, next *templates.DocsLink) {
	for _, sec := range s.Sections {
		chain := append([]*docsPage{sec.Index}, sec.Pages...)
		for i, c := range chain {
			if c != p {
				continue
			}
			link := func(q *docsPage) *templates.DocsLink {
				return &templates.DocsLink{Title: q.Title, URL: q.URL}
			}
			if i > 0 {
				prev = link(chain[i-1])
			}
			if i < len(chain)-1 {
				next = link(chain[i+1])
			}
			return prev, next
		}
	}
	return nil, nil
}

type docsSearchEntry struct {
	Title       string   `json:"title"`
	Path        string   `json:"path"`
	Section     string   `json:"section"`
	Description string   `json:"description"`
	Headings    []string `json:"headings"`
}

// searchIndex lists every section index and page, in navigation order.
// The home page is left out: it only links to the sections.
func (s *docsSite) searchIndex() []docsSearchEntry {
	entries := []docsSearchEntry{}
	add := func(p *docsPage) {
		headings := p.Headings
		if headings == nil {
			headings = []string{}
		}
		entries = append(entries, docsSearchEntry{
			Title:       p.Title,
			Path:        p.URL,
			Section:     p.Section,
			Description: p.Description,
			Headings:    headings,
		})
	}
	for _, sec := range s.Sections {
		add(sec.Index)
		for _, p := range sec.Pages {
			add(p)
		}
	}
	return entries
}

func writeDocsSearchIndex(w http.ResponseWriter, site *docsSite) {
	data, err := json.Marshal(site.searchIndex())
	if err != nil {
		http.Error(w, "search index failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
