package persistence

import (
	"backend_crudgo/domain/users/domain/model"
	repoDomain "backend_crudgo/domain/users/domain/repository"
	"backend_crudgo/infrastructure/database"
	"backend_crudgo/infrastructure/kit/enum"
	response "backend_crudgo/types"
	"context"
	"database/sql"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type sqlUserRepo struct {
	Conn *database.DataDB
}

// NewUserRepository Should initialize the dependencies for this service.
func NewUserRepository(Conn *database.DataDB) repoDomain.UserRepository {
	return &sqlUserRepo{
		Conn: Conn,
	}
}

// CreateUserHandler creates a new user in the database.
func (sr *sqlUserRepo) CreateUserHandler(ctx context.Context, user *model.User) (*response.CreateResponse, error) {
	var idResult string

	stmt, err := sr.Conn.DB.PrepareContext(ctx, InsertUser)
	if err != nil {
		return &response.CreateResponse{}, err
	}

	defer func() {
		err = stmt.Close()
		if err != nil {
			log.Error().Msgf("Could not close testament : [error] %s", err.Error())
		}
	}()

	user.UserPassword = hashPassword(user.UserPassword)
	row := stmt.QueryRowContext(ctx, &user.UserID, &user.Name, &user.UserIdentifier, &user.Email,
		&user.UserPassword, &user.UserTypeIdentifier)

	if err = row.Scan(&idResult); err != sql.ErrNoRows {
		return &response.CreateResponse{}, err
	}

	return &response.CreateResponse{
		Message: "User created",
	}, nil
}

// LoginUserHandler logs in a user by checking if their password is correct.
func (sr *sqlUserRepo) LoginUserHandler(ctx context.Context, user *model.User) (*response.GenericUserResponse, error) {
	stmt, err := sr.Conn.DB.PrepareContext(ctx, SelectLoginUser)
	if err != nil {
		return &response.GenericUserResponse{}, err
	}

	defer func() {
		err = stmt.Close()
		if err != nil {
			log.Error().Msgf("Could not close testament : [error] %s", err.Error())
		}
	}()

	row := stmt.QueryRowContext(ctx, user.Name)
	currentUser := &model.User{}

	if err = row.Scan(&currentUser.UserID, &currentUser.Name, &currentUser.Email, &currentUser.UserIdentifier, &currentUser.UserPassword, &currentUser.UserTypeIdentifier); err != nil {
		return &response.GenericUserResponse{Error: err.Error()}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(currentUser.UserPassword), []byte(user.UserPassword))

	if err != nil {
		return &response.GenericUserResponse{Error: "Password incorrect"}, nil
	}
	token, nil := generateToken(currentUser.UserID)

	return &response.GenericUserResponse{
		Message: "Success",
		User:    token,
	}, nil
}

// GetUserHandler retrieves a specific user from the database.
func (sr *sqlUserRepo) GetUserHandler(ctx context.Context, id string) (*response.GenericUserResponse, error) {
	stmt, err := sr.Conn.DB.PrepareContext(ctx, SelectUser)
	if err != nil {
		return &response.GenericUserResponse{}, err
	}

	defer func() {
		err = stmt.Close()
		if err != nil {
			log.Error().Msgf("Could not close testament : [error] %s", err.Error())
		}
	}()

	row := stmt.QueryRowContext(ctx, id)
	user := &model.User{}

	if err = row.Scan(&user.UserID, &user.Name, &user.Email, &user.UserIdentifier, &user.UserPassword, &user.DateCreated, &user.UserModify, &user.DateModify); err != nil {
		return &response.GenericUserResponse{Error: err.Error()}, err
	}

	return &response.GenericUserResponse{
		Message: "Get user success",
		User:    user,
	}, nil
}

// GetUsersHandler retrieves a list of all users from the database.
func (sr *sqlUserRepo) GetUsersHandler(ctx context.Context) (*response.GenericUserResponse, error) {
	stmt, err := sr.Conn.DB.PrepareContext(ctx, SelectUsers)
	if err != nil {
		return &response.GenericUserResponse{}, nil
	}

	defer func() {
		err = stmt.Close()
		if err != nil {
			log.Error().Msgf("Could not close testament : [error] %s", err.Error())
		}
	}()
	row, err := sr.Conn.DB.QueryContext(ctx, SelectUsers)

	var users []*model.User
	for row.Next() {
		var user = &model.User{}
		if err = row.Scan(&user.UserID, &user.Name, &user.Email, &user.UserIdentifier, &user.UserPassword, &user.DateCreated, &user.UserModify, &user.DateModify); err != nil {
			return &response.GenericUserResponse{Error: err.Error()}, err
		}
		users = append(users, user)
	}

	return &response.GenericUserResponse{
		Message: "Get user success",
		User:    users,
	}, nil
}

// hashPassword hashes a plain text password.
func hashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Msgf("Could not hash password: [error] %s", err.Error())
	}
	return string(hashedPassword)
}

// generateToken generates a new JWT token.
func generateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := os.Getenv(enum.SecretKey)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return enum.EmptyString, err
	}

	return signedToken, nil
}
