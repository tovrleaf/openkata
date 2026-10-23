package main

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"os"
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
			writeDocsPage(w, site.byURL[p])
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

// writeDocsPage renders a bare page. Task 5 replaces it with the docs
// templates.
func writeDocsPage(w http.ResponseWriter, p *docsPage) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<!doctype html>\n<title>%s</title>\n<h1>%s</h1>\n%s",
		html.EscapeString(p.Title), html.EscapeString(p.Title), p.Body)
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
