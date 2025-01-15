package repository

import (
	"context"
	"database/sql"

	"github.com/ariefro/threads-server/internal/entity"
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
