package test

import (
	"net/http/httptest"
	"testing"
	"ticketing-api/auth"
	"ticketing-api/types"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateJWT(t *testing.T) {
	account := &types.Account{ID: 1, Role: "admin"}
	tokenString, err := auth.GenerateJWT(account)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	token, err := auth.ValidateJWT(tokenString)
	if err != nil || !token.Valid {
		t.Fatalf("expected valid token, got: %v", err)
	}

	claims := token.Claims.(jwt.MapClaims)
	if claims["id"].(float64) != float64(account.ID) || claims["role"].(string) != string(account.Role) {
		t.Fatalf("token claims do not match expected values")
	}
}

func TestValidateJWTInvalidToken(t *testing.T) {
	_, err := auth.ValidateJWT("invalid")
	if err == nil {
		t.Fatalf("expected error for invalid token, got none")
	}
}

func TestGetRole(t *testing.T) {
	account := &types.Account{ID: 1, Role: "admin"}
	tokenString, err := auth.GenerateJWT(account)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Add("Authorization", "Bearer "+tokenString)

	role, err := auth.GetRole(r)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if role != account.Role {
		t.Fatalf("expected role to be %s, got %s", account.Role, role)
	}
}

func TestGetRoleInvalidToken(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Add("Authorization", "invalid")

	_, err := auth.GetRole(r)
	if err == nil {
		t.Fatalf("expected error for invalid token, got none")
	}
}

func TestGetAccountID(t *testing.T) {
	account := &types.Account{ID: 1, Role: "admin"}
	tokenString, err := auth.GenerateJWT(account)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Add("Authorization", "Bearer "+tokenString)

	id, err := auth.GetAccountID(r)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if id != account.ID {
		t.Fatalf("expected ID to be %d, got %d", account.ID, id)
	}
}

func TestGetAccountIDInvalidToken(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Add("Authorization", "invalid")

	_, err := auth.GetAccountID(r)
	if err == nil {
		t.Fatalf("expected error for invalid token, got none")
	}
}

func TestIsAuthenticated(t *testing.T) {
	account := &types.Account{ID: 1, Role: "admin"}
	tokenString, err := auth.GenerateJWT(account)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Add("Authorization", "Bearer "+tokenString)

	err = auth.IsAuthenticated(r)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestIsAuthenticatedInvalidToken(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Add("Authorization", "invalid")

	err := auth.IsAuthenticated(r)
	if err == nil {
		t.Fatalf("expected error for invalid token, got none")
	}
}

func TestIsRole(t *testing.T) {
	account := &types.Account{ID: 1, Role: "admin"}
	tokenString, err := auth.GenerateJWT(account)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Add("Authorization", "Bearer "+tokenString)

	err = auth.IsRole(r, types.RoleAdmin)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestIsRoleInvalidRole(t *testing.T) {
	account := &types.Account{ID: 1, Role: "editor"}
	tokenString, err := auth.GenerateJWT(account)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Add("Authorization", "Bearer "+tokenString)

	err = auth.IsRole(r, types.RoleAdmin)
	if err == nil {
		t.Fatalf("expected error for invalid role, got none")
	}
}

func TestIsAccountID(t *testing.T) {
	account := &types.Account{ID: 1, Role: "admin"}
	tokenString, err := auth.GenerateJWT(account)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Add("Authorization", "Bearer "+tokenString)

	err = auth.IsAccountID(r, 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestIsAccountIDInvalidID(t *testing.T) {
	account := &types.Account{ID: 1, Role: "admin"}
	tokenString, err := auth.GenerateJWT(account)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Add("Authorization", "Bearer "+tokenString)

	err = auth.IsAccountID(r, 2)
	if err == nil {
		t.Fatalf("expected error for invalid ID, got none")
	}
}

func TestIsAccountIDInvalidIDValidRole(t *testing.T) {
	account := &types.Account{ID: 1, Role: "admin"}
	tokenString, err := auth.GenerateJWT(account)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Add("Authorization", "Bearer "+tokenString)

	err = auth.IsAccountID(r, 2, types.RoleAdmin)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestCompareHashAndPassword(t *testing.T) {
	password := "password"
	hash, err := auth.CreateHash(password)
	if err != nil {
		t.Fatalf("failed to create hash: %s", err)
	}

	err = auth.CompareHashAndPassword(hash, password)
	if err != nil {
		t.Errorf("expected valid comparison, got: %s", err)
	}
}

func TestCompareHashAndPasswordInvalid(t *testing.T) {
	password := "password"
	wrongPassword := "wrongpassword"
	hash, err := auth.CreateHash(password)
	if err != nil {
		t.Fatalf("Failed to create hash: %s", err)
	}

	err = auth.CompareHashAndPassword(hash, wrongPassword)
	if err == nil {
		t.Error("Expected error for incorrect password, got nil")
	}

	if _, ok := err.(*types.Unauthorized); !ok {
		t.Errorf("Expected error of type *types.Unauthorized, got %T", err)
	}
}
