package main

import (
	"net/http"

	"github.com/darrelhong/micro-url/store"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

func NewServer(urlStore store.UrlStore, oauth2Conf *oauth2.Config, sessionStore *sessions.CookieStore) http.Handler {
	mux := http.NewServeMux()

	addRoutes(mux, urlStore, oauth2Conf, sessionStore)

	var handler http.Handler = mux

	return handler
}
