package service

import (
	"context"
	"testing"

	"backend_crudgo/domain/users/domain/model"
	"backend_crudgo/infrastructure/kit/apperrors"
	"backend_crudgo/mocks"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	assertions := assert.New(t)
	mockRepo := mocks.NewUserRepository(t)
	userService := NewUserService(mockRepo)

	user := &model.User{
		Name:               "Test User",
		Email:              "test@example.com",
		UserIdentifier:     1,
		UserPassword:       "plaintext",
		UserTypeIdentifier: 1,
	}

	created := &model.User{
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

func TestCreateUserConflict(t *testing.T) {
	assertions := assert.New(t)
	mockRepo := mocks.NewUserRepository(t)
	userService := NewUserService(mockRepo)

	user := &model.User{Name: "Test User"}

	mockRepo.On("CreateUser", context.Background(), user).Return(nil, apperrors.ErrConflict)

	res, err := userService.CreateUser(context.Background(), user)
	assertions.ErrorIs(err, apperrors.ErrConflict)
	assertions.Nil(res)
}

func TestGetUser(t *testing.T) {
	assertions := assert.New(t)
	mockRepo := mocks.NewUserRepository(t)
	userService := NewUserService(mockRepo)

	user := &model.User{UserID: 1, Name: "Test User"}

	mockRepo.On("GetUser", context.Background(), int64(1)).Return(user, nil)

	res, err := userService.GetUser(context.Background(), 1)
	assertions.NoError(err)
	assertions.Equal(user, res)
}

func TestGetUserNotFound(t *testing.T) {
	assertions := assert.New(t)
	mockRepo := mocks.NewUserRepository(t)
	userService := NewUserService(mockRepo)

	mockRepo.On("GetUser", context.Background(), int64(2)).Return(nil, apperrors.ErrNotFound)

	res, err := userService.GetUser(context.Background(), 2)
	assertions.ErrorIs(err, apperrors.ErrNotFound)
	assertions.Nil(res)
}

func TestGetUsers(t *testing.T) {
	assertions := assert.New(t)
	mockRepo := mocks.NewUserRepository(t)
	userService := NewUserService(mockRepo)

	users := []*model.User{
		{UserID: 1, Name: "Test User 1"},
		{UserID: 2, Name: "Test User 2"},
	}

	mockRepo.On("GetUsers", context.Background()).Return(users, nil)

	res, err := userService.GetUsers(context.Background())
	assertions.NoError(err)
	assertions.Equal(users, res)
}

func TestLoginUser(t *testing.T) {
	assertions := assert.New(t)
	mockRepo := mocks.NewUserRepository(t)
	userService := NewUserService(mockRepo)

	loginResponse := &model.LoginResponse{Token: "signed-jwt"}

	mockRepo.On("LoginUser", context.Background(), "testuser", "secret").Return(loginResponse, nil)

	res, err := userService.LoginUser(context.Background(), "testuser", "secret")
	assertions.NoError(err)
	assertions.Equal(loginResponse, res)
}

func TestLoginUserUnauthorized(t *testing.T) {
	assertions := assert.New(t)
	mockRepo := mocks.NewUserRepository(t)
	userService := NewUserService(mockRepo)

	mockRepo.On("LoginUser", context.Background(), "testuser", "wrong").Return(nil, apperrors.ErrUnauthorized)

	res, err := userService.LoginUser(context.Background(), "testuser", "wrong")
	assertions.ErrorIs(err, apperrors.ErrUnauthorized)
	assertions.Nil(res)
}
