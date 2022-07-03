package repository

import (
	"database/sql"
	"flhansen/fitter-login-service/src/database"
	"time"

	_ "github.com/lib/pq"
)

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

func NewAccountRepository(host string, port int, username string, password string, databaseName string) (AccountRepository, error) {
	dsn := database.DataSourceName(host, port, username, password, databaseName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return &accountRepository{
		db: db,
	}, nil
}

func (repo *accountRepository) CreateAccount(account Account) (int, error) {
	row := repo.db.QueryRow(QUERY_CREATE_ACCOUNT, account.Username, account.Password, account.Email, account.CreationDate)

	var id int
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
	_, err := repo.db.Exec(QUERY_DELETE_ACCOUNT_BY_ID, id)
	return err
}

func (repo *accountRepository) DeleteAccounts() error {
	_, err := repo.db.Exec(QUERY_DELETE_ACCOUNTS)
	return err
}
