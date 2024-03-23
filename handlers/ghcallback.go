package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

func HandleGhCallback(oauth2Conf *oauth2.Config, sessionStore *sessions.CookieStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get code and state from query params
		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")

		session, err := sessionStore.Get(r, "state")

		if err != nil {
			http.Error(w, "Something went wrong, please try again", http.StatusInternalServerError)
		}

		if state != session.Values["state"] {
			http.Error(w, "Invalid state", http.StatusBadRequest)
			return
		}

		tok, err := oauth2Conf.Exchange(r.Context(), code)

		if err != nil {
			http.Error(w, "Something went wrong, please try again", http.StatusInternalServerError)
			return
		}

		client := oauth2Conf.Client(r.Context(), tok)

		resp, err := client.Get("https://api.github.com/user/emails")

		if err != nil {
			http.Error(w, "Something went wrong, please try again", http.StatusInternalServerError)
			return
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Something went wrong, please try again", http.StatusInternalServerError)
			return
		}

		var emails []struct {
			Email    string `json:"email"`
			Primary  bool   `json:"primary"`
			Verified bool   `json:"verified"`
		}

		err = json.Unmarshal(body, &emails)

		if err != nil {
			http.Error(w, "Something went wrong, please try again", http.StatusInternalServerError)
			return
		}

		if len(emails) == 0 {
			http.Error(w, "No email found", http.StatusBadRequest)
			return
		}

		if !emails[0].Verified {
			http.Error(w, "Email not verified with GitHub", http.StatusBadRequest)
			return
		}

		session, err = sessionStore.Get(r, "login")
		if err != nil {
			http.Error(w, "Something went wrong, please try again", http.StatusInternalServerError)
			return
		}

		session.Values["email"] = emails[0].Email
		session.Save(r, w)

		http.Redirect(w, r, "/", http.StatusFound)
	})
}
