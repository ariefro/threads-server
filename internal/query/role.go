package query

const (
	GetRoleByName = `
		SELECT id, name, description, level
		FROM roles
		WHERE name = $1
	`
)
