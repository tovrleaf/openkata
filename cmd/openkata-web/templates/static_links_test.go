package templates

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var (
	// literalStatic matches a /static/ URL written as text: in double or
	// single quotes, in backticks, or straight after "=" in an attribute.
	literalStatic = regexp.MustCompile("[\"'`=]/static/[^\"'`\\s>]*")
	// assetURL matches a call to assets.URL and captures the file it names.
	assetURL = regexp.MustCompile(`assets\.URL\("([^"]+)"\)`)
)

// TestStaticLinksCarryHashes checks that templates and the Go helpers next to
// them link static files through assets.URL, so each URL carries the file's
// hash, and that every file they name exists. Generated _templ.go files
// repeat their .templ source and are skipped.
func TestStaticLinksCarryHashes(t *testing.T) {
	files, err := sourceFiles()
	if err != nil {
		t.Fatalf("sourceFiles() error = %v", err)
	}
	if len(files) == 0 {
		t.Fatal("sourceFiles() found no templates")
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

// sourceFiles lists the .templ files and the hand-written .go files of this
// package, leaving out generated code and tests.
func sourceFiles() ([]string, error) {
	templs, err := filepath.Glob("*.templ")
	if err != nil {
		return nil, err
	}
	gos, err := filepath.Glob("*.go")
	if err != nil {
		return nil, err
	}
	for _, f := range gos {
		if !strings.HasSuffix(f, "_templ.go") && !strings.HasSuffix(f, "_test.go") {
			templs = append(templs, f)
		}
	}
	return templs, nil
}

func TestLiteralStaticPattern(t *testing.T) {
	for _, in := range []string{
		`<link href="/static/css/style.css">`,
		`<img src='/static/img/a.png'>`,
		"x := `/static/js/a.js`",
		`<img src=/static/img/a.png>`,
	} {
		if !literalStatic.MatchString(in) {
			t.Errorf("literalStatic does not match %s, want a match", in)
		}
	}
	for _, in := range []string{
		`href={ assets.URL("css/style.css") }`,
		`os.ReadFile("web/static/versions.json")`,
	} {
		if literalStatic.MatchString(in) {
			t.Errorf("literalStatic matches %s, want no match", in)
		}
	}
}
