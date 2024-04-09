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

func (d *AccountAdapter) Create(account *types.Account) (*types.Account, error) {
	id := 0
	err := d.db.QueryRow("INSERT INTO account (username, password, role) VALUES ($1, $2, $3) RETURNING id", account.Username, account.Password, account.Role).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("error creating account: %w", err)
	}

	account.ID = id

	return account, nil
}

func (d *AccountAdapter) Get() ([]*types.Account, error) {
	rows, err := d.db.Query(`SELECT * FROM account`)
	if err != nil {
		return nil, fmt.Errorf("error getting accounts: %w", err)
	}

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

func (d *AccountAdapter) GetByID(id int) (*types.Account, error) {
	rows, err := d.db.Query(`SELECT * FROM account WHERE id = $1`, id)
	if err != nil {
		return nil, fmt.Errorf("error getting account: %w", err)
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account %d not found", id)
}

func (d *AccountAdapter) GetByUsername(username string) (*types.Account, error) {
	rows, err := d.db.Query(`SELECT * FROM account WHERE username = $1`, username)
	if err != nil {
		return nil, fmt.Errorf("error getting account: %w", err)
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account with the username: %s not found", username)
}

func (d *AccountAdapter) Update(account *types.Account) (*types.Account, error) {
	_, err := d.db.Exec(`UPDATE account SET username = $1, password = $2, role = $3 WHERE id = $4`, account.Username, account.Password, account.Role, account.ID)
	if err != nil {
		return nil, fmt.Errorf("error updating account: %w", err)
	}

	return account, nil
}

func (d *AccountAdapter) Delete(id int) error {
	_, err := d.db.Exec(`DELETE FROM account WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("error deleting account: %w", err)
	}

	return nil
}

func scanIntoAccount(rows *sql.Rows) (*types.Account, error) {
	account := &types.Account{}

	err := rows.Scan(&account.ID, &account.Username, &account.Password, &account.Role, &account.CreatedAt, &account.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("error scanning account: %w", err)
	}

	return account, nil
}
