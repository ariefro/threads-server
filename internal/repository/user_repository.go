package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ariefro/threads-server/internal/entity"
	"github.com/ariefro/threads-server/internal/query"
)

// NewUserRepository creates a new instance of UserRepository implementation
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

type userRepository struct {
	db *sql.DB
}

type UserRepository interface {
	Create(context.Context, *entity.User) error
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
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
