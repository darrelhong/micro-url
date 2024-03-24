package handlers

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/darrelhong/micro-url/store"
	"github.com/darrelhong/micro-url/utils"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

func HandleGhCallback(oauth2Conf *oauth2.Config,
	sessionStore *sessions.CookieStore,
	userStore store.UserStore,
	tursoApiClient *utils.TursoApiClient,
) http.Handler {
	creatingUserTmpl := template.Must(template.ParseFiles("templates/partials/creating_user.html"))
	redirectIndexTmpl := template.Must(template.ParseFiles("templates/partials/redirect_index.html"))

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

		oauthClient := oauth2Conf.Client(r.Context(), tok)

		resp, err := oauthClient.Get("https://api.github.com/user/emails")

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

		primaryEmailStruct := emails[0]

		if !primaryEmailStruct.Verified {
			http.Error(w, "Email not verified with GitHub", http.StatusBadRequest)
			return
		}

		session, err = sessionStore.Get(r, "login")
		if err != nil {
			http.Error(w, "Something went wrong, please try again", http.StatusInternalServerError)
			return
		}

		primaryEmail := primaryEmailStruct.Email

		session.Values["email"] = primaryEmail
		session.Save(r, w)

		_, err = userStore.GetUser(primaryEmail)

		if err == store.ErrUserNotFound {
			log.Println("User not found, creating new user and database")

			creatingUserTmpl.Execute(w, nil)
			w.(http.Flusher).Flush()

			dbUrl, err := tursoApiClient.CreateTenantDatabaseAndRunMigrations(primaryEmail)

			if err != nil {
				log.Println("gg error creating databse or running migration", err)
				redirectIndexTmpl.Execute(w, nil)
				return
			}

			err = userStore.CreateUser(primaryEmail, dbUrl)

			if err != nil {
				log.Println("gg error creating user", err)
				redirectIndexTmpl.Execute(w, nil)
				return
			}

			redirectIndexTmpl.Execute(w, nil)
			return
		}

		if err != nil && err != store.ErrUserNotFound {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	})
}
