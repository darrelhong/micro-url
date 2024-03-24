package utils

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/darrelhong/micro-url/store"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type UserDbClient struct {
	dbToken string
}

func NewUserDbClient(dbToken string) *UserDbClient {
	return &UserDbClient{dbToken}
}

func (userDbClient *UserDbClient) GetUserUrlStore(dbUrl string) (store.UrlStore, error) {

	tenantDb, err := sql.Open("libsql", fmt.Sprintf("%s?authToken=%s", dbUrl, userDbClient.dbToken))
	if err != nil {
		log.Println("oh no error opening tenant db", err)
		return nil, err
	}

	tenantDbUrlStore := store.NewDbUrlStore(tenantDb)

	return tenantDbUrlStore, nil
}
