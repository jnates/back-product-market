package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend_crudgo/pkg/kit/customErrors"
	"backend_crudgo/users/mocks"
	"backend_crudgo/users/models"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewUserHandler(t *testing.T) {
	mockSvc := mocks.NewUserService(t)

	h := NewUserHandler(mockSvc)

	assert.NotNil(t, h)
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
		h := NewUserHandler(mockSvc)

		user := models.User{Name: "Test User", UserPassword: "secret"}
		created := &models.User{UserID: 1, Name: "Test User"}

		mockSvc.On("CreateUser", mock.Anything, &user).Return(created, nil)

		c, rec := newUserRequest(t, http.MethodPost, "/users/register", user)

		assert.NoError(t, h.CreateUser(c))
		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("invalid body", func(t *testing.T) {
		mockSvc := mocks.NewUserService(t)
		h := NewUserHandler(mockSvc)

		req := httptest.NewRequest(http.MethodPost, "/users/register", bytes.NewReader([]byte("{invalid")))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		assert.NoError(t, h.CreateUser(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("conflict", func(t *testing.T) {
		mockSvc := mocks.NewUserService(t)
		h := NewUserHandler(mockSvc)

		user := models.User{Name: "Test User"}
		mockSvc.On("CreateUser", mock.Anything, &user).Return(nil, customErrors.ErrConflict)

		c, rec := newUserRequest(t, http.MethodPost, "/users/register", user)

		assert.NoError(t, h.CreateUser(c))
		assert.Equal(t, http.StatusConflict, rec.Code)
	})
}

func TestLoginUserHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := mocks.NewUserService(t)
		h := NewUserHandler(mockSvc)

		loginResponse := &models.LoginResponse{Token: "signed-jwt"}
		mockSvc.On("LoginUser", mock.Anything, "testuser", "secret").Return(loginResponse, nil)

		c, rec := newUserRequest(t, http.MethodPost, "/users/login", LoginRequest{UserName: "testuser", Password: "secret"})

		assert.NoError(t, h.LoginUser(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("invalid body", func(t *testing.T) {
		mockSvc := mocks.NewUserService(t)
		h := NewUserHandler(mockSvc)

		req := httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewReader([]byte("{invalid")))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		assert.NoError(t, h.LoginUser(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		mockSvc := mocks.NewUserService(t)
		h := NewUserHandler(mockSvc)

		mockSvc.On("LoginUser", mock.Anything, "testuser", "wrong").Return(nil, customErrors.ErrUnauthorized)

		c, rec := newUserRequest(t, http.MethodPost, "/users/login", LoginRequest{UserName: "testuser", Password: "wrong"})

		assert.NoError(t, h.LoginUser(c))
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}

func TestGetUsersHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := mocks.NewUserService(t)
		h := NewUserHandler(mockSvc)

		users := []*models.User{{UserID: 1}, {UserID: 2}}
		mockSvc.On("GetUsers", mock.Anything).Return(users, nil)

		c, rec := newUserRequest(t, http.MethodGet, "/users", nil)

		assert.NoError(t, h.GetUsers(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("service error", func(t *testing.T) {
		mockSvc := mocks.NewUserService(t)
		h := NewUserHandler(mockSvc)

		mockSvc.On("GetUsers", mock.Anything).Return(nil, assert.AnError)

		c, rec := newUserRequest(t, http.MethodGet, "/users", nil)

		assert.NoError(t, h.GetUsers(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
