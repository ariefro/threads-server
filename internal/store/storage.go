package store

import (
	"database/sql"
)

type Storage struct {
	Users UserStorage
	Posts PostStorage
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Users: NewUserStorage(db),
		Posts: NewPostStorage(db),
	}
}
