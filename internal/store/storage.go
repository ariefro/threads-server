package store

import (
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("resource not found")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Users    UserStorage
	Posts    PostStorage
	Comments CommentStorage
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Users:    NewUserStorage(db),
		Posts:    NewPostStorage(db),
		Comments: NewCommentStorage(db),
	}
}
