package middlewares_test

import (
	"backend_crudgo/infrastructure/kit/enum"
	"backend_crudgo/infrastructure/middlewares"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	tests := []struct {
		name      string
		token     string
		expectErr bool
	}{
		{
			name:      "Should return OK when valid token is provided",
			token:     getTokenString(enum.SecretKey),
			expectErr: false,
		},
		{
			name:      "Should return Unauthorized when invalid token is provided",
			token:     "invalid-token",
			expectErr: true,
		},
		{
			name:      "Should return Unauthorized when token is missing",
			token:     enum.EmptyString,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertions := assert.New(t)
			handlerFunc := func(w http.ResponseWriter, r *http.Request) {}
			nextHandler := http.HandlerFunc(handlerFunc)
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.token != enum.EmptyString {
				req.Header.Set(enum.Authorization, "Bearer "+tt.token)
			}
			os.Setenv(enum.SecretKey, enum.SecretKey)
			recorder := httptest.NewRecorder()

			middlewares.AuthMiddleware(nextHandler).ServeHTTP(recorder, req)

			if tt.expectErr {
				assertions.Equal(http.StatusUnauthorized, recorder.Code, "expected status code %v, but got %v", http.StatusUnauthorized, recorder.Code)
			} else {
				assertions.Equal(http.StatusOK, recorder.Code, "expected status code %v, but got %v", http.StatusOK, recorder.Code)
			}
		})
	}
}

func getTokenString(secretKey string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{})
	tokenString, _ := token.SignedString([]byte(secretKey))
	return tokenString
}
