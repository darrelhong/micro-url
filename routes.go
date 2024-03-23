package main

import (
	"embed"

	"github.com/darrelhong/micro-url/handlers"
	"github.com/darrelhong/micro-url/store"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"

	"net/http"
)

//go:embed static
var static embed.FS

func addRoutes(mux *http.ServeMux, urlStore store.UrlStore, oauth2Conf *oauth2.Config, sessionStore *sessions.CookieStore) {
	mux.Handle("/", handlers.HandleIndex(oauth2Conf, sessionStore))

	mux.Handle("GET /github/callback", handlers.HandleGhCallback(oauth2Conf, sessionStore))

	mux.Handle("POST /shorten", handlers.HandleShorten(urlStore))

	mux.Handle("GET /{shortUrlId}", handlers.HandleRedirect(urlStore))

	mux.Handle("/static/", http.FileServer(http.FS(static)))
}
