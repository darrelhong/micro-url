package main

import (
	"net/http"

	"github.com/darrelhong/micro-url/store"
	"github.com/gorilla/sessions"
)

func NewServer(urlStore store.UrlStore, ghClientId string, sessionStore *sessions.CookieStore) http.Handler {
	mux := http.NewServeMux()

	addRoutes(mux, urlStore, ghClientId, sessionStore)

	var handler http.Handler = mux

	return handler
}
