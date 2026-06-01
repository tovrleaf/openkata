package static

import "embed"

//go:embed all:css all:js all:img
var FS embed.FS
