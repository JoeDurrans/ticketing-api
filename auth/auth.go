package auth

import (
	"fmt"
	"net/http"
	"os"
	"ticketing-api/types"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func AccountIDAuth(r *http.Request, id int, roles ...types.Role) error {
	claims, err := getClaims(r)
	if err != nil {
		return err
	}

	accountID := int(claims["id"].(float64))

	if accountID == id {
		return nil
	}

	for _, role := range roles {
		if role == types.Role(claims["role"].(string)) {
			return nil
		}
	}

	return &types.Forbidden{}
}

func IsAdmin(r *http.Request) error {
	claims, err := getClaims(r)
	if err != nil {
		return err
	}

	if types.Role(claims["role"].(string)) == types.RoleAdmin {
		return nil
	}

	return &types.Forbidden{}
}

func IsEditor(r *http.Request) error {
	claims, err := getClaims(r)
	if err != nil {
		return err
	}

	if types.Role(claims["role"].(string)) == types.RoleEditor {
		return nil
	}

	return &types.Forbidden{}
}

func IsAuthenticated(r *http.Request) error {
	_, err := getClaims(r)
	if err != nil {
		return err
	}

	return nil
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
		return nil, &types.Unauthorized{
			Message: "invalid token",
		}
	}

	return token.Claims.(jwt.MapClaims), nil
}

func GenerateJWT(a *types.Account) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
		"id":   a.ID,
		"role": a.Role,
	})

	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", fmt.Errorf("failed to generate token")
	}

	return signedToken, nil
}

func ValidateJWT(t string) (*jwt.Token, error) {
	return jwt.Parse(t, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})
}
