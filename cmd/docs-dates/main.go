// Command docs-dates writes the last-updated date of every docs page,
// read from git, to a file that is embedded in the web binary at deploy
// time. The file is generated, gitignored and never committed.
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/tovrleaf/openkata/internal/docsdates"
)

const generatedDir = "generated"

func main() {
	dir := flag.String("dir", "web/docs", "docs directory")
	out := flag.String("out", "web/docs/generated/dates.json", "file to write")
	require := flag.Bool("require", false, "fail if any page has no date")
	flag.Parse()

	if err := run(*dir, *out, *require, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "docs-dates:", err)
		os.Exit(1)
	}
}

// run writes the dates of the Markdown pages under dir to out. Pages git
// has no date for are reported, and are an error with require, so a deploy
// can refuse to ship pages that would show the placeholder.
func run(dir, out string, require bool, log io.Writer) error {
	all, err := docsdates.Lookup(dir)
	if err != nil {
		if errors.Is(err, docsdates.ErrShallowClone) {
			return fmt.Errorf("%w (in GitHub Actions, set fetch-depth: 0 on actions/checkout)", err)
		}
		return err
	}

	pages, err := markdownFiles(dir)
	if err != nil {
		return err
	}
	dates := map[string]string{}
	var missing []string
	for _, page := range pages {
		if date, ok := all[page]; ok {
			dates[page] = date
		} else {
			missing = append(missing, page)
		}
	}

	if len(missing) > 0 {
		msg := fmt.Sprintf("%d page(s) have no date because git does not track them: %s",
			len(missing), strings.Join(missing, ", "))
		if require {
			return errors.New(msg)
		}
		fmt.Fprintln(log, "warning:", msg)
	}

	data, err := json.MarshalIndent(dates, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(out, append(data, '\n'), 0o644); err != nil {
		return err
	}
	fmt.Fprintf(log, "wrote %d dates to %s\n", len(dates), out)
	return nil
}

// markdownFiles lists the .md files under dir, relative to it with forward
// slashes and sorted, leaving out the generated directory.
func markdownFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if d.IsDir() {
			if rel == generatedDir {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(rel, ".md") {
			files = append(files, rel)
		}
		return nil
	})
	sort.Strings(files)
	return files, err
}
