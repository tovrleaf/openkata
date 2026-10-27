package docsdates

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"testing"
)

// testRepo is a temporary git repository with commits at fixed dates.
type testRepo struct {
	t   *testing.T
	dir string
}

func newRepo(t *testing.T) *testRepo {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not installed")
	}
	dir := t.TempDir()
	// Stop git from finding a repository above the temporary directory.
	t.Setenv("GIT_CEILING_DIRECTORIES", filepath.Dir(dir))
	r := &testRepo{t: t, dir: dir}
	r.git("", "init", "-q", "-b", "main")
	return r
}

// git runs a git command in the repository. A non-empty date sets the
// author and committer date of any commit it makes, at noon UTC.
func (r *testRepo) git(date string, args ...string) {
	r.t.Helper()
	cmd := exec.Command("git", append([]string{
		"-c", "user.name=Test", "-c", "user.email=test@example.com",
		"-c", "commit.gpgsign=false",
	}, args...)...)
	cmd.Dir = r.dir
	cmd.Env = append(os.Environ(),
		"GIT_CONFIG_GLOBAL="+os.DevNull, "GIT_CONFIG_NOSYSTEM=1")
	if date != "" {
		stamp := date + "T12:00:00+0000"
		cmd.Env = append(cmd.Env, "GIT_AUTHOR_DATE="+stamp, "GIT_COMMITTER_DATE="+stamp)
	}
	if out, err := cmd.CombinedOutput(); err != nil {
		r.t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

func (r *testRepo) write(path, content string) {
	r.t.Helper()
	full := filepath.Join(r.dir, path)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		r.t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
		r.t.Fatal(err)
	}
}

// commit writes the files and commits everything at the given date.
func (r *testRepo) commit(date string, files map[string]string) {
	r.t.Helper()
	for path, content := range files {
		r.write(path, content)
	}
	r.git(date, "add", "-A")
	r.git(date, "commit", "-q", "-m", "commit on "+date)
}

func TestLookup(t *testing.T) {
	r := newRepo(t)
	r.commit("2026-03-01", map[string]string{
		"docs/a.md":        "a1",
		"docs/guides/b.md": "b1",
		"outside/c.md":     "c1",
	})
	r.commit("2026-03-02", map[string]string{"docs/x.md": "x1"})
	r.git("2026-03-03", "rm", "-q", "docs/x.md")
	r.git("2026-03-03", "commit", "-q", "-m", "delete x")
	r.commit("2026-03-05", map[string]string{"docs/a.md": "a2"})
	r.commit("2026-03-09", map[string]string{
		"docs/guides/d.md": "d1",
		"outside/c.md":     "c2",
	})
	r.git("2026-03-12", "mv", "docs/guides/b.md", "docs/guides/b2.md")
	r.git("2026-03-12", "commit", "-q", "-m", "rename")
	r.commit("2026-03-15", map[string]string{"docs/guides/käyttö.md": "k1"})
	r.commit("2026-03-20", map[string]string{"outside/c.md": "c3"})

	// Uncommitted work must not change or add dates.
	r.write("docs/a.md", "edited but not committed")
	r.write("docs/new.md", "untracked")
	r.write("docs/staged.md", "staged only")
	r.git("", "add", "docs/staged.md")
	r.write("docs/x.md", "recreated, not committed") // path of a deleted file

	got, err := Lookup(filepath.Join(r.dir, "docs"))
	if err != nil {
		t.Fatalf("Lookup(docs) error = %v", err)
	}

	tests := []struct {
		name string
		path string
		want string // empty means the file must not be in the map
	}{
		{"last commit that touched the file, not the last commit", "a.md", "2026-03-05"},
		{"file in a subdirectory", "guides/d.md", "2026-03-09"},
		{"renamed file has the date of the rename", "guides/b2.md", "2026-03-12"},
		{"old name of a renamed file is not listed", "guides/b.md", ""},
		{"non-ASCII file name", "guides/käyttö.md", "2026-03-15"},
		{"untracked file has no date", "new.md", ""},
		{"untracked file at the path of a deleted file has no date", "x.md", ""},
		{"staged but uncommitted file has no date", "staged.md", ""},
		{"file outside the directory is not listed", "../outside/c.md", ""},
		{"file outside, without the prefix", "outside/c.md", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			date, ok := got[tt.path]
			if tt.want == "" {
				if ok {
					t.Errorf("Lookup(docs)[%q] = %q, want no entry", tt.path, date)
				}
				return
			}
			if date != tt.want {
				t.Errorf("Lookup(docs)[%q] = %q, want %q", tt.path, date, tt.want)
			}
		})
	}

	t.Run("every date is YYYY-MM-DD", func(t *testing.T) {
		re := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
		for path, date := range got {
			if !re.MatchString(date) {
				t.Errorf("Lookup(docs)[%q] = %q, want YYYY-MM-DD", path, date)
			}
		}
	})

	t.Run("keys are relative to the directory", func(t *testing.T) {
		whole, err := Lookup(r.dir)
		if err != nil {
			t.Fatalf("Lookup(repo root) error = %v", err)
		}
		if whole["docs/a.md"] != "2026-03-05" || whole["outside/c.md"] != "2026-03-20" {
			t.Errorf("Lookup(repo root) = %v, want docs/a.md and outside/c.md keyed from the root", whole)
		}
	})
}

