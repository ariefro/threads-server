package query

const (
	GetCommentsByPostID = `
		SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, c.updated_at, users.username, users.email, users.id, users.created_at, users.updated_at
		FROM comments c
		JOIN users on users.id = c.user_id
		WHERE c.post_id = $1
		ORDER BY c.created_at DESC
	`

	CreateComment = `
		INSERT INTO comments (post_id, user_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`
)
