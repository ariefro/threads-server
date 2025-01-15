package repository

import (
	"database/sql"
)

type Repositories struct {
	Post PostRepository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Post: &postRepository{db},
	}
}
