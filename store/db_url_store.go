package store

import (
	"database/sql"
	"math/rand"
)

type UrlStore interface {
	CreateShortLink(url string) (string, error)
}

type DbUrlStore struct {
	db *sql.DB
}

func NewDbUrlStore(db *sql.DB) *DbUrlStore {
	return &DbUrlStore{db}
}

func (store *DbUrlStore) CreateShortLink(url string) (string, error) {
	shortURLId := generateShortId()
	return shortURLId, nil
}

const availableChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateShortId() string {
	result := make([]byte, 8)
	for i := range result {
		result[i] = availableChars[rand.Intn(len(availableChars))]
	}
	return string(result)
}
