package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ariefro/threads-server/internal/query"
)

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

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
	Create(context.Context, *User) error
}

func (r *userStorage) Create(ctx context.Context, user *User) error {
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
