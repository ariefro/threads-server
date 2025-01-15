package query

const (
	CreatePost = `
		INSERT INTO posts (title, content, user_id, tags)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, update_at
	`
)
