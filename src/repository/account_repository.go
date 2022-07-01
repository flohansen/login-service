package repository

import (
	"database/sql"
	"flhansen/fitter-login-service/src/database"
	"time"
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
}

type accountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) AccountRepository {
	return &accountRepository{
		db: db,
	}
}

func (repo *accountRepository) CreateAccount(account Account) (int, error) {
	rows, err := repo.db.Query(database.QUERY_CREATE_USER, account.Username, account.Password, account.CreationDate)
	if err != nil {
		return -1, err
	}

	var id int
	err = rows.Scan(&id)
	return id, err
}

func (repo *accountRepository) GetAccountById(id int) (Account, error) {
	row := repo.db.QueryRow(database.QUERY_SELECT_USER_BY_ID, id)

	var account Account
	err := row.Scan(&account.Id, &account.Username, &account.Password, &account.Email, &account.CreationDate)
	return account, err
}

func (repo *accountRepository) GetAccountByUsername(username string) (Account, error) {
	row := repo.db.QueryRow(database.QUERY_SELECT_USER_BY_USERNAME, username)

	var account Account
	err := row.Scan(&account.Id, &account.Username, &account.Password, &account.Email, &account.CreationDate)
	return account, err
}

func (repo *accountRepository) DeleteAccountById(id int) error {
	_, err := repo.db.Exec(database.QUERY_DELETE_USER_BY_ID, id)
	return err
}
