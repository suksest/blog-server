package handlers

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JwtClaims struct {
	Hash string `json:"hash"`
	jwt.StandardClaims
}

func createJwtToken(s string) (string, error) {
	claims := JwtClaims{
		s,
		jwt.StandardClaims{
			Id:        "main_user_id",
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	token, err := rawToken.SignedString([]byte("mySecret"))
	if err != nil {
		return "", err
	}

	return token, nil
}
