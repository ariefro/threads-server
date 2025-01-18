package query

const (
	CreateFollower = `
		INSERT INTO followers (user_id, follower_id)
		VALUES ($1, $2)
	`
)
