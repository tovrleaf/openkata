package assets

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
)

const styleHash = "7c98040a54" // hash("body{}")

func testFS() fstest.MapFS {
	return fstest.MapFS{
		"css/style.css": {Data: []byte("body{}")},
		"js/app.js":     {Data: []byte("alert(1)")},
	}
}

func embedded(t *testing.T, fsys fstest.MapFS) *Server {
	t.Helper()
	s, err := NewEmbedded(fsys)
	if err != nil {
		t.Fatalf("NewEmbedded() error = %v", err)
	}
	return s
}

func TestHash(t *testing.T) {
	if got := hash([]byte("body{}")); got != styleHash {
		t.Errorf("hash(%q) = %q, want %q", "body{}", got, styleHash)
	}
}

func TestURL(t *testing.T) {
	servers := map[string]*Server{
		"embedded": embedded(t, testFS()),
		"live":     NewLive(testFS()),
	}
	tests := []struct {
		name string
		file string
		want string
	}{
		{"existing file has its hash", "css/style.css", "/static/css/style.css?v=" + styleHash},
		{"missing file has no hash", "css/missing.css", "/static/css/missing.css"},
		{"directory has no hash", "css", "/static/css"},
	}
	for kind, s := range servers {
		for _, tt := range tests {
			t.Run(kind+"/"+tt.name, func(t *testing.T) {
				if got := s.URL(tt.file); got != tt.want {
					t.Errorf("URL(%q) = %q, want %q", tt.file, got, tt.want)
				}
			})
		}
	}
}

func TestPackageURL(t *testing.T) {
	old := Default
	t.Cleanup(func() { Default = old })

	Default = nil
	if got, want := URL("css/style.css"), "/static/css/style.css"; got != want {
		t.Errorf("URL() without Default = %q, want %q", got, want)
	}
	Default = embedded(t, testFS())
	if got, want := URL("css/style.css"), "/static/css/style.css?v="+styleHash; got != want {
		t.Errorf("URL() = %q, want %q", got, want)
	}
}

func TestLiveSeesEdits(t *testing.T) {
	fsys := testFS()
	s := NewLive(fsys)
	before := s.URL("css/style.css")
	fsys["css/style.css"] = &fstest.MapFile{Data: []byte("body{color:red}")}
	if after := s.URL("css/style.css"); after == before {
		t.Errorf("URL() after an edit = %q, want a new hash", after)
	}
}

func TestServeHTTP(t *testing.T) {
	tests := []struct {
		name        string
		live        bool
		target      string
		ifNoneMatch string
		wantStatus  int
		wantCache   string
	}{
		{"current hash is cached for a year", false, "/css/style.css?v=" + styleHash, "", http.StatusOK, cacheHashed},
		{"no hash is cached briefly", false, "/css/style.css", "", http.StatusOK, cacheShort},
		{"old hash is cached briefly", false, "/css/style.css?v=0000000000", "", http.StatusOK, cacheShort},
		{"matching ETag gets 304", false, "/css/style.css", `"` + styleHash + `"`, http.StatusNotModified, cacheShort},
		{"live files are always checked", true, "/css/style.css?v=" + styleHash, "", http.StatusOK, cacheLive},
		{"live matching ETag gets 304", true, "/css/style.css", `"` + styleHash + `"`, http.StatusNotModified, cacheLive},
		{"missing file", false, "/css/missing.css", "", http.StatusNotFound, ""},
		{"directory is not listed", false, "/css/", "", http.StatusNotFound, ""},
		{"root is not listed", true, "/", "", http.StatusNotFound, ""},
		{"path outside the root", true, "/../go.mod", "", http.StatusNotFound, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewLive(testFS())
			if !tt.live {
				s = embedded(t, testFS())
			}
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.URL.Path, req.URL.RawQuery, _ = strings.Cut(tt.target, "?")
			if tt.ifNoneMatch != "" {
				req.Header.Set("If-None-Match", tt.ifNoneMatch)
			}
			rec := httptest.NewRecorder()
			s.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("ServeHTTP(%s) status = %d, want %d", tt.target, rec.Code, tt.wantStatus)
			}
			if got := rec.Header().Get("Cache-Control"); got != tt.wantCache {
				t.Errorf("ServeHTTP(%s) Cache-Control = %q, want %q", tt.target, got, tt.wantCache)
			}
			if tt.wantStatus != http.StatusOK {
				return
			}
			if got, want := rec.Header().Get("ETag"), `"`+styleHash+`"`; got != want {
				t.Errorf("ServeHTTP(%s) ETag = %q, want %q", tt.target, got, want)
			}
			if got := rec.Body.String(); got != "body{}" {
				t.Errorf("ServeHTTP(%s) body = %q, want %q", tt.target, got, "body{}")
			}
		})
	}
}
