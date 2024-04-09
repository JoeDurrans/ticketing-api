package auth

import (
	"fmt"
	"net/http"
	"os"
	"ticketing-api/types"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func AccountIDAuth(r *http.Request, id int, roles ...types.Role) bool {
	claims, err := getClaims(r)
	if err != nil {
		return false
	}

	accountID := int(claims["id"].(float64))

	if accountID == id {
		return true
	}

	for _, role := range roles {
		if role == types.Role(claims["role"].(string)) {
			return true
		}
	}

	return false
}

func IsAdmin(r *http.Request) bool {
	claims, err := getClaims(r)
	if err != nil {
		return false
	}

	return types.Role(claims["role"].(string)) == types.RoleAdmin
}

func IsEditor(r *http.Request) bool {
	claims, err := getClaims(r)
	if err != nil {
		return false
	}

	return types.Role(claims["role"].(string)) == types.RoleEditor || types.Role(claims["role"].(string)) == types.RoleAdmin
}

func IsAuthenticated(r *http.Request) bool {
	_, err := getClaims(r)
	return err == nil
}

func GetID(r *http.Request) (int, error) {
	claims, err := getClaims(r)
	if err != nil {
		return 0, err
	}

	return int(claims["id"].(float64)), nil
}

func GetRole(r *http.Request) (types.Role, error) {
	claims, err := getClaims(r)
	if err != nil {
		return "", err
	}

	return types.Role(claims["role"].(string)), nil
}

func getClaims(r *http.Request) (jwt.MapClaims, error) {
	tokenStr := r.Header.Get("Authorization")

	token, err := ValidateJWT(tokenStr)
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token.Claims.(jwt.MapClaims), nil
}

func GenerateJWT(a *types.Account) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
		"id":   a.ID,
		"role": a.Role,
	}).SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func ValidateJWT(t string) (*jwt.Token, error) {
	return jwt.Parse(t, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})
}
