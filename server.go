package main

import (
	"net/http"

	"github.com/darrelhong/micro-url/store"
)

func NewServer(urlStore store.UrlStore, ghClientId string) http.Handler {
	mux := http.NewServeMux()

	addRoutes(mux, urlStore, ghClientId)

	var handler http.Handler = mux

	return handler
}
