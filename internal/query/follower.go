package query

const (
	CreateFollower = `
		INSERT INTO followers (user_id, follower_id)
		VALUES ($1, $2)
	`

	DeleteFollowerByID = `
		DELETE FROM followers 
		WHERE user_id = $1 AND follower_id = $2
	`
)
