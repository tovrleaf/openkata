package templates

import (
	"context"
	"strings"
	"testing"
)

func TestNavLinksToTheDocs(t *testing.T) {
	var b strings.Builder
	if err := Nav().Render(context.Background(), &b); err != nil {
		t.Fatalf("Nav().Render() error = %v", err)
	}
	got := b.String()
	for _, want := range []string{
		`<a href="/docs/" class="nav-link">Docs</a>`,
		`<a href="/docs/" class="nav-menu-link">Docs</a>`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("Nav() = %s\nwant it to contain %q", got, want)
		}
	}
}
