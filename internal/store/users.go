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
	GetByID(context.Context, int64) (*User, error)
}

func (r *userStorage) Create(ctx context.Context, user *User) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

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

func (s *userStorage) GetByID(ctx context.Context, userID int64) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := &User{}
	err := s.db.QueryRowContext(
		ctx,
		query.GetUserByID,
		userID,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}
