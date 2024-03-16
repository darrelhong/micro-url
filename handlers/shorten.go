package handlers

import (
	"fmt"
	"net/http"

	"github.com/darrelhong/micro-url/store"
)

func HandleShorten(urlStore store.UrlStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.FormValue("url")
		fmt.Println("URL to shorten:", url)

		shortUrl, err := urlStore.CreateShortLink(url)

		fmt.Println("Shortened URL", shortUrl)

		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
	})
}
