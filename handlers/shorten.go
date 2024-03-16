package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/darrelhong/micro-url/store"
)

func HandleShorten(urlStore store.UrlStore) http.Handler {
	domainName := os.Getenv("DOMAIN_NAME")

	tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/created.html", "templates/partials/head.html"))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.FormValue("url")
		log.Println("URL to shorten:", url)

		shortUrlId, err := urlStore.CreateShortLink(url)

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
