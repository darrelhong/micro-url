package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/darrelhong/micro-url/store"
	"github.com/darrelhong/micro-url/utils"
	"github.com/gorilla/sessions"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func HandleShorten(
	urlStore store.UrlStore,
	sessionStore sessions.Store,
	userStore store.UserStore,
	userDbClient *utils.UserDbClient,
) http.Handler {
	domainName := os.Getenv("DOMAIN_NAME")

	tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/created.html", "templates/partials/head.html"))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		urlToShorten := r.FormValue("url")
		log.Println("URL to shorten:", urlToShorten)

		if urlToShorten == "" {
			http.Error(w, "URL is required", http.StatusBadRequest)
			return
		}

		_, err := url.ParseRequestURI(urlToShorten)
		if err != nil {
			http.Error(w, "URL is not valid", http.StatusBadRequest)
			return
		}

		var shortUrlId string

		loginSession, err := sessionStore.Get(r, "login")

		if err != nil {
			log.Println("Error getting login session, welp skips", err)
		}

		email, ok := loginSession.Values["email"].(string)

		if email != "" && ok {
			user, err := userStore.GetUser(email)

			if err != nil {
				log.Println("erm couldn't get user from email in session?", err)
				http.Error(w, "Something went wrong.", http.StatusUnauthorized)
				return
			}

			userDbUrlStore, err := userDbClient.GetUserUrlStore(user.DbUrl)
			if err != nil {
				log.Println("oh no error opening tenant db", err)
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
				return
			}

			shortUrlId, err = userDbUrlStore.CreateShortLink(urlToShorten)

			if err != nil {
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
				return
			}
		} else {
			shortUrlId, err = urlStore.CreateShortLink(urlToShorten)
		}

		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}

		log.Println("Shortened URL", shortUrlId)

		tmpl.ExecuteTemplate(w, "base", struct {
			ShortUrl string
		}{
			ShortUrl: fmt.Sprintf("%s/%s", domainName, shortUrlId),
		})
	})
}
