package main

import (
	"embed"

	"github.com/darrelhong/micro-url/handlers"

	"net/http"
)

//go:embed static
var static embed.FS

func addRoutes(mux *http.ServeMux) {
	mux.Handle("/", handlers.HandleIndex())

	mux.Handle("POST /shorten", handlers.HandleShorten())

	mux.Handle("/static/", http.FileServer(http.FS(static)))
}
