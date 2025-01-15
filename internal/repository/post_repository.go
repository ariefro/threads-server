package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ariefro/threads-server/internal/entity"
	"github.com/ariefro/threads-server/internal/query"
	"github.com/lib/pq"
)

// NewPostRepository creates a new instance of PostRepository implementation
func NewPostRepository(db *sql.DB) PostRepository {
	return &postRepository{
		db: db,
	}
}

type postRepository struct {
	db *sql.DB
}

type PostRepository interface {
	Create(context.Context, *entity.Post) error
}

func (r *postRepository) Create(ctx context.Context, post *entity.Post) error {
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
