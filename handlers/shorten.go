package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
)

func HandleShorten() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.FormValue("url")
		fmt.Println("URL to shorten:", url)
		shortURLId := generateShortId()
		fmt.Println("Shortened URL ID:", shortURLId)
	})
}

const availableChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateShortId() string {
	result := make([]byte, 8)
	for i := range result {
		result[i] = availableChars[rand.Intn(len(availableChars))]
	}
	return string(result)
}
