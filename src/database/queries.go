package database

const (
	QUERY_CREATE_USER = `
	INSERT INTO Account (username, password, email, creation_date)
	VALUES ($1, $2, $3, $4)
	RETURNING id`

	QUERY_SELECT_USER_BY_ID = `
	SELECT id, username, password, email, creation_date
	FROM Account
	WHERE id = $1
	LIMIT 1`

	QUERY_SELECT_USER_BY_USERNAME = `
	SELECT id, username, password, email, creation_date
	FROM Account
	WHERE username = $1
	LIMIT 1`

	QUERY_DELETE_USER_BY_ID = `
	DELETE FROM Account
	WHERE id = $1`
)
