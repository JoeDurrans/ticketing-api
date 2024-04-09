package api

import (
	"fmt"
	"net/http"
	"ticketing-api/auth"
	"ticketing-api/types"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	req := CreateAccountRequest{}

	err := DecodeRequest(r, &req)
	if err != nil {
		return err
	}

	account, err := CreateAccount(req.Username, req.Password, req.Role)
	if err != nil {
		return err
	}

	account, err = s.db.Account.Create(account)
	if err != nil {
		return err
	}

	return EncodeResponse(w, http.StatusOK, &APIResponse{Status: http.StatusOK, Message: "account created", Data: account})
}

func (s *APIServer) handleGetAccounts(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.db.Account.Get()
	if err != nil {
		return err
	}

	return EncodeResponse(w, http.StatusOK, &APIResponse{Status: http.StatusOK, Message: "accounts found", Data: accounts})
}

func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	if ok := auth.AccountIDAuth(r, id, types.RoleAdmin, types.RoleEditor); !ok {
		return fmt.Errorf("permission denied")
	}

	account, err := s.db.Account.GetByID(id)
	if err != nil {
		return err
	}

	return EncodeResponse(w, http.StatusOK, &APIResponse{Status: http.StatusOK, Message: "account found", Data: account})
}

func (s *APIServer) handleUpdateAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	if ok := auth.AccountIDAuth(r, id, types.RoleAdmin); !ok {
		return fmt.Errorf("permission denied")
	}

	account, err := s.db.Account.GetByID(id)
	if err != nil {
		return err
	}

	req := UpdateAccountRequest{}

	err = DecodeRequest(r, &req)
	if err != nil {
		return err
	}

	if req.Username != "" {
		account.Username = req.Username
	}

	if req.Password != "" {
		account.Password = req.Password
	}

	if req.Role != "" {
		account.Role = req.Role
	}

	account, err = s.db.Account.Update(account)
	if err != nil {
		return err
	}

	return EncodeResponse(w, http.StatusOK, &APIResponse{Status: http.StatusOK, Message: "account updated", Data: account})
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	if ok := auth.AccountIDAuth(r, id, types.RoleAdmin); !ok {
		return fmt.Errorf("permission denied")
	}

	err = s.db.Account.Delete(id)
	if err != nil {
		return err
	}

	return EncodeResponse(w, http.StatusNoContent, &APIResponse{Status: http.StatusNoContent, Message: "account deleted"})
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	req := LoginRequest{}
	err := DecodeRequest(r, &req)
	if err != nil {
		return err
	}

	account, err := s.db.Account.GetByUsername(req.Username)
	if err != nil {
		return err
	}

	if !account.CheckPasswordHash(req.Password) {
		return fmt.Errorf("invalid password")
	}

	token, err := auth.GenerateJWT(account)
	if err != nil {
		return err
	}

	return EncodeResponse(w, http.StatusOK, &APIResponse{Status: http.StatusOK, Message: "login successful", Data: &LoginResponse{Token: token}})
}

func CreateAccount(username string, password string, role types.Role) (*types.Account, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &types.Account{
		Username:  username,
		Password:  string(bytes),
		Role:      role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

type CreateAccountRequest struct {
	Username string     `json:"username"`
	Password string     `json:"password"`
	Role     types.Role `json:"role"`
}

type UpdateAccountRequest struct {
	Username string     `json:"username"`
	Password string     `json:"password"`
	Role     types.Role `json:"role"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
