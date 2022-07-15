package repository

import (
	"database/sql"
	"flhansen/fitter-login-service/src/database"
	"time"

	_ "github.com/lib/pq"
)

type DatabaseConfig struct {
	Host         string
	Port         int
	Username     string
	Password     string
	DatabaseName string
}

type Account struct {
	Id           int
	Username     string
	Password     string
	Email        string
	CreationDate time.Time
}

type AccountRepository interface {
	CreateAccount(account Account) (int, error)
	GetAccountById(id int) (Account, error)
	GetAccountByUsername(username string) (Account, error)
	DeleteAccountById(id int) error
	DeleteAccounts() error
}

type accountRepository struct {
	db *sql.DB
}

func NewAccountRepository(config DatabaseConfig) (AccountRepository, error) {
	dsn := database.DataSourceName(config.Host, config.Port, config.Username, config.Password, config.DatabaseName)
	db, _ := sql.Open("postgres", dsn)

	return &accountRepository{
		db: db,
	}, nil
}

func (repo *accountRepository) CreateAccount(account Account) (int, error) {
	row := repo.db.QueryRow(QUERY_CREATE_ACCOUNT, account.Username, account.Password, account.Email, account.CreationDate)

	id := -1
	err := row.Scan(&id)
	return id, err
}

func (repo *accountRepository) GetAccountById(id int) (Account, error) {
	row := repo.db.QueryRow(QUERY_SELECT_ACCOUNT_BY_ID, id)

	var account Account
	err := row.Scan(&account.Id, &account.Username, &account.Password, &account.Email, &account.CreationDate)
	return account, err
}

func (repo *accountRepository) GetAccountByUsername(username string) (Account, error) {
	row := repo.db.QueryRow(QUERY_SELECT_ACCOUNT_BY_USERNAME, username)

	var account Account
	err := row.Scan(&account.Id, &account.Username, &account.Password, &account.Email, &account.CreationDate)
	return account, err
}

func (repo *accountRepository) DeleteAccountById(id int) error {
	row := repo.db.QueryRow(QUERY_DELETE_ACCOUNT_BY_ID, id)

	deletedId := -1
	err := row.Scan(&deletedId)
	return err
}

func (repo *accountRepository) DeleteAccounts() error {
	_, err := repo.db.Exec(QUERY_DELETE_ACCOUNTS)
	return err
}
