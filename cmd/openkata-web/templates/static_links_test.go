package templates

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

var (
	// literalStatic matches a /static/ URL written into a template as text.
	literalStatic = regexp.MustCompile(`"/static/[^"]*"`)
	// assetURL matches a call to assets.URL and captures the file it names.
	assetURL = regexp.MustCompile(`assets\.URL\("([^"]+)"\)`)
)

// TestStaticLinksCarryHashes checks that templates link static files through
// assets.URL, so each URL carries the file's hash, and that every file they
// name exists.
func TestStaticLinksCarryHashes(t *testing.T) {
	files, err := filepath.Glob("*.templ")
	if err != nil {
		t.Fatalf("Glob(*.templ) error = %v", err)
	}
	if len(files) == 0 {
		t.Fatal("Glob(*.templ) found no templates")
	}
	static := filepath.Join("..", "..", "..", "web", "static")

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("ReadFile(%s) error = %v", file, err)
		}
		for _, lit := range literalStatic.FindAll(data, -1) {
			t.Errorf("%s links %s directly, want assets.URL so the URL carries a hash", file, lit)
		}
		for _, m := range assetURL.FindAllSubmatch(data, -1) {
			name := string(m[1])
			if _, err := os.Stat(filepath.Join(static, name)); err != nil {
				t.Errorf("%s links static file %q, which does not exist in web/static", file, name)
			}
		}
	}
}
