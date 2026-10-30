package templates

import "github.com/tovrleaf/openkata/internal/siteversion"

const releasesURL = "https://github.com/tovrleaf/openkata/releases/tag/"

// releaseURL returns the GitHub page of a release, or "" for a build that
// has none, so the footer shows the version as plain text.
func releaseURL(version string) string {
	if version == "" || version == siteversion.Dev {
		return ""
	}
	return releasesURL + version
}
