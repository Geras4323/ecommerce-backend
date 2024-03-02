package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type JwtLoginClaims struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

type JwtChangePasswordClaims struct {
	ID uint `json:"id"`
	jwt.RegisteredClaims
}

func SignToken(claims jwt.Claims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func VerifyToken(token string, claims jwt.Claims, key string) (*jwt.Token, error) {
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {
		return nil, err
	}

	return parsedToken, nil
}
