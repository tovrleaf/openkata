package main

import (
	"regexp"
	"runtime/debug"
	"strings"
	"testing"

	"github.com/tovrleaf/openkata/internal/siteversion"
)

func TestVersionWithCommit(t *testing.T) {
	const full = "4f2a9c1d8e3b5a6f7c0d1e2f3a4b5c6d7e8f9a0b"

	tests := []struct {
		name     string
		settings []debug.BuildSetting
		want     string
	}{
		{
			"clean build uses the short commit",
			[]debug.BuildSetting{{Key: "vcs.revision", Value: full}, {Key: "vcs.modified", Value: "false"}},
			"0.3.0+4f2a9c1",
		},
		{
			"uncommitted changes are marked",
			[]debug.BuildSetting{{Key: "vcs.revision", Value: full}, {Key: "vcs.modified", Value: "true"}},
			"0.3.0+4f2a9c1.dirty",
		},
		{
			"modified flag is optional",
			[]debug.BuildSetting{{Key: "vcs.revision", Value: full}},
			"0.3.0+4f2a9c1",
		},
		{
			"a revision shorter than seven characters is kept whole",
			[]debug.BuildSetting{{Key: "vcs.revision", Value: "abc12"}},
			"0.3.0+abc12",
		},
		{
			"other settings are ignored",
			[]debug.BuildSetting{{Key: "GOARCH", Value: "arm64"}, {Key: "vcs.revision", Value: full}},
			"0.3.0+4f2a9c1",
		},
		{"no settings", nil, "0.3.0+dev"},
		{"settings without a revision", []debug.BuildSetting{{Key: "GOOS", Value: "linux"}}, "0.3.0+dev"},
		{
			"an empty revision counts as none",
			[]debug.BuildSetting{{Key: "vcs.revision", Value: ""}, {Key: "vcs.modified", Value: "true"}},
			"0.3.0+dev",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := versionWithCommit("0.3.0", tt.settings); got != tt.want {
				t.Errorf("versionWithCommit(%q, %v) = %q, want %q", "0.3.0", tt.settings, got, tt.want)
			}
		})
	}
}

// semverBuild matches a semantic version with optional build metadata, the
// form clients parse.
var semverBuild = regexp.MustCompile(`^\d+\.\d+\.\d+(\+[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?$`)

func TestServerVersion(t *testing.T) {
	got := serverVersion()
	if !strings.HasPrefix(got, baseVersion) {
		t.Errorf("serverVersion() = %q, want it to start with the base version %q", got, baseVersion)
	}
	if !semverBuild.MatchString(got) {
		t.Errorf("serverVersion() = %q, want a semantic version with build metadata", got)
	}
}

func TestReleaseBase(t *testing.T) {
	tests := []struct{ release, want string }{
		{"v0.14.0", "0.14.0"},
		{"v1.2.3", "1.2.3"},
		{"0.14.0", "0.14.0"},
		{"dev", baseVersion},
		{"", baseVersion},
		{"v0.14", baseVersion},
		{"v0.14.0-rc1", baseVersion},
		{"skills/grill-me/v1.0.1", baseVersion},
	}
	for _, tt := range tests {
		if got := releaseBase(tt.release); got != tt.want {
			t.Errorf("releaseBase(%q) = %q, want %q", tt.release, got, tt.want)
		}
	}
}

func TestServerVersionOfARelease(t *testing.T) {
	old := siteversion.Release
	t.Cleanup(func() { siteversion.Release = old })

	siteversion.Release = "v0.14.0"
	got := serverVersion()
	if !strings.HasPrefix(got, "0.14.0+") {
		t.Errorf("serverVersion() = %q, want it to start with 0.14.0+", got)
	}
	if !semverBuild.MatchString(got) {
		t.Errorf("serverVersion() = %q, want a semantic version with build metadata", got)
	}

	siteversion.Release = ""
	if got := serverVersion(); !strings.HasPrefix(got, baseVersion+"+") {
		t.Errorf("serverVersion() without a release = %q, want it to start with %s+", got, baseVersion)
	}
}

func TestVersionStringsAreValidSemver(t *testing.T) {
	for _, settings := range [][]debug.BuildSetting{
		nil,
		{{Key: "vcs.revision", Value: "4f2a9c1d8e3b5a6f"}},
		{{Key: "vcs.revision", Value: "4f2a9c1d8e3b5a6f"}, {Key: "vcs.modified", Value: "true"}},
	} {
		if got := versionWithCommit(baseVersion, settings); !semverBuild.MatchString(got) {
			t.Errorf("versionWithCommit(%v) = %q, want a semantic version with build metadata", settings, got)
		}
	}
}
