// Package service contains the user domain's business logic.
package service

import (
	"context"

	"backend_crudgo/users/interfaces"
	"backend_crudgo/users/models"
)

// userService implements interfaces.UserService, delegating to the user repository.
type userService struct {
	userRepository interfaces.UserRepository
}

// NewUserService builds a UserService backed by the given repository.
func NewUserService(userRepository interfaces.UserRepository) interfaces.UserService {
	return &userService{
		userRepository: userRepository,
	}
}

// CreateUser creates a new user by delegating to the repository.
//
// Parameters:
//   - user: the user to create, with UserPassword set to the plaintext password
//
// Returns:
//   - the created user (without password)
//   - customErrors.ErrConflict if the username or email already exist
func (s *userService) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	return s.userRepository.CreateUser(ctx, user)
}

// GetUser retrieves a single user by ID by delegating to the repository.
//
// Parameters:
//   - id: the user identifier
//
// Returns:
//   - the matching user (without password)
//   - customErrors.ErrNotFound if no user matches id
func (s *userService) GetUser(ctx context.Context, id int64) (*models.User, error) {
	return s.userRepository.GetUser(ctx, id)
}

// LoginUser verifies credentials and returns a signed JWT by delegating to the repository.
//
// Parameters:
//   - userName: the username to authenticate
//   - password: the plaintext password to verify
//
// Returns:
//   - a LoginResponse carrying the signed JWT
//   - customErrors.ErrUnauthorized if the credentials are invalid
func (s *userService) LoginUser(ctx context.Context, userName, password string) (*models.LoginResponse, error) {
	return s.userRepository.LoginUser(ctx, userName, password)
}

// GetUsers retrieves every user by delegating to the repository.
//
// Returns:
//   - the list of users (without passwords), empty if none exist
func (s *userService) GetUsers(ctx context.Context) ([]*models.User, error) {
	return s.userRepository.GetUsers(ctx)
}
