package middlewares

import (
	"backend_crudgo/infrastructure/kit/enum"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
)

// AuthMiddleware is a middleware function that validates a JSON Web Token (JWT) in an HTTP request.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == enum.EmptyString {
			http.Error(w, "Authorization token missing", http.StatusUnauthorized)
			return
		}
		tokenString = tokenString[len("Bearer "):]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			secretKey := os.Getenv(enum.SecretKey)
			return []byte(secretKey), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
