// Package repository implements the user repository against PostgreSQL.
package repository

import (
	"context"
	"errors"
	"time"

	"backend_crudgo/pkg/kit/customErrors"
	"backend_crudgo/pkg/kit/enums"
	"backend_crudgo/users/constants"
	"backend_crudgo/users/interfaces"
	"backend_crudgo/users/mapper"
	"backend_crudgo/users/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	pgxtool "github.com/jnates/go-toolkit/tools/sqlconnection/pgx"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jnates/go-toolkit/tools/env"
	"github.com/jnates/go-toolkit/tools/querybuilder"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

const (
	uniqueViolationCode = "23505"
	tokenTTL            = 24 * time.Hour
)

type userRepository struct {
	pool pgxtool.DBPool
}

// NewUserRepository builds a UserRepository backed by the given connection pool.
func NewUserRepository(pool pgxtool.DBPool) interfaces.UserRepository {
	return &userRepository{pool: pool}
}

// CreateUser hashes the password and persists a new user.
//
// Parameters:
//   - user: the user to create; UserPassword must contain the plaintext password
//
// Returns:
//   - the created user (without password)
//   - customErrors.ErrConflict if user_name or user_email already exist
func (r *userRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	subLogger := log.With().Str("repository", "UserRepository").Str("method", "CreateUser").Logger()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.UserPassword), bcrypt.DefaultCost)
	if err != nil {
		subLogger.Error().Err(err).Msg("error hashing password")
		return nil, err
	}

	query, args, errBuilder := querybuilder.NewInsertBuilder().
		Into(constants.UsersTable).
		Columns(constants.ColumnUserName, constants.ColumnUserIdentifier, constants.ColumnUserEmail,
			constants.ColumnUserPassword, constants.ColumnUserTypeIdentifier).
		Values(user.Name, user.UserIdentifier, user.Email, string(hashedPassword), user.UserTypeIdentifier).
		Returning(constants.ColumnUserID).
		Build()
	if errBuilder != nil {
		subLogger.Error().Err(errBuilder).Msg("error building insert query")
		return nil, errBuilder
	}

	var userID int64
	if err := r.pool.QueryRow(ctx, query, args...).Scan(&userID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode {
			return nil, customErrors.ErrConflict
		}
		subLogger.Error().Err(err).Msg("error executing insert query")
		return nil, err
	}

	created := *user
	created.UserID = userID
	created.UserPassword = ""

	return &created, nil
}

// GetUser retrieves a single user by ID.
//
// Parameters:
//   - id: the user identifier
//
// Returns:
//   - the matching user (without password)
//   - customErrors.ErrNotFound if no user matches id
func (r *userRepository) GetUser(ctx context.Context, id int64) (*models.User, error) {
	subLogger := log.With().Str("repository", "UserRepository").Str("method", "GetUser").Logger()

	query, args, errBuilder := querybuilder.NewSelectBuilder().
		Select(constants.Columns...).
		From(constants.UsersTable).
		Where(constants.ColumnUserID, querybuilder.OpEqual, id).
		Limit(1).
		Build()
	if errBuilder != nil {
		subLogger.Error().Err(errBuilder).Msg("error building select query")
		return nil, errBuilder
	}

	user, err := scanUserDB(r.pool.QueryRow(ctx, query, args...))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, customErrors.ErrNotFound
		}
		subLogger.Error().Err(err).Msg("error executing select query")
		return nil, err
	}

	return mapper.ToUserModel(&user), nil
}

// GetUsers retrieves every user.
//
// Returns:
//   - the list of users (without passwords), empty if none exist
func (r *userRepository) GetUsers(ctx context.Context) ([]*models.User, error) {
	subLogger := log.With().Str("repository", "UserRepository").Str("method", "GetUsers").Logger()

	query, args, errBuilder := querybuilder.NewSelectBuilder().
		Select(constants.Columns...).
		From(constants.UsersTable).
		OrderBy(constants.ColumnUserID, querybuilder.Asc).
		Build()
	if errBuilder != nil {
		subLogger.Error().Err(errBuilder).Msg("error building select query")
		return nil, errBuilder
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		subLogger.Error().Err(err).Msg("error executing select query")
		return nil, err
	}
	defer rows.Close()

	items, err := pgxtool.ScanRows(rows, subLogger, scanUserDB)
	if err != nil {
		return nil, err
	}

	users := make([]*models.User, 0, len(items))
	for i := range items {
		users = append(users, mapper.ToUserModel(&items[i]))
	}

	return users, nil
}

// scanUserDB scans a single row into a UserDB, following column order in constants.Columns.
func scanUserDB(s pgxtool.Scanner) (models.UserDB, error) {
	var u models.UserDB
	err := s.Scan(&u.UserID, &u.UserName, &u.UserIdentifier, &u.UserEmail, &u.UserPassword, &u.UserTypeIdentifier)
	return u, err
}

// LoginUser verifies the given credentials and returns a signed JWT on success.
// It never reveals whether the username or the password was the cause of a failure.
//
// Parameters:
//   - userName: the username to authenticate
//   - password: the plaintext password to verify
//
// Returns:
//   - a LoginResponse carrying the signed JWT
//   - customErrors.ErrUnauthorized if the credentials are invalid
func (r *userRepository) LoginUser(ctx context.Context, userName, password string) (*models.LoginResponse, error) {
	subLogger := log.With().Str("repository", "UserRepository").Str("method", "LoginUser").Logger()

	query, args, errBuilder := querybuilder.NewSelectBuilder().
		Select(constants.Columns...).
		From(constants.UsersTable).
		Where(constants.ColumnUserName, querybuilder.OpEqual, userName).
		Limit(1).
		Build()
	if errBuilder != nil {
		subLogger.Error().Err(errBuilder).Msg("error building select query")
		return nil, errBuilder
	}

	user, err := scanUserDB(r.pool.QueryRow(ctx, query, args...))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, customErrors.ErrUnauthorized
		}
		subLogger.Error().Err(err).Msg("error executing select query")
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.UserPassword), []byte(password)); err != nil {
		return nil, customErrors.ErrUnauthorized
	}

	token, err := generateToken(user.UserID)
	if err != nil {
		subLogger.Error().Err(err).Msg("error generating token")
		return nil, err
	}

	return &models.LoginResponse{Token: token}, nil
}

// generateToken signs a JWT for the given user ID using the configured secret key.
func generateToken(userID int64) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub": userID,
		"iat": now.Unix(),
		"exp": now.Add(tokenTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(env.GetString(enums.SecretKey, "")))
}
