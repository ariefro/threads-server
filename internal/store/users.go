package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ariefro/threads-server/internal/entity"
	"github.com/ariefro/threads-server/internal/query"
)

// NewUserStorage creates a new instance of UserStorage implementation
func NewUserStorage(db *sql.DB) UserStorage {
	return &userStorage{
		db: db,
	}
}

type userStorage struct {
	db *sql.DB
}

type UserStorage interface {
	Create(context.Context, *entity.User) error
}

func (r *userStorage) Create(ctx context.Context, user *entity.User) error {
	err := r.db.QueryRowContext(
		ctx,
		query.CreateUser,
		user.Username,
		user.Email,
		user.Password,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}
