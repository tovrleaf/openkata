// Package assets serves the site's static files with cache headers, and
// builds their URLs with a content hash so a changed file gets a new URL.
package assets

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io/fs"
	"net/http"
	"strings"
	"time"
)

// Cache-Control values. A URL that carries the file's current hash never
// changes content, so it is cached for a year. Any other URL is cached for
// five minutes: files that CSS imports or scripts fetch have no hash.
const (
	cacheHashed = "public, max-age=31536000, immutable"
	cacheShort  = "public, max-age=300"
	cacheLive   = "no-cache"
)

// Server serves the files of an fs.FS and knows each file's hash.
type Server struct {
	fsys   fs.FS
	live   bool
	hashes map[string]string
}

// Default is the server that URL uses. main sets it at start-up.
var Default *Server

// NewEmbedded returns a server for files that do not change while it runs,
// such as an embed.FS. It hashes every file once.
func NewEmbedded(fsys fs.FS) (*Server, error) {
	hashes := make(map[string]string)
	err := fs.WalkDir(fsys, ".", func(name string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		data, err := fs.ReadFile(fsys, name)
		if err != nil {
			return err
		}
		hashes[name] = hash(data)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &Server{fsys: fsys, hashes: hashes}, nil
}

// NewLive returns a server for files that are edited while it runs, as in
// local development. It hashes a file each time it is asked, and tells
// browsers to check for a new version on every load.
func NewLive(fsys fs.FS) *Server {
	return &Server{fsys: fsys, live: true}
}

// URL returns the path of a static file under /static/ with its hash, using
// Default. Without Default, or for a file that does not exist, the URL has no
// hash.
func URL(name string) string {
	if Default == nil {
		return "/static/" + name
	}
	return Default.URL(name)
}

// URL returns the path of a static file under /static/ with its hash. For a
// file that does not exist, the URL has no hash.
func (s *Server) URL(name string) string {
	if h, ok := s.hash(name); ok {
		return "/static/" + name + "?v=" + h
	}
	return "/static/" + name
}

// hash returns the hash of the named file, and whether the file exists.
func (s *Server) hash(name string) (string, bool) {
	if !s.live {
		h, ok := s.hashes[name]
		return h, ok
	}
	data, err := s.read(name)
	if err != nil {
		return "", false
	}
	return hash(data), true
}

// read returns the contents of a regular file. Directories are not served.
func (s *Server) read(name string) ([]byte, error) {
	if !fs.ValidPath(name) {
		return nil, fs.ErrInvalid
	}
	info, err := fs.Stat(s.fsys, name)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return nil, fs.ErrNotExist
	}
	return fs.ReadFile(s.fsys, name)
}

// ServeHTTP serves the file named by the request path, which must already
// have the /static/ prefix removed. The hash is the ETag, so a browser that
// has the file gets a 304.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/")
	data, err := s.read(name)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	h := hash(data)
	if !s.live {
		h = s.hashes[name]
	}

	w.Header().Set("ETag", `"`+h+`"`)
	switch {
	case s.live:
		w.Header().Set("Cache-Control", cacheLive)
	case r.URL.Query().Get("v") == h:
		w.Header().Set("Cache-Control", cacheHashed)
	default:
		w.Header().Set("Cache-Control", cacheShort)
	}
	http.ServeContent(w, r, name, time.Time{}, bytes.NewReader(data))
}

// hash returns the first 10 hex characters of the SHA-256 of data.
func hash(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])[:10]
}
