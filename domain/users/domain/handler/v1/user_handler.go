// Package v1 exposes the HTTP handlers for the users domain.
package v1

import (
	"net/http"

	"backend_crudgo/domain/users/domain/model"
	"backend_crudgo/domain/users/domain/service"
	"backend_crudgo/domain/users/infrastructure/persistence"
	"backend_crudgo/infrastructure/kit/apperrors"
	"backend_crudgo/infrastructure/middlewares"

	"github.com/jnates/go-toolkit/tools/customserver"
	pgxtool "github.com/jnates/go-toolkit/tools/sqlconnection/pgx"
	"github.com/labstack/echo/v4"
)

// LoginRequest is the request body accepted by POST /users/login.
type LoginRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"user_password"`
}

// UserHandler exposes the HTTP endpoints for users.
type UserHandler struct {
	Service service.UserService
}

// NewUserHandler wires the user repository, service and handler.
func NewUserHandler(pool pgxtool.DBPool) *UserHandler {
	return &UserHandler{
		Service: service.NewUserService(persistence.NewUserRepository(pool)),
	}
}

// CreateUser handles POST /users/register.
//
// @Description Register a new user
// @Tags Users
// @Accept json
// @Produce json
// @ID CreateUser
// @Param UserRequest body model.User true "User data"
// @Success 201 {object} customserver.GenericResponse{data=model.User}
// @Failure 400 {object} customserver.GenericResponse
// @Failure 409 {object} customserver.GenericResponse
// @Failure 500 {object} customserver.GenericResponse
// @Router /users/register [POST]
func (h *UserHandler) CreateUser(c echo.Context) error {
	ctx := c.Request().Context()

	var user model.User
	if err := c.Bind(&user); err != nil {
		return middlewares.HandleError(c, apperrors.ErrInvalidInput)
	}

	created, err := h.Service.CreateUser(ctx, &user)
	if err != nil {
		return middlewares.HandleError(c, err)
	}

	return c.JSON(http.StatusCreated, customserver.GenerateSuccessGenericResponse(created))
}

// LoginUser handles POST /users/login.
//
// @Description Authenticate a user and issue a JWT
// @Tags Users
// @Accept json
// @Produce json
// @ID LoginUser
// @Param LoginRequest body LoginRequest true "Login credentials"
// @Success 200 {object} customserver.GenericResponse{data=model.LoginResponse}
// @Failure 400 {object} customserver.GenericResponse
// @Failure 401 {object} customserver.GenericResponse
// @Failure 500 {object} customserver.GenericResponse
// @Router /users/login [POST]
func (h *UserHandler) LoginUser(c echo.Context) error {
	ctx := c.Request().Context()

	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return middlewares.HandleError(c, apperrors.ErrInvalidInput)
	}

	loginResponse, err := h.Service.LoginUser(ctx, req.UserName, req.Password)
	if err != nil {
		return middlewares.HandleError(c, err)
	}

	return c.JSON(http.StatusOK, customserver.GenerateSuccessGenericResponse(loginResponse))
}

// GetUsers handles GET /users.
//
// @Description Get every user
// @Tags Users
// @Produce json
// @ID GetUsers
// @Security BearerAuth
// @Success 200 {object} customserver.GenericResponse{data=[]model.User}
// @Failure 401 {object} customserver.GenericResponse
// @Failure 500 {object} customserver.GenericResponse
// @Router /users [GET]
func (h *UserHandler) GetUsers(c echo.Context) error {
	ctx := c.Request().Context()

	users, err := h.Service.GetUsers(ctx)
	if err != nil {
		return middlewares.HandleError(c, err)
	}

	return c.JSON(http.StatusOK, customserver.GenerateSuccessGenericResponse(users))
}
