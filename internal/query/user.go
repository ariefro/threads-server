package query

const (
	CreateUser = `
		INSERT INTO users (username, email, password)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	GetUserByID = `
		SELECT id, username, email, password, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	CreateUserInvitation = `
		INSERT INTO user_invitations (token, user_id, expiry)
		VALUES ($1, $2, $3)
	`
)
