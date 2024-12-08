package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
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
	fmt.Println("Generated JWT:", token_string)
	return token_string, nil
}

func ValidateJWT(token_string string) (*Claims, error) {
	fmt.Println(token_string)
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(token_string, claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("incorrect signing method")
			}
			return key, nil
		})

	if err != nil {
		fmt.Println("Error parsing token:", err)
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

type ctxkey string

const userIDkey ctxkey = "userID"

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		claims, err := ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDkey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
