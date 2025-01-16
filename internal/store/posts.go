package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ariefro/threads-server/internal/query"
	"github.com/lib/pq"
)

type Post struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserID    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewPostStorage creates a new instance of PostStorage implementation
func NewPostStorage(db *sql.DB) PostStorage {
	return &postStorage{
		db: db,
	}
}

type postStorage struct {
	db *sql.DB
}

type PostStorage interface {
	Create(context.Context, *Post) error
}

func (r *postStorage) Create(ctx context.Context, post *Post) error {
	err := r.db.QueryRowContext(
		ctx,
		query.CreatePost,
		post.Title,
		post.Content,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}

	return nil
}
