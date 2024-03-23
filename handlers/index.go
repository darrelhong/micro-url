package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"html/template"
	"net/http"

	"github.com/gorilla/sessions"
)

func HandleIndex(ghClientId string, sessionStore *sessions.CookieStore) http.Handler {
	tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/index.html", "templates/partials/head.html"))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b := make([]byte, 16)
		_, err := rand.Read(b)
		if err != nil {
			http.Error(w, "Something went wrong, please try again", http.StatusInternalServerError)
		}

		state := base64.StdEncoding.EncodeToString(b)

		session, err := sessionStore.Get(r, "state")

		if err != nil {
			http.Error(w, "Something went wrong, please try again", http.StatusInternalServerError)
		}

		session.Values["state"] = state
		session.Save(r, w)

		err = tmpl.ExecuteTemplate(w, "base", struct {
			GhClientId string
			State      string
		}{
			GhClientId: ghClientId,
			State:      state,
		})

		if err != nil {
			http.Error(w, "Something went wrong, please try again", http.StatusInternalServerError)
		}
	})
}
