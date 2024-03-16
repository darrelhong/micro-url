package main

import "net/http"

func NewServer() http.Handler {
	mux := http.NewServeMux()

	addRoutes(mux)

	var handler http.Handler = mux

	return handler
}
