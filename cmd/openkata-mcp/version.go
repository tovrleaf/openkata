package main

import (
	"regexp"
	"runtime/debug"
	"strings"

	"github.com/tovrleaf/openkata/internal/siteversion"
)

// baseVersion is the version of a build that does not come from a release,
// such as go run. A deployed server reports the release it was built from.
const baseVersion = "0.0.0"

var releaseNumber = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+$`)

// releaseBase returns the version of a release tag such as "v0.14.0" without
// the v, or baseVersion for anything that is not a release, such as "dev".
func releaseBase(release string) string {
	number := strings.TrimPrefix(release, "v")
	if !releaseNumber.MatchString(number) {
		return baseVersion
	}
	return number
}

// serverVersion returns the version that clients see in the handshake: the
// release the binary was built from, followed by the commit, as semantic
// version build metadata, for example "0.14.0+4f2a9c1". A build from a
// working tree with uncommitted changes adds ".dirty", a build without
// version control information, such as go run, gets "+dev", and a build
// that is not from a release starts with baseVersion.
func serverVersion() string {
	base := releaseBase(siteversion.Current())
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return versionWithCommit(base, nil)
	}
	return versionWithCommit(base, info.Settings)
}

// versionWithCommit builds the version string from the build settings that
// the Go toolchain records, "vcs.revision" and "vcs.modified".
func versionWithCommit(base string, settings []debug.BuildSetting) string {
	var revision string
	var modified bool
	for _, s := range settings {
		switch s.Key {
		case "vcs.revision":
			revision = s.Value
		case "vcs.modified":
			modified = s.Value == "true"
		}
	}
	if revision == "" {
		return base + "+dev"
	}
	if len(revision) > 7 {
		revision = revision[:7]
	}
	if modified {
		return base + "+" + revision + ".dirty"
	}
	return base + "+" + revision
}
