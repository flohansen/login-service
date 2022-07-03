package database

import "fmt"

func DataSourceName(host string, port int, username string, password string, database string) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, username, password, database)
}
