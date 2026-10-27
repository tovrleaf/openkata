package main

import "runtime/debug"

// baseVersion is the version of the server's tools and behavior. Change it
// by hand when a change deserves a new number. The build adds the commit.
const baseVersion = "0.3.0"

// serverVersion returns the version that clients see in the handshake:
// baseVersion followed by the commit the binary was built from, as semantic
// version build metadata, for example "0.3.0+4f2a9c1". A build from a
// working tree with uncommitted changes adds ".dirty", and a build without
// version control information, such as go run, gets "+dev".
func serverVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return versionWithCommit(baseVersion, nil)
	}
	return versionWithCommit(baseVersion, info.Settings)
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
