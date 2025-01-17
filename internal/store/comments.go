package store

import (
	"context"
	"database/sql"

	"github.com/ariefro/threads-server/internal/query"
)

type Comment struct {
	ID        int64  `json:"id"`
	PostID    int64  `json:"post_id"`
	UserID    int64  `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	User      User   `json:"user"`
}

// NewCommentStorage creates a new instance of CommentStorage implementation
func NewCommentStorage(db *sql.DB) CommentStorage {
	return &commentStorage{
		db: db,
	}
}

type commentStorage struct {
	db *sql.DB
}

type CommentStorage interface {
	GetByPostID(context.Context, int64) ([]Comment, error)
}

func (s *commentStorage) GetByPostID(ctx context.Context, postID int64) ([]Comment, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query.GetCommentsByPostID, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []Comment{}
	for rows.Next() {
		var c Comment
		c.User = User{}
		err := rows.Scan(
			&c.ID,
			&c.PostID,
			&c.UserID,
			&c.Content,
			&c.CreatedAt,
			&c.User.Username,
			&c.User.ID,
		)
		if err != nil {
			return nil, err
		}

		comments = append(comments, c)
	}

	return comments, nil
}
