// Package docsdates finds when documentation files were last updated,
// from git history.
package docsdates

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	// ErrGitUnavailable means the git binary was not found.
	ErrGitUnavailable = errors.New("git is not available")
	// ErrNotRepository means the directory is not inside a git repository.
	ErrNotRepository = errors.New("not a git repository")
	// ErrShallowClone means the repository has truncated history, so the
	// last commit that touched a file cannot be known. In a clone with a
	// single commit every file would get that commit's date.
	ErrShallowClone = errors.New("shallow clone, full history is needed")
)

const dateLayout = "2006-01-02"

// Lookup returns the last-updated date of every file under dir that git
// has history for. The date is the author date of the last commit that
// changed the file, as YYYY-MM-DD, so an unrelated commit does not move
// it. Keys are slash-separated paths relative to dir.
//
// Only files git currently tracks are in the map. A new file that is not
// committed yet, or one that reuses the path of a deleted file, has no
// date; a file with uncommitted edits keeps the date of its last commit.
// Three git processes check the repository, read the history and list the
// tracked files, however many files there are.
//
// The map is never nil. When git is missing, dir is not in a repository
// or the repository is a shallow clone, the map is empty and the error
// says why, so a caller that can live without dates may ignore it. A
// repository without commits is not an error.
func Lookup(dir string) (map[string]string, error) {
	dates := map[string]string{}

	git, err := exec.LookPath("git")
	if err != nil {
		return dates, fmt.Errorf("%w: %v", ErrGitUnavailable, err)
	}

	shallow, err := runGit(git, dir, "rev-parse", "--is-shallow-repository")
	if err != nil {
		return dates, err
	}
	if strings.TrimSpace(string(shallow)) == "true" {
		return dates, fmt.Errorf("%w: %s", ErrShallowClone, dir)
	}

	// Newest commit first, so the first time a path is seen is its last
	// change. Merge commits list no files, and renames list the new path
	// only. core.quotepath keeps non-ASCII names readable.
	out, err := runGit(git, dir,
		"-c", "core.quotepath=off",
		"log", "-M", "--relative",
		"--format=%x01%as", "--name-only",
		"--", ".")
	if err != nil {
		if errors.Is(err, errNoCommits) {
			return dates, nil
		}
		return dates, err
	}

	var date string
	for _, line := range strings.Split(string(out), "\n") {
		switch {
		case line == "":
		case strings.HasPrefix(line, "\x01"):
			date = strings.TrimPrefix(line, "\x01")
			if _, err := time.Parse(dateLayout, date); err != nil {
				return map[string]string{}, fmt.Errorf("git log in %s: unexpected date %q", dir, date)
			}
		case date != "":
			if _, seen := dates[line]; !seen {
				dates[line] = date
			}
		}
	}

	// History also names files that were renamed or deleted since. Keep
	// only the ones git tracks now, so a new file at an old path has no
	// date of its own.
	listed, err := runGit(git, dir, "-c", "core.quotepath=off", "ls-files", "-z", "--", ".")
	if err != nil {
		return map[string]string{}, err
	}
	tracked := map[string]bool{}
	for _, path := range strings.Split(string(listed), "\x00") {
		if path != "" {
			tracked[path] = true
		}
	}
	for path := range dates {
		if !tracked[path] {
			delete(dates, path)
		}
	}
	return dates, nil
}

// errNoCommits is returned by runGit for a repository without commits.
var errNoCommits = errors.New("repository has no commits")

// runGit runs git in dir and maps the failures callers care about to
// errors. The C locale keeps the messages it matches on stable.
func runGit(git, dir string, args ...string) ([]byte, error) {
	cmd := exec.Command(git, args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "LC_ALL=C")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	out, err := cmd.Output()
	if err == nil {
		return out, nil
	}
	msg := stderr.String()
	switch {
	case strings.Contains(msg, "not a git repository"):
		return nil, fmt.Errorf("%w: %s", ErrNotRepository, dir)
	case strings.Contains(msg, "does not have any commits yet"):
		return nil, errNoCommits
	}
	return nil, fmt.Errorf("git %s in %s: %w: %s", args[len(args)-1], dir, err, strings.TrimSpace(msg))
}
