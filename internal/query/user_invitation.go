package query

const (
	CreateUserInvitation = `
		INSERT INTO user_invitations (token, user_id, expiry)
		VALUES ($1, $2, $3)
	`

	DeleteUserInvitation = `
		DELETE FROM user_invitations
		WHERE user_id = $1
	`
)
