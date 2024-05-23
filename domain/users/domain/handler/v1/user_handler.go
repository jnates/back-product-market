package v1

import (
	"backend_crudgo/infrastructure/kit/enum"
	"backend_crudgo/infrastructure/kit/tool"
	"encoding/json"
	"fmt"
	"net/http"

	"backend_crudgo/domain/users/domain/model"
	"backend_crudgo/domain/users/domain/service"
	"backend_crudgo/domain/users/infrastructure/persistence"
	"backend_crudgo/infrastructure/database"
)

// UserRouter is a struct that contains a UserService instance. It is used to create an HTTP router for user-related endpoints.
type UserRouter struct {
	Service service.UserService
}

// NewUserHandler initializes the dependencies for this service.
func NewUserHandler(db *database.DataDB) *UserRouter {
	return &UserRouter{
		Service: service.NewUserService(persistence.NewUserRepository(db)),
	}
}

// CreateUserHandler handles user creation.
func (ur *UserRouter) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user model.User
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		tool.WriteJSONResponseWithMarshalling(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := ur.Service.CreateUser(ctx, &user)
	if err != nil {
		tool.WriteJSONResponseWithMarshalling(w, http.StatusConflict, err.Error())
		return
	}

	w.Header().Add(enum.Location, fmt.Sprintf("%s%s", r.URL.String(), result))
	tool.WriteJSONResponseWithMarshalling(w, http.StatusCreated, result)
}

// LoginUserHandler handles user login.
func (ur *UserRouter) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var user model.User
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		tool.WriteJSONResponseWithMarshalling(w, http.StatusBadRequest, err.Error())
		return
	}

	userResponse, err := ur.Service.LoginUser(ctx, &user)
	if err != nil {
		tool.WriteJSONResponseWithMarshalling(w, http.StatusInternalServerError, err.Error())
		return
	}

	tool.WriteJSONResponseWithMarshalling(w, http.StatusOK, userResponse)
}

// GetUsersHandler handles retrieving users.
func (ur *UserRouter) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userResponse, err := ur.Service.GetUsers(ctx)
	if err != nil {
		tool.WriteJSONResponseWithMarshalling(w, http.StatusInternalServerError, err.Error())
		return
	}

	tool.WriteJSONResponseWithMarshalling(w, http.StatusOK, userResponse)
}

// AuthHandler maneja la generación de tokens de autenticación.
func (ur *UserRouter) AuthHandler(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	tokenResponse, err := ur.Service.GenerateToken(r.Context(), &user)
	if err != nil {
		http.Error(w, "Error generando el token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"token":     tokenResponse.Token,
		"tokenType": tokenResponse.TokenType,
	})
}

// ProtectedHandler handles requests to protected routes.
func (ur *UserRouter) ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"message": "This is a protected route"}
	tool.WriteJSONResponseWithMarshalling(w, http.StatusOK, response)
}
