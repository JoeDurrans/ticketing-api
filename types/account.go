package types

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Role string

const (
	RoleAdmin  Role = "admin"
	RoleUser   Role = "user"
	RoleEditor Role = "editor"
)

type Account struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (a *Account) CheckPasswordHash(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password))
	if err != nil {
		return &Unauthorized{
			Message: "invalid password",
		}
	}

	return nil
}

func CreateAccount(username string, password string, role Role) (*Account, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error creating account")
	}

	return &Account{
		Username:  username,
		Password:  string(bytes),
		Role:      role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}
