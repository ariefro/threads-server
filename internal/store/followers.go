package store

import (
	"context"
	"database/sql"

	"github.com/ariefro/threads-server/internal/query"
	"github.com/lib/pq"
)

type Follower struct {
	UserID     int64  `json:"user_id"`
	FollowerID int64  `json:"follower_id"`
	CreatedAt  string `json:"created_at"`
}

// NewFollowerStorage creates a new instance of FollowerStorage implementation
func NewFollowerStorage(db *sql.DB) FollowerStorage {
	return &followeStorage{
		db: db,
	}
}

type followeStorage struct {
	db *sql.DB
}

type FollowerStorage interface {
	Follow(ctx context.Context, followerID, userID int64) error
	Unfollow(ctx context.Context, followerID, userID int64) error
}

func (s *followeStorage) Follow(ctx context.Context, followerID, userID int64) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query.CreateFollower, userID, followerID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return ErrConflict
		}
	}

	return nil
}

func (s *followeStorage) Unfollow(ctx context.Context, followerID, userID int64) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query.DeleteFollowerByID, userID, followerID)
	return err
}
