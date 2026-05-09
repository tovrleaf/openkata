package main

import (
	"net/http"

	"github.com/tovrleaf/openkata/cmd/openkata-web/templates"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	templates.Home().Render(r.Context(), w)
}
