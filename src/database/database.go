package database

import (
	"database/sql"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/rds/rdsutils"
	_ "github.com/lib/pq"
)

type DbUser struct {
	Id           int       `db:"id"`
	Username     string    `db:"username"`
	Password     string    `db:"password"`
	Email        string    `db:"email"`
	CreationDate time.Time `db:"creation_date"`
}

type SqlDatabase interface {
	CreateUser(user DbUser) (int, error)
	DeleteUserById(id int) error
	GetUserById(id int) (DbUser, error)
	GetUserByUsername(username string) (DbUser, error)
}

type PostgresDatabase struct {
	Db *sql.DB
}

type DatabaseConfig struct {
	Host   string
	Port   int
	Name   string
	User   string
	Region string
}

// Creates a new database handle.
func New(cfg DatabaseConfig) (*PostgresDatabase, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	credentialsPath := path.Join(homeDir, ".aws", "credentials")
	creds := credentials.NewSharedCredentials(credentialsPath, "default")
	endpoint := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	authToken, err := rdsutils.BuildAuthToken(endpoint, cfg.Region, cfg.User, creds)
	if err != nil {
		return nil, err
	}

	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", cfg.Host, cfg.Port, cfg.User, authToken, cfg.Name)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	return &PostgresDatabase{
		Db: db,
	}, nil
}

func (db *PostgresDatabase) Close() error {
	return db.Db.Close()
}

// Inserts a given user into the database.
func (db *PostgresDatabase) CreateUser(user DbUser) (int, error) {
	rows, err := db.Db.Query(QUERY_CREATE_USER, user.Username, user.Password, user.Email, user.CreationDate)
	if err != nil {
		return -1, err
	}

	var id int
	err = rows.Scan(&id)
	return id, err
}

// Deletes a user identified by the given id.
func (db *PostgresDatabase) DeleteUserById(id int) error {
	_, err := db.Db.Query(QUERY_DELETE_USER_BY_ID, id)
	return err
}

// Returns a user identified by its identifier.
func (db *PostgresDatabase) GetUserById(id int) (DbUser, error) {
	row := db.Db.QueryRow(QUERY_SELECT_USER_BY_ID, id)

	var user DbUser
	err := row.Scan(&user)

	return user, err
}

// Returns a user identified by its username.
func (db *PostgresDatabase) GetUserByUsername(username string) (DbUser, error) {
	row := db.Db.QueryRow(QUERY_SELECT_USER_BY_USERNAME, username)

	var user DbUser
	err := row.Scan(&user.Id, &user.Username, &user.Password, &user.Email, &user.CreationDate)

	return user, err
}
