package docs

import "embed"

// FS holds the documentation pages served under /docs/.
//
//go:embed all:guides all:reference all:architecture all:contributing index.md
var FS embed.FS
