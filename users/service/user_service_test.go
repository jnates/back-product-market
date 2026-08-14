package service

import (
	"context"
	"testing"

	"backend_crudgo/pkg/kit/customErrors"
	"backend_crudgo/users/mocks"
	"backend_crudgo/users/models"

	"github.com/stretchr/testify/assert"
)

func TestUserService_CreateUser(t *testing.T) {
	assertions := assert.New(t)
	mockRepo := mocks.NewUserRepository(t)
	userService := NewUserService(mockRepo)

	user := &models.User{
		Name:               "Test User",
		Email:              "test@example.com",
		UserIdentifier:     1,
		UserPassword:       "plaintext",
		UserTypeIdentifier: 1,
	}

	created := &models.User{
		UserID:             1,
		Name:               "Test User",
		Email:              "test@example.com",
		UserIdentifier:     1,
		UserTypeIdentifier: 1,
	}

	mockRepo.On("CreateUser", context.Background(), user).Return(created, nil)

	res, err := userService.CreateUser(context.Background(), user)
	assertions.NoError(err)
	assertions.Equal(created, res)
}

func TestUserService_CreateUser_Conflict(t *testing.T) {
	assertions := assert.New(t)
	mockRepo := mocks.NewUserRepository(t)
	userService := NewUserService(mockRepo)

	user := &models.User{Name: "Test User"}

	mockRepo.On("CreateUser", context.Background(), user).Return(nil, customErrors.ErrConflict)

	res, err := userService.CreateUser(context.Background(), user)
	assertions.ErrorIs(err, customErrors.ErrConflict)
	assertions.Nil(res)
}

func TestUserService_GetUser(t *testing.T) {
	assertions := assert.New(t)
	mockRepo := mocks.NewUserRepository(t)
	userService := NewUserService(mockRepo)

	user := &models.User{UserID: 1, Name: "Test User"}

	mockRepo.On("GetUser", context.Background(), int64(1)).Return(user, nil)

	res, err := userService.GetUser(context.Background(), 1)
	assertions.NoError(err)
	assertions.Equal(user, res)
}

func TestUserService_GetUser_NotFound(t *testing.T) {
	assertions := assert.New(t)
	mockRepo := mocks.NewUserRepository(t)
	userService := NewUserService(mockRepo)

	mockRepo.On("GetUser", context.Background(), int64(2)).Return(nil, customErrors.ErrNotFound)

	res, err := userService.GetUser(context.Background(), 2)
	assertions.ErrorIs(err, customErrors.ErrNotFound)
	assertions.Nil(res)
}

func TestUserService_GetUsers(t *testing.T) {
	assertions := assert.New(t)
	mockRepo := mocks.NewUserRepository(t)
	userService := NewUserService(mockRepo)

	users := []*models.User{
		{UserID: 1, Name: "Test User 1"},
		{UserID: 2, Name: "Test User 2"},
	}

	mockRepo.On("GetUsers", context.Background()).Return(users, nil)

	res, err := userService.GetUsers(context.Background())
	assertions.NoError(err)
	assertions.Equal(users, res)
}

func TestUserService_LoginUser(t *testing.T) {
	assertions := assert.New(t)
	mockRepo := mocks.NewUserRepository(t)
	userService := NewUserService(mockRepo)

	loginResponse := &models.LoginResponse{Token: "signed-jwt"}

	mockRepo.On("LoginUser", context.Background(), "testuser", "secret").Return(loginResponse, nil)

	res, err := userService.LoginUser(context.Background(), "testuser", "secret")
	assertions.NoError(err)
	assertions.Equal(loginResponse, res)
}

func TestUserService_LoginUser_Unauthorized(t *testing.T) {
	assertions := assert.New(t)
	mockRepo := mocks.NewUserRepository(t)
	userService := NewUserService(mockRepo)

	mockRepo.On("LoginUser", context.Background(), "testuser", "wrong").Return(nil, customErrors.ErrUnauthorized)

	res, err := userService.LoginUser(context.Background(), "testuser", "wrong")
	assertions.ErrorIs(err, customErrors.ErrUnauthorized)
	assertions.Nil(res)
}
