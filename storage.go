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
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore () (*PostgresStore, error){
	connectStr := "user=postgres dbname=bank host=localhost port=5432 password=TU01-BE213-0634/2019 sslmode=disable"
	db, err := sql.Open("postgres", connectStr)

	if err != nil{
		return nil, err
	}

	if err := db.Ping(); err != nil{
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}


func(s *PostgresStore) Init() error{
	return s.CreateAccountTable()
}


func(s *PostgresStore) CreateAccountTable() error{
	query := `CREATE TABLE if not exists account (
		id serial primary key,
		firstname VARCHAR(50),
		lastname VARVCHAR(50),
		number serial,
		balance serial,
		created_at timestamp
	);
	`

	_, err := s.db.Exec(query)

	return err
}



func(s *PostgresStore) CreateAccount(acc *Account) error {

	query := `insert into account
	(firstname,lastname,number,balance, created_at)
	values ($1, $2, $3, $4, $5)`

	resp, err := s.db.Query(query, 
				acc.FirstName,
				acc.LastName,
				acc.Number,
				acc.Balance,
				acc.CreatedAt,
	)

	if err != nil {
		return err 
	}

	fmt.Printf("%+v", resp)

	return nil
}

func(s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

func(s *PostgresStore) DeleteAccount(id int) error {
	_, err := s.db.Query("DELETE FROM account WHERE id = $1", id)
	return err
}

func(s *PostgresStore) GetAccounts() ([]*Account,error) {

	rows, err := s.db.Query("SELECT * FROM account")

	if err != nil{
		return nil, err
	}

	accounts := []*Account{}

	for rows.Next(){
		account, err := ScanIntoAccount(rows)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}
	return accounts, nil
}

func(s *PostgresStore) GetAccountByID(id int) (*Account, error) {

	rows, err := s.db.Query("SELECT * FROM account WHERE id = $1", id)

	if err != nil{
		return nil,err
	}

	for rows.Next(){
		return ScanIntoAccount(rows)
	}
	return nil, fmt.Errorf("Account %d not found", id)
}

func ScanIntoAccount(rows *sql.Rows)(*Account, error){
	account := new(Account)
	rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt,
		)

	return account, nil
}