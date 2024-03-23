package handlers

import (
	"html/template"
	"net/http"
)

func HandleIndex(ghClientId string) http.Handler {
	tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/index.html", "templates/partials/head.html"))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.ExecuteTemplate(w, "base", struct {
			GhClientId string
		}{
			GhClientId: ghClientId,
		})

		if err != nil {
			http.Error(w, "Something went wrong executing template", http.StatusInternalServerError)
		}
	})
}
