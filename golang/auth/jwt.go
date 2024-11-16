package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var key = []byte("pfftesting")

func (account *Account) GenerateJWT() (string, error) {
	expiration_time := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID: account.UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration_time),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token_string, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return token_string, nil
}

func (account *Account) ValidateJWT(token_string string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(token_string, claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("incorrect signing method")
			}
			return key, nil
		})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token has expired")
		}
		return nil, fmt.Errorf("invalid token")
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