func TestLookupMergeKeepsTheCommitDate(t *testing.T) {
	r := newRepo(t)
	r.commit("2026-04-01", map[string]string{"docs/base.md": "base"})
	r.git("", "switch", "-q", "-c", "feature")
	r.commit("2026-04-10", map[string]string{"docs/m.md": "from the branch"})
	r.git("", "switch", "-q", "main")
	r.commit("2026-04-11", map[string]string{"docs/other.md": "unrelated"})
	r.git("2026-04-20", "merge", "-q", "--no-ff", "-m", "merge feature", "feature")

	got, err := Lookup(filepath.Join(r.dir, "docs"))
	if err != nil {
		t.Fatalf("Lookup(docs) error = %v", err)
	}
	if got["m.md"] != "2026-04-10" {
		t.Errorf("Lookup(docs)[m.md] = %q, want the branch commit date 2026-04-10, not the merge date", got["m.md"])
	}
	if got["other.md"] != "2026-04-11" || got["base.md"] != "2026-04-01" {
		t.Errorf("Lookup(docs) = %v, want other.md on 2026-04-11 and base.md on 2026-04-01", got)
	}
}

func TestLookupWithoutDates(t *testing.T) {
	t.Run("repository without commits", func(t *testing.T) {
		r := newRepo(t)
		got, err := Lookup(r.dir)
		if err != nil {
			t.Errorf("Lookup(empty repo) error = %v, want nil", err)
		}
		if got == nil || len(got) != 0 {
			t.Errorf("Lookup(empty repo) = %v, want an empty map", got)
		}
	})

	t.Run("directory that is not a repository", func(t *testing.T) {
		if _, err := exec.LookPath("git"); err != nil {
			t.Skip("git is not installed")
		}
		dir := t.TempDir()
		t.Setenv("GIT_CEILING_DIRECTORIES", filepath.Dir(dir))
		got, err := Lookup(dir)
		if !errors.Is(err, ErrNotRepository) {
			t.Errorf("Lookup(non-repo) error = %v, want %v", err, ErrNotRepository)
		}
		if got == nil || len(got) != 0 {
			t.Errorf("Lookup(non-repo) = %v, want an empty map", got)
		}
	})

	t.Run("git is not installed", func(t *testing.T) {
		t.Setenv("PATH", t.TempDir())
		got, err := Lookup(t.TempDir())
		if !errors.Is(err, ErrGitUnavailable) {
			t.Errorf("Lookup without git error = %v, want %v", err, ErrGitUnavailable)
		}
		if got == nil || len(got) != 0 {
			t.Errorf("Lookup without git = %v, want an empty map", got)
		}
	})

	t.Run("directory that does not exist", func(t *testing.T) {
		if _, err := exec.LookPath("git"); err != nil {
			t.Skip("git is not installed")
		}
		got, err := Lookup(filepath.Join(t.TempDir(), "missing"))
		if err == nil {
			t.Error("Lookup(missing dir) error = nil, want an error")
		}
		if got == nil || len(got) != 0 {
			t.Errorf("Lookup(missing dir) = %v, want an empty map", got)
		}
	})
}
