package store

import (
	"database/sql"
)

type Storage struct {
	User UserStorage
	Post PostStorage
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		User: NewUserStorage(db),
		Post: NewPostStorage(db),
	}
}
