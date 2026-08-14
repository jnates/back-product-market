// Package handler exposes the HTTP handlers for the users domain.
package handler

import (
	"net/http"

	"backend_crudgo/pkg/kit/customErrors"
	"backend_crudgo/pkg/kit/validation"
	"backend_crudgo/users/interfaces"
	"backend_crudgo/users/models"

	"github.com/jnates/go-toolkit/tools/customserver"
	"github.com/labstack/echo/v4"
)

// LoginRequest is the request body accepted by POST /users/login.
type LoginRequest struct {
	UserName string `json:"user_name" validate:"required"`
	Password string `json:"user_password" validate:"required"`
}

// CreateUserRequest is the request body accepted by POST /users/register.
type CreateUserRequest struct {
	Name               string `json:"user_name" validate:"required,min=2,max=100"`
	Email              string `json:"user_email" validate:"required,email"`
	UserIdentifier     int64  `json:"user_identifier" validate:"required"`
	UserPassword       string `json:"user_password" validate:"required,min=8"`
	UserTypeIdentifier int64  `json:"user_type_identifier" validate:"required"`
}

// userHandler implements interfaces.UserHandler, delegating to the user service.
type userHandler struct {
	service interfaces.UserService
}

// NewUserHandler builds a UserHandler backed by the given service.
func NewUserHandler(service interfaces.UserService) interfaces.UserHandler {
	return &userHandler{service: service}
}

// CreateUser handles POST /users/register.
//
// Parameters:
//   - c: echo context containing the HTTP request and response
//
// Returns:
//   - error: HTTP error or nil on success
//
// @Description Register a new user
// @Tags Users
// @Accept json
// @Produce json
// @ID CreateUser
// @Param UserRequest body CreateUserRequest true "User data"
// @Success 201 {object} customserver.GenericResponse{data=models.User}
// @Failure 400 {object} customserver.GenericResponse
// @Failure 409 {object} customserver.GenericResponse
// @Failure 422 {object} customserver.GenericResponse
// @Failure 500 {object} customserver.GenericResponse
// @Router /users/register [POST]
func (h *userHandler) CreateUser(c echo.Context) error {
	ctx := c.Request().Context()

	var req CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return customErrors.HandleError(c, customErrors.ErrInvalidInput)
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, customserver.GenerateErrorGenericResponse(
			http.StatusUnprocessableEntity, "validation failed", validation.FormatErrors(err)))
	}

	user := models.User{
		Name:               req.Name,
		Email:              req.Email,
		UserIdentifier:     req.UserIdentifier,
		UserPassword:       req.UserPassword,
		UserTypeIdentifier: req.UserTypeIdentifier,
	}

	created, err := h.service.CreateUser(ctx, &user)
	if err != nil {
		return customErrors.HandleError(c, err)
	}

	return c.JSON(http.StatusCreated, customserver.GenerateSuccessGenericResponse(created))
}

// LoginUser handles POST /users/login.
//
// Parameters:
//   - c: echo context containing the HTTP request and response
//
// Returns:
//   - error: HTTP error or nil on success
//
// @Description Authenticate a user and issue a JWT
// @Tags Users
// @Accept json
// @Produce json
// @ID LoginUser
// @Param LoginRequest body LoginRequest true "Login credentials"
// @Success 200 {object} customserver.GenericResponse{data=models.LoginResponse}
// @Failure 400 {object} customserver.GenericResponse
// @Failure 401 {object} customserver.GenericResponse
// @Failure 422 {object} customserver.GenericResponse
// @Failure 500 {object} customserver.GenericResponse
// @Router /users/login [POST]
func (h *userHandler) LoginUser(c echo.Context) error {
	ctx := c.Request().Context()

	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return customErrors.HandleError(c, customErrors.ErrInvalidInput)
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, customserver.GenerateErrorGenericResponse(
			http.StatusUnprocessableEntity, "validation failed", validation.FormatErrors(err)))
	}

	loginResponse, err := h.service.LoginUser(ctx, req.UserName, req.Password)
	if err != nil {
		return customErrors.HandleError(c, err)
	}

	return c.JSON(http.StatusOK, customserver.GenerateSuccessGenericResponse(loginResponse))
}

// GetUsers handles GET /users.
//
// Parameters:
//   - c: echo context containing the HTTP request and response
//
// Returns:
//   - error: HTTP error or nil on success
//
// @Description Get every user
// @Tags Users
// @Produce json
// @ID GetUsers
// @Security BearerAuth
// @Success 200 {object} customserver.GenericResponse{data=[]models.User}
// @Failure 401 {object} customserver.GenericResponse
// @Failure 500 {object} customserver.GenericResponse
// @Router /users [GET]
func (h *userHandler) GetUsers(c echo.Context) error {
	ctx := c.Request().Context()

	users, err := h.service.GetUsers(ctx)
	if err != nil {
		return customErrors.HandleError(c, err)
	}

	return c.JSON(http.StatusOK, customserver.GenerateSuccessGenericResponse(users))
}
