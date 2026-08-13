// Package persistence implements the user repository against PostgreSQL.
package persistence

import (
	"context"
	"errors"
	"time"

	"backend_crudgo/domain/users/constants"
	"backend_crudgo/domain/users/domain/model"
	repoDomain "backend_crudgo/domain/users/domain/repository"
	"backend_crudgo/infrastructure/kit/apperrors"
	"backend_crudgo/infrastructure/kit/enum"

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

type sqlUserRepo struct {
	pool pgxtool.DBPool
}

// NewUserRepository builds a UserRepository backed by the given connection pool.
func NewUserRepository(pool pgxtool.DBPool) repoDomain.UserRepository {
	return &sqlUserRepo{pool: pool}
}

// CreateUser hashes the password and persists a new user.
//
// Parameters:
//   - user: the user to create; UserPassword must contain the plaintext password
//
// Returns:
//   - the created user (without password)
//   - apperrors.ErrConflict if user_name or user_email already exist
func (sr *sqlUserRepo) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
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
	if err := sr.pool.QueryRow(ctx, query, args...).Scan(&userID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode {
			return nil, apperrors.ErrConflict
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
//   - apperrors.ErrNotFound if no user matches id
func (sr *sqlUserRepo) GetUser(ctx context.Context, id int64) (*model.User, error) {
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

	user, err := scanUserDB(sr.pool.QueryRow(ctx, query, args...))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		subLogger.Error().Err(err).Msg("error executing select query")
		return nil, err
	}

	return user.ToModel(), nil
}

// GetUsers retrieves every user.
//
// Returns:
//   - the list of users (without passwords), empty if none exist
func (sr *sqlUserRepo) GetUsers(ctx context.Context) ([]*model.User, error) {
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

	rows, err := sr.pool.Query(ctx, query, args...)
	if err != nil {
		subLogger.Error().Err(err).Msg("error executing select query")
		return nil, err
	}
	defer rows.Close()

	items, err := pgxtool.ScanRows(rows, subLogger, scanUserDB)
	if err != nil {
		return nil, err
	}

	users := make([]*model.User, 0, len(items))
	for i := range items {
		users = append(users, items[i].ToModel())
	}

	return users, nil
}

// scanUserDB scans a single row into a UserDB, following column order in constants.Columns.
func scanUserDB(s pgxtool.Scanner) (UserDB, error) {
	var u UserDB
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
//   - apperrors.ErrUnauthorized if the credentials are invalid
func (sr *sqlUserRepo) LoginUser(ctx context.Context, userName, password string) (*model.LoginResponse, error) {
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

	user, err := scanUserDB(sr.pool.QueryRow(ctx, query, args...))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrUnauthorized
		}
		subLogger.Error().Err(err).Msg("error executing select query")
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.UserPassword), []byte(password)); err != nil {
		return nil, apperrors.ErrUnauthorized
	}

	token, err := generateToken(user.UserID)
	if err != nil {
		subLogger.Error().Err(err).Msg("error generating token")
		return nil, err
	}

	return &model.LoginResponse{Token: token}, nil
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
	return token.SignedString([]byte(env.GetString(enum.SecretKey, "")))
}
