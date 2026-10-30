package main

import (
	"os/exec"
	"strings"
	"testing"
)

// repoWithTag makes a repository with one commit and the given tags.
func repoWithTag(t *testing.T, tags ...string) string {
	t.Helper()
	dir := t.TempDir()
	steps := [][]string{{"init", "-q"}, {"commit", "-q", "--allow-empty", "-m", "first"}}
	for _, tag := range tags {
		steps = append(steps, []string{"tag", tag})
	}
	for _, args := range steps {
		cmd := exec.Command("git", append([]string{
			"-c", "user.name=t", "-c", "user.email=t@example.com", "-c", "commit.gpgsign=false",
		}, args...)...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	return dir
}

func TestRun(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not installed")
	}

	t.Run("prints the release", func(t *testing.T) {
		var out strings.Builder
		if err := run(repoWithTag(t, "v0.2.0", "v0.10.0"), true, &out); err != nil {
			t.Fatalf("run() error = %v", err)
		}
		if got := out.String(); got != "v0.10.0\n" {
			t.Errorf("run() printed %q, want v0.10.0", got)
		}
	})

	t.Run("without a tag it prints dev", func(t *testing.T) {
		var out strings.Builder
		if err := run(repoWithTag(t, "skills/grill-me/v1.0.1"), false, &out); err != nil {
			t.Fatalf("run() error = %v", err)
		}
		if got := out.String(); got != "dev\n" {
			t.Errorf("run() printed %q, want dev", got)
		}
	})

	t.Run("without a tag, required, it fails and prints nothing", func(t *testing.T) {
		var out strings.Builder
		if err := run(repoWithTag(t), true, &out); err == nil {
			t.Error("run() error = nil, want an error when no release tag exists")
		}
		if out.Len() != 0 {
			t.Errorf("run() printed %q, want nothing on failure", out.String())
		}
	})
}
