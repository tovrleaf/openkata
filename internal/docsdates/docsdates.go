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
// Two git processes read the history and the file list, however many
// files there are.
//
// The map is never nil. When git is missing or dir is not in a repository
// the map is empty and the error says why, so a caller that can live
// without dates may ignore it. A repository without commits is not an
// error.
func Lookup(dir string) (map[string]string, error) {
	dates := map[string]string{}

	git, err := exec.LookPath("git")
	if err != nil {
		return dates, fmt.Errorf("%w: %v", ErrGitUnavailable, err)
	}

	// Newest commit first, so the first time a path is seen is its last
	// change. Merge commits list no files, and renames list the new path
	// only. core.quotepath keeps non-ASCII names readable.
	cmd := exec.Command(git,
		"-c", "core.quotepath=off",
		"log", "-M", "--relative",
		"--format=%x01%as", "--name-only",
		"--", ".")
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "LC_ALL=C")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	out, err := cmd.Output()
	if err != nil {
		msg := stderr.String()
		switch {
		case strings.Contains(msg, "not a git repository"):
			return dates, fmt.Errorf("%w: %s", ErrNotRepository, dir)
		case strings.Contains(msg, "does not have any commits yet"):
			return dates, nil
		}
		return dates, fmt.Errorf("git log in %s: %w: %s", dir, err, strings.TrimSpace(msg))
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
	tracked, err := trackedFiles(git, dir)
	if err != nil {
		return map[string]string{}, err
	}
	for path := range dates {
		if !tracked[path] {
			delete(dates, path)
		}
	}
	return dates, nil
}

// trackedFiles lists the files under dir that git tracks, relative to dir.
func trackedFiles(git, dir string) (map[string]bool, error) {
	cmd := exec.Command(git, "-c", "core.quotepath=off", "ls-files", "-z", "--", ".")
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "LC_ALL=C")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git ls-files in %s: %w", dir, err)
	}
	tracked := map[string]bool{}
	for _, path := range strings.Split(string(out), "\x00") {
		if path != "" {
			tracked[path] = true
		}
	}
	return tracked, nil
}
