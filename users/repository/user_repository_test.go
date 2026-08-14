package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"backend_crudgo/pkg/kit/customErrors"
	"backend_crudgo/pkg/kit/enums"
	"backend_crudgo/users/constants"
	"backend_crudgo/users/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jnates/go-toolkit/tools/querybuilder"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// anyArgs builds n wildcard argument matchers.
func anyArgs(n int) []interface{} {
	args := make([]interface{}, n)
	for i := range args {
		args[i] = pgxmock.AnyArg()
	}
	return args
}

func TestUserRepository_CreateUser(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	repo := NewUserRepository(pool)
	user := &models.User{Name: "Test User", Email: "test@example.com", UserIdentifier: 1,
		UserPassword: "plaintext", UserTypeIdentifier: 1}

	query, _, errBuilder := querybuilder.NewInsertBuilder().
		Into(constants.UsersTable).
		Columns(constants.ColumnUserName, constants.ColumnUserIdentifier, constants.ColumnUserEmail,
			constants.ColumnUserPassword, constants.ColumnUserTypeIdentifier).
		Values(user.Name, user.UserIdentifier, user.Email, "hashed", user.UserTypeIdentifier).
		Returning(constants.ColumnUserID).
		Build()
	require.NoError(t, errBuilder)

	pool.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(anyArgs(5)...).
		WillReturnRows(pgxmock.NewRows([]string{constants.ColumnUserID}).AddRow(int64(1)))

	created, err := repo.CreateUser(context.Background(), user)

	require.NoError(t, err)
	assert.Equal(t, int64(1), created.UserID)
	assert.Empty(t, created.UserPassword)
	assert.NoError(t, pool.ExpectationsWereMet())
}

func TestUserRepository_CreateUser_Conflict(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	repo := NewUserRepository(pool)
	user := &models.User{Name: "Test User", UserPassword: "plaintext"}

	pool.ExpectQuery(".*").WithArgs(anyArgs(5)...).WillReturnError(&pgconn.PgError{Code: uniqueViolationCode})

	created, err := repo.CreateUser(context.Background(), user)

	assert.ErrorIs(t, err, customErrors.ErrConflict)
	assert.Nil(t, created)
}

func TestUserRepository_CreateUser_QueryError(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	repo := NewUserRepository(pool)
	user := &models.User{Name: "Test User", UserPassword: "plaintext"}

	pool.ExpectQuery(".*").WithArgs(anyArgs(5)...).WillReturnError(errors.New("db down"))

	created, err := repo.CreateUser(context.Background(), user)

	assert.Error(t, err)
	assert.False(t, errors.Is(err, customErrors.ErrConflict))
	assert.Nil(t, created)
}

func TestUserRepository_GetUser(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	repo := NewUserRepository(pool)

	query, _, errBuilder := querybuilder.NewSelectBuilder().
		Select(constants.Columns...).
		From(constants.UsersTable).
		Where(constants.ColumnUserID, querybuilder.OpEqual, int64(1)).
		Limit(1).
		Build()
	require.NoError(t, errBuilder)

	pool.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(int64(1)).
		WillReturnRows(pgxmock.NewRows(constants.Columns).
			AddRow(int64(1), "Test User", int64(1), "test@example.com", "hashed", int64(1)))

	user, err := repo.GetUser(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, int64(1), user.UserID)
	assert.Empty(t, user.UserPassword)
	assert.NoError(t, pool.ExpectationsWereMet())
}

func TestUserRepository_GetUser_NotFound(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	repo := NewUserRepository(pool)

	pool.ExpectQuery(".*").WithArgs(anyArgs(1)...).WillReturnError(pgx.ErrNoRows)

	user, err := repo.GetUser(context.Background(), 99)

	assert.ErrorIs(t, err, customErrors.ErrNotFound)
	assert.Nil(t, user)
}

func TestUserRepository_GetUser_QueryError(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	repo := NewUserRepository(pool)

	pool.ExpectQuery(".*").WithArgs(anyArgs(1)...).WillReturnError(errors.New("db down"))

	user, err := repo.GetUser(context.Background(), 1)

	assert.Error(t, err)
	assert.False(t, errors.Is(err, customErrors.ErrNotFound))
	assert.Nil(t, user)
}

func TestUserRepository_GetUsers(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	repo := NewUserRepository(pool)

	query, _, errBuilder := querybuilder.NewSelectBuilder().
		Select(constants.Columns...).
		From(constants.UsersTable).
		OrderBy(constants.ColumnUserID, querybuilder.Asc).
		Build()
	require.NoError(t, errBuilder)

	pool.ExpectQuery(regexp.QuoteMeta(query)).
		WillReturnRows(pgxmock.NewRows(constants.Columns).
			AddRow(int64(1), "User 1", int64(1), "u1@example.com", "hashed", int64(1)).
			AddRow(int64(2), "User 2", int64(1), "u2@example.com", "hashed", int64(1)))

	users, err := repo.GetUsers(context.Background())

	require.NoError(t, err)
	assert.Len(t, users, 2)
	assert.NoError(t, pool.ExpectationsWereMet())
}

func TestUserRepository_GetUsers_QueryError(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	repo := NewUserRepository(pool)

	pool.ExpectQuery(".*").WillReturnError(errors.New("db down"))

	users, err := repo.GetUsers(context.Background())

	assert.Error(t, err)
	assert.Nil(t, users)
}

func TestUserRepository_LoginUser(t *testing.T) {
	t.Setenv(enums.SecretKey, "test-secret")

	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	repo := NewUserRepository(pool)

	hashed, err := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
	require.NoError(t, err)

	pool.ExpectQuery(".*").WithArgs("testuser").
		WillReturnRows(pgxmock.NewRows(constants.Columns).
			AddRow(int64(1), "testuser", int64(1), "test@example.com", string(hashed), int64(1)))

	loginResponse, err := repo.LoginUser(context.Background(), "testuser", "correct-password")

	require.NoError(t, err)
	assert.NotEmpty(t, loginResponse.Token)
}

func TestUserRepository_LoginUser_WrongPassword(t *testing.T) {
	t.Setenv(enums.SecretKey, "test-secret")

	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	repo := NewUserRepository(pool)

	hashed, err := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
	require.NoError(t, err)

	pool.ExpectQuery(".*").WithArgs("testuser").
		WillReturnRows(pgxmock.NewRows(constants.Columns).
			AddRow(int64(1), "testuser", int64(1), "test@example.com", string(hashed), int64(1)))

	loginResponse, err := repo.LoginUser(context.Background(), "testuser", "wrong-password")

	assert.ErrorIs(t, err, customErrors.ErrUnauthorized)
	assert.Nil(t, loginResponse)
}

func TestUserRepository_LoginUser_NotFound(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	repo := NewUserRepository(pool)

	pool.ExpectQuery(".*").WithArgs("ghost").WillReturnError(pgx.ErrNoRows)

	loginResponse, err := repo.LoginUser(context.Background(), "ghost", "whatever")

	assert.ErrorIs(t, err, customErrors.ErrUnauthorized)
	assert.Nil(t, loginResponse)
}

func TestUserRepository_LoginUser_QueryError(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	repo := NewUserRepository(pool)

	pool.ExpectQuery(".*").WithArgs("testuser").WillReturnError(errors.New("db down"))

	loginResponse, err := repo.LoginUser(context.Background(), "testuser", "secret")

	assert.Error(t, err)
	assert.False(t, errors.Is(err, customErrors.ErrUnauthorized))
	assert.Nil(t, loginResponse)
}
