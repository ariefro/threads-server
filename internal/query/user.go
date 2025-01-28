package query

const (
	CreateUser = `
		INSERT INTO users (username, email, password, role_id)
		VALUES ($1, $2, $3, (SELECT id FROM roles WHERE name = $4))
		RETURNING id, created_at
	`

	GetUserByID = `
		SELECT users.id, username, email, password, created_at, updated_at, roles.*
		FROM users
		JOIN roles ON (users.role_id = roles.id)
		WHERE users.id = $1 AND is_active = true
	`

	GetUserByEmail = `
		SELECT id, username, email, password, created_at, updated_at
		FROM users
		WHERE email = $1 AND is_active = true
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

	DeleteUserById = `
		DELETE FROM users WHERE id = $1
	`
)
