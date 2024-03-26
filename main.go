package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/darrelhong/micro-url/store"
	"github.com/darrelhong/micro-url/utils"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/sync/errgroup"
)

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

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
		return err
	}

	urlStore := store.NewDbUrlStore(db)
	userStore := store.NewDbUserStore(db)

	tursoApiUrl := os.Getenv("TURSO_API_URL")
	tursoOrgName := os.Getenv("TURSO_ORG_NAME")
	tursoApiToken := os.Getenv("TURSO_API_TOKEN")

	tursoApiClient := utils.NewTursoApiClient(tursoApiUrl, tursoOrgName, tursoApiToken, dbToken)

	userDbClient := utils.NewUserDbClient(dbToken)

	srv := NewServer(urlStore, oauth2Conf, sessionStore, userStore, tursoApiClient, userDbClient)

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: srv,
	}

	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		log.Println("Starting server on :8080")

		return httpServer.ListenAndServe()
	})

	eg.Go(func() error {
		<-egCtx.Done()
		log.Println("Shutting down server")

		shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		return httpServer.Shutdown(shutCtx)
	})

	err = eg.Wait()

	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func main() {
	ctx := context.Background()

	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

}
