// Package siteversion gives the project its release number: the newest
// tag of the form vMAJOR.MINOR.PATCH that the checked-out commit contains.
// Skills, rules and profiles have tags of their own, such as
// skills/<name>/v1.0.0, and are never mistaken for it.
package siteversion

import (
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// Dev is the version of a build that does not come from a release tag.
const Dev = "dev"

// Release is set when the binary is built, with
// -ldflags "-X github.com/tovrleaf/openkata/internal/siteversion.Release=v0.1.0".
var Release string

// Current returns the release the running binary was built from, such as
// "v0.14.0", or Dev when the build did not set one.
func Current() string {
	if Release == "" {
		return Dev
	}
	return Release
}

var tagPattern = regexp.MustCompile(`^v(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$`)

// Newest returns the highest release tag in tags, compared by number so
// that v0.10.0 is newer than v0.9.0, or Dev when there is none. Tags that
// are not exactly vMAJOR.MINOR.PATCH are ignored: artifact tags, a
// pre-release suffix, or a version without the v.
func Newest(tags []string) string {
	best := ""
	var bestNums [3]int
	for _, tag := range tags {
		m := tagPattern.FindStringSubmatch(strings.TrimSpace(tag))
		if m == nil {
			continue
		}
		var nums [3]int
		for i := range nums {
			nums[i], _ = strconv.Atoi(m[i+1])
		}
		if best == "" || newer(nums, bestNums) {
			best, bestNums = m[0], nums
		}
	}
	if best == "" {
		return Dev
	}
	return best
}

func newer(a, b [3]int) bool {
	for i := range a {
		if a[i] != b[i] {
			return a[i] > b[i]
		}
	}
	return false
}

// FromGit returns the newest release tag that the commit checked out in
// dir contains, or Dev when there is none, when dir is not a repository or
// when git is not installed. A tag on another branch is not counted, and a
// tag that was never fetched cannot be seen, so build from a clone with
// its tags.
func FromGit(dir string) string {
	cmd := exec.Command("git", "tag", "--merged", "HEAD", "--list", "v[0-9]*")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return Dev
	}
	return Newest(strings.Fields(string(out)))
}
