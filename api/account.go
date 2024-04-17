package api

import (
	"net/http"
	"ticketing-api/auth"
	"ticketing-api/types"
)

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	req := CreateAccountRequest{}

	err := decodeRequest(r, &req)
	if err != nil {
		return err
	}

	account, err := types.CreateAccount(req.Username, req.Password, req.Role)
	if err != nil {
		return err
	}

	account, err = s.db.Account.Create(account)
	if err != nil {
		return err
	}

	return encodeResponse(w, http.StatusOK, &APIResponse{Status: http.StatusOK, Message: "account created", Data: account})
}

func (s *APIServer) handleGetAccounts(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.db.Account.Get()
	if err != nil {
		return err
	}

	return encodeResponse(w, http.StatusOK, &APIResponse{Status: http.StatusOK, Message: "accounts found", Data: accounts})
}

func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	err = auth.IsAccountID(r, id, types.RoleAdmin, types.RoleEditor)
	if err != nil {
		return err
	}

	account, err := s.db.Account.GetByID(id)
	if err != nil {
		return err
	}

	return encodeResponse(w, http.StatusOK, &APIResponse{Status: http.StatusOK, Message: "account found", Data: account})
}

func (s *APIServer) handleUpdateAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	err = auth.IsAccountID(r, id, types.RoleAdmin)
	if err != nil {
		return err
	}

	account, err := s.db.Account.GetByID(id)
	if err != nil {
		return err
	}

	req := UpdateAccountRequest{}

	err = decodeRequest(r, &req)
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

	return encodeResponse(w, http.StatusOK, &APIResponse{Status: http.StatusOK, Message: "account updated", Data: account})
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	err = auth.IsAccountID(r, id, types.RoleAdmin)
	if err != nil {
		return err
	}

	err = s.db.Account.Delete(id)
	if err != nil {
		return err
	}

	return encodeResponse(w, http.StatusNoContent, &APIResponse{Status: http.StatusNoContent, Message: "account deleted"})
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	req := LoginRequest{}
	err := decodeRequest(r, &req)
	if err != nil {
		return err
	}

	account, err := s.db.Account.GetByUsername(req.Username)
	if err != nil {
		return err
	}

	err = account.CheckPasswordHash(req.Password)
	if err != nil {
		return err
	}

	token, err := auth.GenerateJWT(account)
	if err != nil {
		return err
	}

	return encodeResponse(w, http.StatusOK, &APIResponse{Status: http.StatusOK, Message: "login successful", Data: &LoginResponse{Token: token}})
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
