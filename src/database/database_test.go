package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDataSourceName(t *testing.T) {
	// given
	host := "localhost"
	port := 8000
	username := "username"
	password := "password"
	database := "test"

	// when
	dsn := DataSourceName(host, port, username, password, database)

	// then
	assert.Equal(t, "host=localhost port=8000 user=username password=password dbname=test sslmode=disable", dsn)
}
