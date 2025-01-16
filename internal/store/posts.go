package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ariefro/threads-server/internal/entity"
	"github.com/ariefro/threads-server/internal/query"
	"github.com/lib/pq"
)

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
	Create(context.Context, *entity.Post) error
}

func (r *postStorage) Create(ctx context.Context, post *entity.Post) error {
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
