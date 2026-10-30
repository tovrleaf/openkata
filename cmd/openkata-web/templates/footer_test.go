package templates

import (
	"context"
	"strings"
	"testing"

	"github.com/tovrleaf/openkata/internal/siteversion"
)

func renderFooter(t *testing.T, release string) string {
	t.Helper()
	old := siteversion.Release
	siteversion.Release = release
	t.Cleanup(func() { siteversion.Release = old })

	var b strings.Builder
	if err := Footer().Render(context.Background(), &b); err != nil {
		t.Fatalf("Footer().Render() error = %v", err)
	}
	return b.String()
}

func TestFooterLinksTheRelease(t *testing.T) {
	got := renderFooter(t, "v0.14.0")
	for _, want := range []string{
		`href="https://github.com/tovrleaf/openkata/releases/tag/v0.14.0"`,
		"Release v0.14.0</a>",
		"© 2026 OPENKATA.dev",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("Footer() = %s\nwant it to contain %q", got, want)
		}
	}
}

func TestFooterLinks(t *testing.T) {
	got := renderFooter(t, "v0.14.0")
	for _, want := range []string{
		"Community", `href="https://github.com/tovrleaf/openkata"`,
		"Explore", `<a href="/docs/">Documentation</a>`,
		"MIT License", `blob/main/LICENSE`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("Footer() = %s\nwant it to contain %q", got, want)
		}
	}
}

func TestFooterWithoutARelease(t *testing.T) {
	got := renderFooter(t, "")
	if !strings.Contains(got, "Release dev") {
		t.Errorf("Footer() = %s\nwant it to say Release dev", got)
	}
	if strings.Contains(got, "releases/tag") {
		t.Errorf("Footer() = %s\nwant no link to a release for a dev build", got)
	}
}

func TestReleaseURL(t *testing.T) {
	tests := []struct{ in, want string }{
		{"v0.14.0", "https://github.com/tovrleaf/openkata/releases/tag/v0.14.0"},
		{siteversion.Dev, ""},
		{"", ""},
	}
	for _, tt := range tests {
		if got := releaseURL(tt.in); got != tt.want {
			t.Errorf("releaseURL(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
