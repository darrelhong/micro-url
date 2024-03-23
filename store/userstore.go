package store

import (
	"database/sql"
	"errors"
)

type User struct {
	Id    int
	Email string
	DbUrl string
}

type UserStore interface {
	GetUser(email string) (User, error)
}

type DbUserStore struct {
	db *sql.DB
}

var ErrUserNotFound = errors.New("User not found")

func NewDbUserStore(db *sql.DB) *DbUserStore {
	return &DbUserStore{db}
}

func (store *DbUserStore) GetUser(email string) (User, error) {
	var user User

	err := store.db.QueryRow("SELECT id, email, db_url FROM users WHERE email = ?", email).Scan(&user.Id, &user.Email, &user.DbUrl)

	if err == sql.ErrNoRows {
		return User{}, ErrUserNotFound
	}

	if err != nil {
		return User{}, err
	}

	return user, nil
}
