package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"html/template"
	"net/http"

	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

func HandleIndex(oauth2Conf *oauth2.Config, sessionStore *sessions.CookieStore) http.Handler {
	tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/index.html", "templates/partials/head.html"))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b := make([]byte, 16)
		_, err := rand.Read(b)
		if err != nil {
			http.Error(w, "Something went wrong, please try again", http.StatusInternalServerError)
			return
		}

		state := base64.StdEncoding.EncodeToString(b)

		loginSession, err := sessionStore.Get(r, "login")

		if err != nil {
			http.Error(w, "Something went wrong, please try again", http.StatusInternalServerError)
			return
		}

		var oauth2Url string

		if loginSession.Values["email"] == nil {
			oauth2Url = oauth2Conf.AuthCodeURL(state)

			session, err := sessionStore.Get(r, "state")

			if err != nil {
				http.Error(w, "Something went wrong, please try again", http.StatusInternalServerError)
				return
			}

			session.Values["state"] = state
			session.Save(r, w)
		}

		err = tmpl.ExecuteTemplate(w, "base", struct {
			Oauth2Url string
		}{
			Oauth2Url: oauth2Url,
		})

		if err != nil {
			http.Error(w, "Something went wrong, please try again", http.StatusInternalServerError)
			return
		}
	})
}
