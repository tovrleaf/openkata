package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tovrleaf/openkata/internal/docsdates"
)

// gitRepo creates a repository whose commits have fixed dates.
type gitRepo struct {
	t   *testing.T
	dir string
}

func newGitRepo(t *testing.T) *gitRepo {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not installed")
	}
	dir := t.TempDir()
	t.Setenv("GIT_CEILING_DIRECTORIES", filepath.Dir(dir))
	r := &gitRepo{t: t, dir: dir}
	r.git("", "init", "-q", "-b", "main")
	return r
}

func (r *gitRepo) git(date string, args ...string) {
	r.t.Helper()
	cmd := exec.Command("git", append([]string{
		"-c", "user.name=Test", "-c", "user.email=test@example.com",
		"-c", "commit.gpgsign=false",
	}, args...)...)
	cmd.Dir = r.dir
	cmd.Env = append(os.Environ(), "GIT_CONFIG_GLOBAL="+os.DevNull, "GIT_CONFIG_NOSYSTEM=1")
	if date != "" {
		stamp := date + "T12:00:00+0000"
		cmd.Env = append(cmd.Env, "GIT_AUTHOR_DATE="+stamp, "GIT_COMMITTER_DATE="+stamp)
	}
	if out, err := cmd.CombinedOutput(); err != nil {
		r.t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

func (r *gitRepo) write(path, content string) {
	r.t.Helper()
	full := filepath.Join(r.dir, path)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		r.t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
		r.t.Fatal(err)
	}
}

func (r *gitRepo) commit(date string, files map[string]string) {
	r.t.Helper()
	for path, content := range files {
		r.write(path, content)
	}
	r.git(date, "add", "-A")
	r.git(date, "commit", "-q", "-m", "commit on "+date)
}

// docsRepo has a docs tree with pages committed on different days and a
// generated directory that must never be listed.
func docsRepo(t *testing.T) *gitRepo {
	r := newGitRepo(t)
	r.commit("2026-06-01", map[string]string{
		"web/docs/index.md":           "# Docs",
		"web/docs/guides/index.md":    "# Guides",
		"web/docs/guides/intro.md":    "# Intro",
		"web/docs/embed.go":           "package docs",
		"web/docs/generated/.gitkeep": "",
	})
	r.commit("2026-06-08", map[string]string{"web/docs/guides/intro.md": "# Intro, edited"})
	return r
}

func readDates(t *testing.T, path string) map[string]string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var dates map[string]string
	if err := json.Unmarshal(data, &dates); err != nil {
		t.Fatalf("%s is not valid JSON: %v\n%s", path, err, data)
	}
	return dates
}

func TestRunWritesDates(t *testing.T) {
	r := docsRepo(t)
	dir := filepath.Join(r.dir, "web/docs")
	out := filepath.Join(dir, "generated/dates.json")

	var log bytes.Buffer
	if err := run(dir, out, true, &log); err != nil {
		t.Fatalf("run() error = %v", err)
	}

	got := readDates(t, out)
	want := map[string]string{
		"index.md":        "2026-06-01",
		"guides/index.md": "2026-06-01",
		"guides/intro.md": "2026-06-08",
	}
	if len(got) != len(want) {
		t.Errorf("dates = %v, want exactly %v (Markdown pages only)", got, want)
	}
	for page, date := range want {
		if got[page] != date {
			t.Errorf("dates[%q] = %q, want %q", page, got[page], date)
		}
	}
	if !strings.Contains(log.String(), "wrote 3 dates") {
		t.Errorf("log = %q, want it to say 3 dates were written", log.String())
	}
}

