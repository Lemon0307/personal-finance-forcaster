package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// runs this code every request

// change this bit
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
