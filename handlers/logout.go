package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

func HandleLogout(sessionStore *sessions.CookieStore) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loginSession, err := sessionStore.Get(r, "login")

		if err != nil {
			log.Println("hmm couldn't get session")
			http.Error(w, "Something went wrong, please try again", http.StatusInternalServerError)
			return
		}

		loginSession.Options.MaxAge = -1

		err = loginSession.Save(r, w)
		if err != nil {
			log.Println("hmm couldn't save session")
			http.Error(w, "Something went wrong, please try again", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
}
