package main

import (
	"net/http"

	"github.com/darrelhong/micro-url/store"
	"github.com/darrelhong/micro-url/utils"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

func NewServer(
	urlStore store.UrlStore,
	oauth2Conf *oauth2.Config,
	sessionStore *sessions.CookieStore,
	userStore store.UserStore,
	tursoApiClient *utils.TursoApiClient,
	userDbClient *utils.UserDbClient,
) http.Handler {
	mux := http.NewServeMux()

	addRoutes(mux, urlStore, oauth2Conf, sessionStore, userStore, tursoApiClient, userDbClient)

	var handler http.Handler = mux

	return handler
}
