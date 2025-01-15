package repository

import (
	"context"
	"database/sql"

	"github.com/ariefro/threads-server/internal/entity"
)

type userRepository struct {
	db *sql.DB
}

type UserRepository interface {
	Create(context.Context, *entity.User) error
}
