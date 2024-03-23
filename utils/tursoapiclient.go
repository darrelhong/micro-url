package utils

import (
	"bytes"
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

//go:embed create_urls_table.sql
var createUsersTableMigration string

type TursoApiClient struct {
	tursoApiUrl   string
	tursoOrgName  string
	tursoApiToken string
	dbToken       string
	client        *http.Client
}

func NewTursoApiClient(
	tursoApiUrl string,
	tursoOrgName string,
	tursoApiToken string,
	dbToken string,
) *TursoApiClient {
	return &TursoApiClient{
		tursoApiUrl:   tursoApiUrl,
		tursoOrgName:  tursoOrgName,
		tursoApiToken: tursoApiToken,
		dbToken:       dbToken,
		client:        &http.Client{},
	}
}

func (c *TursoApiClient) CreateTenantDatabaseAndRunMigrations(email string) (string, error) {
	username := StripNonAlphaNumeric(strings.Split(email, "@")[0])
	dbName := "micro-url-db-" + username

	log.Println("Database name:", dbName)

	// create turso database
	body := []byte(`{
		"group":"default",
		"name":"` + dbName + `" 
		}`)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/organizations/%s/databases", c.tursoApiUrl, c.tursoOrgName), bytes.NewBuffer(body))

	if err != nil {
		log.Println("how did request creation fail?", err)
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+c.tursoApiToken)

	resp, err := c.client.Do(req)

	if err != nil {
		log.Println("woops turso db creation failed", err)
		return "", err
	}

	type CreateDatabaseResponse struct {
		Database struct {
			DbId     string `json:"DbId"`
			Hostname string `json:"Hostname"`
			Name     string `json:"Name"`
		} `json:"database"`
	}

	defer resp.Body.Close()

	createDatabaseResponse := CreateDatabaseResponse{}
	err = json.NewDecoder(resp.Body).Decode(&createDatabaseResponse)

	if err != nil {
		log.Println("somehow json is wrong", err)
		return "", err
	}

	dbUrl := "https://" + createDatabaseResponse.Database.Hostname

	log.Println("Database URL:", dbUrl)

	db, err := sql.Open("libsql", fmt.Sprintf("%s?authToken=%s", dbUrl, c.dbToken))

	if err != nil {
		log.Println("damn it cannot connect to newly created db", err)
		return "", err
	}

	// give turso some time to ready the database
	time.Sleep(10 * time.Second)

	log.Println("Running migrations")
	_, err = db.Exec(createUsersTableMigration)

	if err != nil {
		log.Println("shit migration on new db failed", err)
		return "", err
	}

	return dbUrl, nil
}
