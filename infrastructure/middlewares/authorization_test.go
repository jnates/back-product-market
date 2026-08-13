package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"backend_crudgo/infrastructure/kit/enum"
	"backend_crudgo/infrastructure/middlewares"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthMiddleware(t *testing.T) {
	secretKey := "test-secret"
	t.Setenv(enum.SecretKey, secretKey)

	tests := []struct {
		name           string
		token          string
		expectedStatus int
	}{
		{
			name:           "Should return OK when valid token is provided",
			token:          validToken(t, secretKey),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Should return Unauthorized when invalid token is provided",
			token:          "invalid-token",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Should return Unauthorized when token is missing",
			token:          "",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authMiddleware, err := middlewares.NewAuthMiddleware()
			require.NoError(t, err)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.token != enum.EmptyString {
				req.Header.Set(enum.Authorization, "Bearer "+tt.token)
			}
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			handler := authMiddleware(func(c echo.Context) error {
				return c.NoContent(http.StatusOK)
			})

			_ = handler(c)

			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}

func validToken(t *testing.T, secretKey string) string {
	t.Helper()

	claims := jwt.MapClaims{
		"sub": int64(1),
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(secretKey))
	require.NoError(t, err)

	return signed
}
