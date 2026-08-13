package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend_crudgo/domain/users/domain/model"
	"backend_crudgo/infrastructure/kit/apperrors"
	"backend_crudgo/mocks"

	"github.com/labstack/echo/v4"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewUserHandler(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	h := NewUserHandler(pool)

	assert.NotNil(t, h)
	assert.NotNil(t, h.Service)
}

func newUserRequest(t *testing.T, method, target string, body any) (echo.Context, *httptest.ResponseRecorder) {
	t.Helper()

	b, err := json.Marshal(body)
	assert.NoError(t, err)

	req := httptest.NewRequest(method, target, bytes.NewReader(b))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	return echo.New().NewContext(req, rec), rec
}

func TestCreateUserHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := mocks.NewUserService(t)
		h := &UserHandler{Service: mockSvc}

		user := model.User{Name: "Test User", UserPassword: "secret"}
		created := &model.User{UserID: 1, Name: "Test User"}

		mockSvc.On("CreateUser", mock.Anything, &user).Return(created, nil)

		c, rec := newUserRequest(t, http.MethodPost, "/users/register", user)

		assert.NoError(t, h.CreateUser(c))
		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("invalid body", func(t *testing.T) {
		mockSvc := mocks.NewUserService(t)
		h := &UserHandler{Service: mockSvc}

		req := httptest.NewRequest(http.MethodPost, "/users/register", bytes.NewReader([]byte("{invalid")))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		assert.NoError(t, h.CreateUser(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("conflict", func(t *testing.T) {
		mockSvc := mocks.NewUserService(t)
		h := &UserHandler{Service: mockSvc}

		user := model.User{Name: "Test User"}
		mockSvc.On("CreateUser", mock.Anything, &user).Return(nil, apperrors.ErrConflict)

		c, rec := newUserRequest(t, http.MethodPost, "/users/register", user)

		assert.NoError(t, h.CreateUser(c))
		assert.Equal(t, http.StatusConflict, rec.Code)
	})
}

func TestLoginUserHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := mocks.NewUserService(t)
		h := &UserHandler{Service: mockSvc}

		loginResponse := &model.LoginResponse{Token: "signed-jwt"}
		mockSvc.On("LoginUser", mock.Anything, "testuser", "secret").Return(loginResponse, nil)

		c, rec := newUserRequest(t, http.MethodPost, "/users/login", LoginRequest{UserName: "testuser", Password: "secret"})

		assert.NoError(t, h.LoginUser(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("invalid body", func(t *testing.T) {
		mockSvc := mocks.NewUserService(t)
		h := &UserHandler{Service: mockSvc}

		req := httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewReader([]byte("{invalid")))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		assert.NoError(t, h.LoginUser(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		mockSvc := mocks.NewUserService(t)
		h := &UserHandler{Service: mockSvc}

		mockSvc.On("LoginUser", mock.Anything, "testuser", "wrong").Return(nil, apperrors.ErrUnauthorized)

		c, rec := newUserRequest(t, http.MethodPost, "/users/login", LoginRequest{UserName: "testuser", Password: "wrong"})

		assert.NoError(t, h.LoginUser(c))
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}

func TestGetUsersHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := mocks.NewUserService(t)
		h := &UserHandler{Service: mockSvc}

		users := []*model.User{{UserID: 1}, {UserID: 2}}
		mockSvc.On("GetUsers", mock.Anything).Return(users, nil)

		c, rec := newUserRequest(t, http.MethodGet, "/users", nil)

		assert.NoError(t, h.GetUsers(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("service error", func(t *testing.T) {
		mockSvc := mocks.NewUserService(t)
		h := &UserHandler{Service: mockSvc}

		mockSvc.On("GetUsers", mock.Anything).Return(nil, assert.AnError)

		c, rec := newUserRequest(t, http.MethodGet, "/users", nil)

		assert.NoError(t, h.GetUsers(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
