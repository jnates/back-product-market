package middlewares_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend_crudgo/infrastructure/kit/apperrors"
	"backend_crudgo/infrastructure/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleError(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedStatus int
	}{
		{
			name:           "not found",
			err:            apperrors.ErrNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "invalid input",
			err:            apperrors.ErrInvalidInput,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "unauthorized",
			err:            apperrors.ErrUnauthorized,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "conflict",
			err:            apperrors.ErrConflict,
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "wrapped not found",
			err:            errors.New("wrapped: " + apperrors.ErrNotFound.Error()),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "unknown error",
			err:            errors.New("boom"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := middlewares.HandleError(c, tt.err)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}
