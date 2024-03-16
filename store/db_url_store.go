package store

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"strings"
)

type UrlStore interface {
	CreateShortLink(url string) (string, error)
	GetOriginalUrl(shortURLId string) (string, error)
}

type DbUrlStore struct {
	db *sql.DB
}

func NewDbUrlStore(db *sql.DB) *DbUrlStore {
	return &DbUrlStore{db}
}

const retries = 3

func (store *DbUrlStore) CreateShortLink(url string) (string, error) {
	for range retries {
		shortURLId := generateShortId()

		_, err := store.db.Exec("INSERT INTO urls (original_url, short_url_id) VALUES (?, ?)", url, shortURLId)

		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed: urls.short_url_id") {
				log.Printf("Short URL ID %s already exists, trying again", shortURLId)
				continue
			}

			log.Println("Error inserting URL into database:", err)
			return "", err
		}

		return shortURLId, nil
	}

	log.Printf("failed to generate unique short URL after %d attempts", retries)
	return "", fmt.Errorf("failed to generate unique short URL after %d attempts", retries)
}

func (store *DbUrlStore) GetOriginalUrl(shortURLId string) (string, error) {
	var originalUrl string

	err := store.db.QueryRow("SELECT original_url FROM urls WHERE short_url_id = ?", shortURLId).Scan(&originalUrl)

	if err == sql.ErrNoRows {
		return "", nil
	}

	if err != nil {
		log.Println("Error getting original URL from database:", err)

		return "", err
	}

	return originalUrl, nil

}

const availableChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateShortId() string {
	result := make([]byte, 8)
	for i := range result {
		result[i] = availableChars[rand.Intn(len(availableChars))]
	}
	return string(result)
}
