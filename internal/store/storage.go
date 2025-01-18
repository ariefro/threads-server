package store

import (
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("resource not found")
	ErrConflict          = errors.New("resource already exists")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Users     UserStorage
	Posts     PostStorage
	Comments  CommentStorage
	Followers FollowerStorage
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Users:     NewUserStorage(db),
		Posts:     NewPostStorage(db),
		Comments:  NewCommentStorage(db),
		Followers: NewFollowerStorage(db),
	}
}
