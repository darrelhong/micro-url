package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/darrelhong/micro-url/store"
	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUrl := os.Getenv("DB_URL")
	dbToken := os.Getenv("DB_TOKEN")

	db, err := sql.Open("libsql", fmt.Sprintf("%s?authToken=%s", dbUrl, dbToken))

	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	srv := NewServer(store.NewDbUrlStore(db), os.Getenv("GH_CLIENT_ID"))

	log.Fatal(http.ListenAndServe(":8080", srv))
}
