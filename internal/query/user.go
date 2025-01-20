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

	GetUserByInvitation = `
		SELECT u.id, u.username, u.email, u.is_active, u.created_at, u.updated_at
		FROM users u
		JOIN user_invitations ui ON u.id = ui.user_id
		WHERE ui.token = $1 AND ui.expiry > $2
	`

	UpdateUser = `
		UPDATE users
		SET username = $1, email = $2, is_active = $3
		WHERE id = $4
	`
)
