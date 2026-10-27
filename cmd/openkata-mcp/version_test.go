package main

import (
	"regexp"
	"runtime/debug"
	"strings"
	"testing"
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
