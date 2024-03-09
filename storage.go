package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccounts() ([]*Account, error)
	GetAccountByID(int) (*Account, error)
	GetAccountByNumber(int) (*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

// Create and initialize a new instance of PostgresStore struct
func NewPostgresStore() (*PostgresStore, error) {

	//Connection string

	connStr := "user=postgres dbname=postgres password=gobank sslmode=disable"
	//Establish the connection to the PostgreSQL database
	//returns database handle *sql.DB and error
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	//If database is opened successfully, call db.Ping() to verify if database is reachable and responsive
	//This ensures that the database is ready for the queries
	if err := db.Ping(); err != nil {
		return nil, err
	}

	//Database connection and ping is successful, so create a new instance of PostgresStore struct
	//Initialize db field with the *sql.DB obtained earlier
	//Then return this struct with nil (for error)
	return &PostgresStore{
		db: db,
	}, nil
}

// Initialize the PostgresStore instance and call a function to create the necessary database tables
func (s *PostgresStore) Init() error {
	return s.createAccountTable()
}

/*
Create the necessary tables
Exec() returns result after executing the query which is ignored here.
Only errors are returned if encountered
*/
func (s *PostgresStore) createAccountTable() error {
	query := `create table if not exists account (
		id serial primary key,
		first_name varchar(50),
		last_name varchar(50),
		number serial,
		encrypted_password varchar(100),
		balance serial,
		created_at timestamp
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	query := `insert into account 
	(first_name, last_name, number, encrypted_password, balance, created_at)
	values ($1, $2, $3, $4, $5, $6)`
	_, err := s.db.Query(
		query,
		acc.FirstName, acc.LastName, acc.Number, acc.EncryptedPassword, acc.Balance, acc.CreatedAt)

	if err != nil {
		return err
	}

	//fmt.Printf("%+v\n", resp)

	return nil
}

func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	_, err := s.db.Query("delete from account where id = $1", id)
	return err
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query(`select * from account`)
	if err != nil {
		return nil, err
	}

	accounts := []*Account{}
	for rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	rows, err := s.db.Query(`select * from account where id = $1`, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}
	return nil, fmt.Errorf("account %d not found", id)
}

func (s *PostgresStore) GetAccountByNumber(number int) (*Account, error) {
	rows, err := s.db.Query(`select * from account where number = $1`, number)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}
	return nil, fmt.Errorf("account with number %d not found", number)
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.EncryptedPassword,
		&account.Balance,
		&account.CreatedAt)

	return account, err
}
