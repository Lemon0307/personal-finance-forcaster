package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var key = []byte("pfftesting")

func (account *Account) GenerateJWT() (string, error) {
	// set JWT expiration date
	expiration_time := time.Now().Add(720 * time.Hour)

	// set up information that is stored in the JWT
	claims := &Claims{
		UserID: account.UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration_time),
		},
	}
	// encode claims and sign with secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token_string, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	fmt.Println("Generated JWT:", token_string)
	return token_string, nil
}

func ValidateJWT(token_string string) (*Claims, error) {
	claims := &Claims{}
	// decode jwt to get user_id
	token, err := jwt.ParseWithClaims(token_string, claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("incorrect signing method")
			}
			return key, nil
		})

	if err != nil {
		fmt.Println("Error parsing token:", err)
		// check if token has expired
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

const UserIDkey ctxkey = "user_id"

// runs this code every request
func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// find token from authorization header
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			fmt.Println("Missing token")
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		// remove Bearer from the authorization string to get token only
		// also decodes the jwt to get user_id
		claims, err := ValidateJWT(strings.TrimPrefix(tokenString, "Bearer "))
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// store user_id in context
		ctx := context.WithValue(r.Context(), UserIDkey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
