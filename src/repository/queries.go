package repository

const (
	QUERY_DELETE_ACCOUNTS = `
	DELETE FROM account`

	QUERY_CREATE_ACCOUNT_TABLE = `
	CREATE TABLE account (
		id SERIAL PRIMARY KEY,
		username VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(64) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		creation_date TIMESTAMP WITH TIME ZONE DEFAULT now()
	)`

	QUERY_CREATE_ACCOUNT = `
	INSERT INTO Account (username, password, email, creation_date)
	VALUES ($1, $2, $3, $4)
	RETURNING id`

	QUERY_SELECT_ACCOUNT_BY_ID = `
	SELECT id, username, password, email, creation_date
	FROM Account
	WHERE id = $1
	LIMIT 1`

	QUERY_SELECT_ACCOUNT_BY_USERNAME = `
	SELECT id, username, password, email, creation_date
	FROM Account
	WHERE username = $1
	LIMIT 1`

	QUERY_DELETE_ACCOUNT_BY_ID = `
	DELETE FROM Account
	WHERE id = $1
	RETURNING id`
)
