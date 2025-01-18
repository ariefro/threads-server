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

	GetUserFeed = `
		SELECT p.id, p.user_id, p.title, p.content, p.version, p.tags, p.created_at, p.updated_at, u.username, COUNT(c.id) AS comments_count
		FROM posts p
		LEFT JOIN comments c ON c.post_id = p.id
		LEFT JOIN users u ON p.user_id = u.id
		JOIN followers f ON f.follower_id = p.user_id OR p.user_id = $1
		WHERE 
			f.user_id = $1 AND
			(p.title ILIKE '%' || $4 || '%' OR p.content ILIKE '%' || $4 || '%') AND
			(p.tags @> $5 OR $5 = '{}')
		GROUP BY p.id, u.username
		ORDER BY p.created_at %s
		LIMIT $2 OFFSET $3
	`
)
