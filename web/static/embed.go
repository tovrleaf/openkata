package static

import "embed"

//go:embed all:css all:js all:img versions.json favicon.ico favicon.svg favicon-96x96.png apple-touch-icon.png web-app-manifest-192x192.png web-app-manifest-512x512.png site.webmanifest
var FS embed.FS
