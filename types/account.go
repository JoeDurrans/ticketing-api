package types

import (
	"time"
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

func CreateAccount(username string, password string, role Role) (*Account, error) {
	return &Account{
		Username:  username,
		Password:  password,
		Role:      role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}
