package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/darrelhong/micro-url/store"
)

func HandleShorten(urlStore store.UrlStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.FormValue("url")
		log.Println("URL to shorten:", url)

		shortUrlId, err := urlStore.CreateShortLink(url)

		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}

		log.Println("Shortened URL", shortUrlId)

		fmt.Fprintf(w, "Shortened URL: %s", shortUrlId)
	})
}
