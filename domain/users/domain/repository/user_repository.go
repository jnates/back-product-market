// Package repository defines the persistence contract for users.
package repository

import (
	"context"

	"backend_crudgo/domain/users/domain/model"
)

// UserRepository defines data access operations for users.
// Implementations must map apperrors.ErrNotFound / apperrors.ErrConflict / apperrors.ErrUnauthorized
// for their respective failure conditions.
type UserRepository interface {
	// CreateUser persists a new user, hashing the password before storage.
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	// GetUser retrieves a single user by ID.
	GetUser(ctx context.Context, id int64) (*model.User, error)
	// GetUsers retrieves every user.
	GetUsers(ctx context.Context) ([]*model.User, error)
	// LoginUser verifies credentials and returns a signed JWT on success.
	LoginUser(ctx context.Context, userName, password string) (*model.LoginResponse, error)
}
