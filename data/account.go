package data

import (
	"database/sql"
	"fmt"
	"ticketing-api/types"

	_ "github.com/lib/pq"
)

type AccountAdapter struct {
	db *sql.DB
}

func CreateAccountAdapter(db *sql.DB) *AccountAdapter {
	return &AccountAdapter{
		db: db,
	}
}

func (a *AccountAdapter) Create(account *types.Account) (*types.Account, error) {
	id := 0
	err := a.db.QueryRow("INSERT INTO account (username, password, role) VALUES ($1, $2, $3) RETURNING id", account.Username, account.Password, account.Role).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("error creating account")
	}

	account.ID = id

	return account, nil
}

func (a *AccountAdapter) Get() ([]*types.Account, error) {
	rows, err := a.db.Query(`SELECT * FROM account`)
	if err != nil {
		return nil, fmt.Errorf("error getting accounts")
	}
	defer rows.Close()

	accounts := []*types.Account{}

	for rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (a *AccountAdapter) GetByID(id int) (*types.Account, error) {
	rows, err := a.db.Query(`SELECT * FROM account WHERE id = $1`, id)
	if err != nil {
		return nil, fmt.Errorf("error getting account")
	}
	defer rows.Close()

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account with id: %d not found", id)
}

func (a *AccountAdapter) GetByUsername(username string) (*types.Account, error) {
	rows, err := a.db.Query(`SELECT * FROM account WHERE username = $1`, username)
	if err != nil {
		return nil, fmt.Errorf("error getting account")
	}
	defer rows.Close()

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account with the username: %s not found", username)
}

func (a *AccountAdapter) Update(account *types.Account) (*types.Account, error) {
	_, err := a.db.Exec(`UPDATE account SET username = $1, password = $2, role = $3 WHERE id = $4`, account.Username, account.Password, account.Role, account.ID)
	if err != nil {
		return nil, fmt.Errorf("error updating account")
	}

	return account, nil
}

func (a *AccountAdapter) Delete(id int) error {
	_, err := a.db.Exec(`DELETE FROM account WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("error deleting account")
	}

	return nil
}

func scanIntoAccount(rows *sql.Rows) (*types.Account, error) {
	account := &types.Account{}

	err := rows.Scan(&account.ID, &account.Username, &account.Password, &account.Role, &account.CreatedAt, &account.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("error reading account")
	}

	return account, nil
}
