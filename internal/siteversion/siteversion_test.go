package siteversion

import (
	"os"
	"os/exec"
	"testing"
)

func TestNewest(t *testing.T) {
	tests := []struct {
		name string
		tags []string
		want string
	}{
		{"no tags", nil, Dev},
		{"one tag", []string{"v0.1.0"}, "v0.1.0"},
		{"compared by number, not text", []string{"v0.9.0", "v0.10.0", "v0.2.0"}, "v0.10.0"},
		{"major beats minor", []string{"v1.0.0", "v0.99.99"}, "v1.0.0"},
		{"patch is compared", []string{"v0.1.2", "v0.1.10", "v0.1.9"}, "v0.1.10"},
		{"order does not matter", []string{"v0.3.0", "v0.1.0", "v0.2.0"}, "v0.3.0"},
		{
			"artifact tags are ignored",
			[]string{"skills/grill-me/v9.9.9", "rules/bash-style/v2.0.0", "v0.4.0"},
			"v0.4.0",
		},
		{"only artifact tags", []string{"skills/grill-me/v1.0.1", "profiles/x/v1.0.0"}, Dev},
		{"a pre-release suffix is ignored", []string{"v0.5.0-rc1", "v0.4.0"}, "v0.4.0"},
		{"build metadata is ignored", []string{"v0.5.0+build", "v0.4.0"}, "v0.4.0"},
		{"a version without the v is ignored", []string{"0.5.0", "v0.4.0"}, "v0.4.0"},
		{"too few or too many parts are ignored", []string{"v1", "v1.2", "v1.2.3.4", "v0.4.0"}, "v0.4.0"},
		{"leading zeros are ignored", []string{"v01.2.3", "v0.4.0"}, "v0.4.0"},
		{"empty and padded entries", []string{"", "  v0.2.0\n"}, "v0.2.0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Newest(tt.tags); got != tt.want {
				t.Errorf("Newest(%q) = %q, want %q", tt.tags, got, tt.want)
			}
		})
	}
}

func TestCurrent(t *testing.T) {
	old := Release
	t.Cleanup(func() { Release = old })

	Release = ""
	if got := Current(); got != Dev {
		t.Errorf("Current() with no release set = %q, want %q", got, Dev)
	}
	Release = "v0.14.0"
	if got := Current(); got != "v0.14.0" {
		t.Errorf("Current() = %q, want v0.14.0", got)
	}
}

// gitIn runs git in dir and fails the test if it does not succeed.
func gitIn(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{
		"-c", "user.name=t", "-c", "user.email=t@example.com", "-c", "commit.gpgsign=false",
	}, args...)...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

func TestFromGit(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not installed")
	}

	t.Run("a repository without release tags", func(t *testing.T) {
		dir := t.TempDir()
		gitIn(t, dir, "init", "-q")
		gitIn(t, dir, "commit", "-q", "--allow-empty", "-m", "first")
		gitIn(t, dir, "tag", "skills/grill-me/v1.0.1")
		if got := FromGit(dir); got != Dev {
			t.Errorf("FromGit() = %q, want %q", got, Dev)
		}
	})

	t.Run("the newest release the commit contains", func(t *testing.T) {
		dir := t.TempDir()
		gitIn(t, dir, "init", "-q")
		gitIn(t, dir, "commit", "-q", "--allow-empty", "-m", "first")
		gitIn(t, dir, "tag", "v0.9.0")
		gitIn(t, dir, "commit", "-q", "--allow-empty", "-m", "second")
		gitIn(t, dir, "tag", "v0.10.0")
		gitIn(t, dir, "tag", "rules/bash-style/v3.0.0")
		gitIn(t, dir, "commit", "-q", "--allow-empty", "-m", "third, not released")
		if got := FromGit(dir); got != "v0.10.0" {
			t.Errorf("FromGit() = %q, want v0.10.0", got)
		}
	})

	t.Run("a tag on a commit that is not contained is not counted", func(t *testing.T) {
		dir := t.TempDir()
		gitIn(t, dir, "init", "-q")
		gitIn(t, dir, "commit", "-q", "--allow-empty", "-m", "first")
		gitIn(t, dir, "tag", "v0.1.0")
		gitIn(t, dir, "checkout", "-q", "-b", "other")
		gitIn(t, dir, "commit", "-q", "--allow-empty", "-m", "on another branch")
		gitIn(t, dir, "tag", "v0.2.0")
		gitIn(t, dir, "checkout", "-q", "-")
		if got := FromGit(dir); got != "v0.1.0" {
			t.Errorf("FromGit() = %q, want v0.1.0, the release this commit contains", got)
		}
	})

	t.Run("not a repository", func(t *testing.T) {
		if got := FromGit(t.TempDir()); got != Dev {
			t.Errorf("FromGit() = %q, want %q", got, Dev)
		}
	})

	t.Run("a directory that does not exist", func(t *testing.T) {
		if got := FromGit("/no/such/dir/for/siteversion"); got != Dev {
			t.Errorf("FromGit() = %q, want %q", got, Dev)
		}
	})
}

func TestMain(m *testing.M) {
	// Keep the tests independent of the user's git configuration.
	os.Setenv("GIT_CONFIG_GLOBAL", os.DevNull)
	os.Setenv("GIT_CONFIG_SYSTEM", os.DevNull)
	os.Exit(m.Run())
}
