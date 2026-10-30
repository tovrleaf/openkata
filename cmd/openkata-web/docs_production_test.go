package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

func TestCachedSiteBuildsOnce(t *testing.T) {
	var builds int32
	load := cachedSite(func() (*docsSite, error) {
		atomic.AddInt32(&builds, 1)
		return &docsSite{}, nil
	})

	first, err := load()
	if err != nil {
		t.Fatalf("load() error = %v", err)
	}
	for i := 0; i < 3; i++ {
		got, err := load()
		if err != nil || got != first {
			t.Fatalf("load() = %p, %v, want the site built the first time %p", got, err, first)
		}
	}
	if builds != 1 {
		t.Errorf("site was built %d times, want once", builds)
	}
}

func TestCachedSiteDoesNotKeepAFailure(t *testing.T) {
	var builds int
	load := cachedSite(func() (*docsSite, error) {
		builds++
		if builds == 1 {
			return nil, errors.New("broken dates")
		}
		return &docsSite{}, nil
	})

	if _, err := load(); err == nil {
		t.Fatal("load() error = nil, want the build error")
	}
	if site, err := load(); err != nil || site == nil {
		t.Fatalf("load() = %v, %v, want a site after the failure", site, err)
	}
	if builds != 2 {
		t.Errorf("site was built %d times, want 2: the failure must not be cached", builds)
	}
}

func TestCachedSiteConcurrentRequests(t *testing.T) {
	var builds int32
	load := cachedSite(func() (*docsSite, error) {
		atomic.AddInt32(&builds, 1)
		return &docsSite{}, nil
	})

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := load(); err != nil {
				t.Errorf("load() error = %v", err)
			}
		}()
	}
	wg.Wait()
	if builds != 1 {
		t.Errorf("site was built %d times for 20 simultaneous requests, want once", builds)
	}
}

func TestDocsSiteLoaderInLambda(t *testing.T) {
	load := docsSiteLoader(true)
	first, err := load()
	if err != nil {
		t.Fatalf("docsSiteLoader(true)() error = %v", err)
	}
	if first.Home == nil || len(first.Sections) == 0 {
		t.Fatalf("embedded site has home %v and %d sections, want both", first.Home, len(first.Sections))
	}
	second, _ := load()
	if second != first {
		t.Error("docsSiteLoader(true) built the site twice, want it built once")
	}
}

// TestProductionDocsRoute serves the docs the way Lambda does: from the
// embedded tree, through the same handler main registers.
func TestProductionDocsRoute(t *testing.T) {
	handler := docsHandler(docsSiteLoader(true))
	get := func(method, path string) *httptest.ResponseRecorder {
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, httptest.NewRequest(method, path, nil))
		return rec
	}

	t.Run("home page, with the nav and footer", func(t *testing.T) {
		rec := get(http.MethodGet, "/docs/")
		if rec.Code != http.StatusOK {
			t.Fatalf("GET /docs/ = %d, want 200", rec.Code)
		}
		body := rec.Body.String()
		for _, want := range []string{`href="/docs/" class="nav-link"`, `<a href="/docs/">Documentation</a>`, "Last updated"} {
			if !strings.Contains(body, want) {
				t.Errorf("GET /docs/ body does not contain %q", want)
			}
		}
	})

	t.Run("every embedded page is served", func(t *testing.T) {
		site, err := loadEmbeddedSite()
		if err != nil {
			t.Fatalf("loadEmbeddedSite() error = %v", err)
		}
		for url := range site.byURL {
			if rec := get(http.MethodGet, url); rec.Code != http.StatusOK {
				t.Errorf("GET %s = %d, want 200", url, rec.Code)
			}
		}
	})

	t.Run("search index", func(t *testing.T) {
		rec := get(http.MethodGet, "/docs/search-index.json")
		if rec.Code != http.StatusOK || !strings.Contains(rec.Header().Get("Content-Type"), "json") {
			t.Errorf("GET /docs/search-index.json = %d %q, want 200 json", rec.Code, rec.Header().Get("Content-Type"))
		}
	})

	t.Run("unknown page and wrong method", func(t *testing.T) {
		if rec := get(http.MethodGet, "/docs/no-such-page/"); rec.Code != http.StatusNotFound {
			t.Errorf("GET /docs/no-such-page/ = %d, want 404", rec.Code)
		}
		if rec := get(http.MethodPost, "/docs/"); rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("POST /docs/ = %d, want 405", rec.Code)
		}
	})
}