func TestRunOutputIsStable(t *testing.T) {
	r := docsRepo(t)
	dir := filepath.Join(r.dir, "web/docs")
	out := filepath.Join(dir, "generated/dates.json")

	if err := run(dir, out, true, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	first, _ := os.ReadFile(out)
	if err := run(dir, out, true, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	second, _ := os.ReadFile(out)
	if !bytes.Equal(first, second) {
		t.Errorf("two runs wrote different files:\n%s\n---\n%s", first, second)
	}
	if !bytes.HasSuffix(first, []byte("}\n")) {
		t.Errorf("dates.json = %q, want a trailing newline", first)
	}
}

func TestRunCreatesTheOutputDirectory(t *testing.T) {
	r := docsRepo(t)
	dir := filepath.Join(r.dir, "web/docs")
	out := filepath.Join(r.dir, "somewhere", "new", "dates.json")
	if err := run(dir, out, true, &bytes.Buffer{}); err != nil {
		t.Fatalf("run() error = %v", err)
	}
	if _, err := os.Stat(out); err != nil {
		t.Errorf("os.Stat(%s) error = %v, want the file", out, err)
	}
}

func TestRunPagesWithoutDates(t *testing.T) {
	r := docsRepo(t)
	r.write("web/docs/guides/draft.md", "# Draft, not committed")
	dir := filepath.Join(r.dir, "web/docs")
	out := filepath.Join(dir, "generated/dates.json")

	t.Run("require fails and names the page", func(t *testing.T) {
		err := run(dir, out, true, &bytes.Buffer{})
		if err == nil {
			t.Fatal("run(require) error = nil, want an error for the uncommitted page")
		}
		if !strings.Contains(err.Error(), "guides/draft.md") {
			t.Errorf("run(require) error = %q, want it to name guides/draft.md", err)
		}
	})

	t.Run("without require it warns and writes the rest", func(t *testing.T) {
		var log bytes.Buffer
		if err := run(dir, out, false, &log); err != nil {
			t.Fatalf("run() error = %v", err)
		}
		if !strings.Contains(log.String(), "warning:") || !strings.Contains(log.String(), "guides/draft.md") {
			t.Errorf("log = %q, want a warning naming guides/draft.md", log.String())
		}
		got := readDates(t, out)
		if _, ok := got["guides/draft.md"]; ok {
			t.Errorf("dates has an entry for the uncommitted page: %v", got)
		}
		if got["guides/intro.md"] != "2026-06-08" {
			t.Errorf("dates[guides/intro.md] = %q, want 2026-06-08", got["guides/intro.md"])
		}
	})
}

func TestRunFailsWithoutUsableHistory(t *testing.T) {
	t.Run("shallow clone points to the checkout depth", func(t *testing.T) {
		src := docsRepo(t)
		shallow := filepath.Join(t.TempDir(), "shallow")
		clone := exec.Command("git", "clone", "-q", "--depth", "1", "file://"+src.dir, shallow)
		clone.Env = append(os.Environ(), "GIT_CONFIG_GLOBAL="+os.DevNull, "GIT_CONFIG_NOSYSTEM=1")
		if out, err := clone.CombinedOutput(); err != nil {
			t.Fatalf("git clone --depth 1: %v\n%s", err, out)
		}
		out := filepath.Join(t.TempDir(), "dates.json")

		err := run(filepath.Join(shallow, "web/docs"), out, true, &bytes.Buffer{})
		if !errors.Is(err, docsdates.ErrShallowClone) {
			t.Fatalf("run() on a shallow clone error = %v, want %v", err, docsdates.ErrShallowClone)
		}
		if !strings.Contains(err.Error(), "fetch-depth: 0") {
			t.Errorf("run() error = %q, want a hint about fetch-depth: 0", err)
		}
		if _, statErr := os.Stat(out); !errors.Is(statErr, os.ErrNotExist) {
			t.Errorf("os.Stat(%s) error = %v, want no file written", out, statErr)
		}
	})

	t.Run("not a repository", func(t *testing.T) {
		if _, err := exec.LookPath("git"); err != nil {
			t.Skip("git is not installed")
		}
		dir := t.TempDir()
		t.Setenv("GIT_CEILING_DIRECTORIES", filepath.Dir(dir))
		err := run(dir, filepath.Join(dir, "dates.json"), true, &bytes.Buffer{})
		if !errors.Is(err, docsdates.ErrNotRepository) {
			t.Errorf("run() outside a repository error = %v, want %v", err, docsdates.ErrNotRepository)
		}
	})
}
