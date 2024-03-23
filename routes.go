package main

import (
	"embed"

	"github.com/darrelhong/micro-url/handlers"
	"github.com/darrelhong/micro-url/store"

	"net/http"
)

//go:embed static
var static embed.FS

func addRoutes(mux *http.ServeMux, urlStore store.UrlStore, ghClientId string) {
	mux.Handle("/", handlers.HandleIndex(ghClientId))

	mux.Handle("POST /shorten", handlers.HandleShorten(urlStore))

	mux.Handle("GET /{shortUrlId}", handlers.HandleRedirect(urlStore))

	mux.Handle("/static/", http.FileServer(http.FS(static)))
}
