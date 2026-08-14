// Package interfaces defines the hexagonal-architecture ports for the users
// domain: UserRepository, UserService, and UserHandler.
//
// revive:disable:var-naming The package name "interfaces" is meaningful in this context.
package interfaces

import (
	"context"

	"backend_crudgo/users/models"

	"github.com/labstack/echo/v4"
)

// UserRepository defines data access operations for users.
// Implementations must map customErrors.ErrNotFound / customErrors.ErrConflict /
// customErrors.ErrUnauthorized for their respective failure conditions.
type UserRepository interface {
	// CreateUser persists a new user, hashing the password before storage.
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	// GetUser retrieves a single user by ID.
	GetUser(ctx context.Context, id int64) (*models.User, error)
	// GetUsers retrieves every user.
	GetUsers(ctx context.Context) ([]*models.User, error)
	// LoginUser verifies credentials and returns a signed JWT on success.
	LoginUser(ctx context.Context, userName, password string) (*models.LoginResponse, error)
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
	//   - customErrors.ErrConflict if the username or email already exist
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)

	// LoginUser verifies credentials and returns a signed JWT on success.
	//
	// Parameters:
	//   - userName: the username to authenticate
	//   - password: the plaintext password to verify
	//
	// Returns:
	//   - a LoginResponse carrying the signed JWT
	//   - customErrors.ErrUnauthorized if the credentials are invalid
	LoginUser(ctx context.Context, userName, password string) (*models.LoginResponse, error)

	// GetUser retrieves a single user by ID.
	//
	// Parameters:
	//   - id: the user identifier
	//
	// Returns:
	//   - the matching user (without password)
	//   - customErrors.ErrNotFound if no user matches id
	GetUser(ctx context.Context, id int64) (*models.User, error)

	// GetUsers retrieves every user.
	//
	// Returns:
	//   - the list of users (without passwords), empty if none exist
	GetUsers(ctx context.Context) ([]*models.User, error)
}

// UserHandler exposes the HTTP endpoints for users.
type UserHandler interface {
	// CreateUser handles POST /users/register.
	CreateUser(c echo.Context) error
	// LoginUser handles POST /users/login.
	LoginUser(c echo.Context) error
	// GetUsers handles GET /users.
	GetUsers(c echo.Context) error
}
