package store

import (
	"context"
	"database/sql"

	"github.com/ariefro/threads-server/internal/query"
)

// NewRoleStorage creates a new instance of RoleStorage implementation
func NewRoleStorage(db *sql.DB) RoleStorage {
	return &roleStorage{
		db: db,
	}
}

type roleStorage struct {
	db *sql.DB
}

type RoleStorage interface {
	GetByName(context.Context, string) (*Role, error)
}

type Role struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Level       int    `json:"level"`
}

func (s *roleStorage) GetByName(ctx context.Context, slug string) (*Role, error) {
	role := &Role{}
	err := s.db.QueryRowContext(ctx, query.GetRoleByName, slug).
		Scan(&role.ID, &role.Name, &role.Description, &role.Level)
	if err != nil {
		return nil, err
	}

	return role, nil
}
