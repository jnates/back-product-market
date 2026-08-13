// Package service contains the user domain's business logic.
package service

import (
	"context"

	"backend_crudgo/domain/users/domain/model"
	"backend_crudgo/domain/users/domain/repository"
)

type userService struct {
	UserRepository repository.UserRepository
}

// UserService exposes the user use cases consumed by the HTTP handlers.
type UserService interface {
	// CreateUser creates a new user.
	//
	// Parameters:
	//   - user: the user to create, with UserPassword set to the plaintext password
	//
	// Returns:
	//   - the created user (without password)
	//   - apperrors.ErrConflict if the username or email already exist
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)

	// LoginUser verifies credentials and returns a signed JWT on success.
	//
	// Parameters:
	//   - userName: the username to authenticate
	//   - password: the plaintext password to verify
	//
	// Returns:
	//   - a LoginResponse carrying the signed JWT
	//   - apperrors.ErrUnauthorized if the credentials are invalid
	LoginUser(ctx context.Context, userName, password string) (*model.LoginResponse, error)

	// GetUser retrieves a single user by ID.
	//
	// Parameters:
	//   - id: the user identifier
	//
	// Returns:
	//   - the matching user (without password)
	//   - apperrors.ErrNotFound if no user matches id
	GetUser(ctx context.Context, id int64) (*model.User, error)

	// GetUsers retrieves every user.
	//
	// Returns:
	//   - the list of users (without passwords), empty if none exist
	GetUsers(ctx context.Context) ([]*model.User, error)
}

// NewUserService builds a UserService backed by the given repository.
func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{
		UserRepository: userRepository,
	}
}

// CreateUser creates a new user by delegating to the repository.
//
// Parameters:
//   - user: the user to create, with UserPassword set to the plaintext password
//
// Returns:
//   - the created user (without password)
//   - apperrors.ErrConflict if the username or email already exist
func (us *userService) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	return us.UserRepository.CreateUser(ctx, user)
}

// GetUser retrieves a single user by ID by delegating to the repository.
//
// Parameters:
//   - id: the user identifier
//
// Returns:
//   - the matching user (without password)
//   - apperrors.ErrNotFound if no user matches id
func (us *userService) GetUser(ctx context.Context, id int64) (*model.User, error) {
	return us.UserRepository.GetUser(ctx, id)
}

// LoginUser verifies credentials and returns a signed JWT by delegating to the repository.
//
// Parameters:
//   - userName: the username to authenticate
//   - password: the plaintext password to verify
//
// Returns:
//   - a LoginResponse carrying the signed JWT
//   - apperrors.ErrUnauthorized if the credentials are invalid
func (us *userService) LoginUser(ctx context.Context, userName, password string) (*model.LoginResponse, error) {
	return us.UserRepository.LoginUser(ctx, userName, password)
}

// GetUsers retrieves every user by delegating to the repository.
//
// Returns:
//   - the list of users (without passwords), empty if none exist
func (us *userService) GetUsers(ctx context.Context) ([]*model.User, error) {
	return us.UserRepository.GetUsers(ctx)
}
