package docs

import "embed"

// FS holds the documentation pages served under /docs/, and the
// last-updated dates in generated/dates.json. The deploy writes that file
// before the build; it is not committed, and generated/.gitkeep keeps the
// directory, and so this pattern, valid without it.
//
//go:embed all:guides all:reference all:architecture all:contributing all:generated index.md
var FS embed.FS
