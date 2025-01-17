package query

const (
	CreatePost = `
		INSERT INTO posts (title, content, user_id, tags)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	GetPostByID = `
		SELECT id, user_id, title, content, tags, version, created_at, updated_at
		FROM posts
		WHERE id = $1
	`

	DeletePostByID = `
		DELETE FROM posts
		WHERE id = $1
	`

	UpdatePostByID = `
		UPDATE posts
		SET title = $1, content = $2, version = version + 1
		WHERE id = $3 AND version = $4
		RETURNING version
	`
)
