package handlers

import (
	"net/http"

	"github.com/darrelhong/micro-url/store"
)

func HandleRedirect(urlStore store.UrlStore) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		shortUrlId := r.PathValue("shortUrlId")

		originalUrl, err := urlStore.GetOriginalUrl(shortUrlId)

		if originalUrl == "" {

			if err != nil {
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
				return
			}

			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		http.Redirect(w, r, originalUrl, http.StatusFound)
	})
}
