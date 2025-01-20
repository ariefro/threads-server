package query

const (
	DeleteUserInvitation = `
		DELETE FROM user_invitations
		WHERE user_id = $1
	`
)
