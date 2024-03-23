package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/darrelhong/micro-url/store"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	oauth2Conf := &oauth2.Config{
		ClientID:     os.Getenv("GH_CLIENT_ID"),
		ClientSecret: os.Getenv("GH_CLIENT_SECRET"),
		RedirectURL:  fmt.Sprintf("%s/github/callback", os.Getenv("DOMAIN_NAME")),
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}

	sessionStore := sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))

	dbUrl := os.Getenv("DB_URL")
	dbToken := os.Getenv("DB_TOKEN")

	db, err := sql.Open("libsql", fmt.Sprintf("%s?authToken=%s", dbUrl, dbToken))

	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	srv := NewServer(store.NewDbUrlStore(db), oauth2Conf, sessionStore)

	log.Fatal(http.ListenAndServe(":8080", srv))
}
