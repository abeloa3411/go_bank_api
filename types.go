package main

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Number int64 	`json:"number"`
	Password string 	`json:"password"`
}

type CreateAccountRequest struct {
	FirstName string	`json:"firstname"`
	LastName string		`json:"lastname"`
	Password string		`json:"password"`
}

type TranferRequest struct {
	ToAccount string  `json:"toAccount"`
	Amount int  `json:"amount"`
}

type Account struct {
	ID        int	`json:"id"`
	FirstName string	`json:"firstname"`
	LastName  string	`json:"lastname"`
	EncryptedPassword	string		`json:"-"`
	Number    int64		`json:"number"`
	Balance   int64		`json:"balance"`
	CreatedAt 	time.Time 	`json:"createdAt"`
}

func NewAccount(firstName, lastName, password string) (*Account, error) {
	encpass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	return &Account{
		FirstName: firstName,
		LastName:  lastName,
		EncryptedPassword: string(encpass),
		Number:   int64(rand.Intn(100000)),
		CreatedAt: time.Now().UTC(),
	}, nil
}