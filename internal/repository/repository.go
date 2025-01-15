package repository

import (
	"database/sql"
)

type Repositories struct {
	User UserRepository
	Post PostRepository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		User: NewUserRepository(db),
		Post: NewPostRepository(db),
	}
}
