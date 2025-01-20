package store

import (
	"context"
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

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
