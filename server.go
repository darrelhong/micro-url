package main

import (
	"net/http"

	"github.com/darrelhong/micro-url/store"
)

func NewServer(urlStore store.UrlStore) http.Handler {
	mux := http.NewServeMux()

	addRoutes(mux, urlStore)

	var handler http.Handler = mux

	return handler
}
