package static

import "embed"

//go:embed all:css all:js all:img versions.json
var FS embed.FS
