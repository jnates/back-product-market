package v1

import (
	"backend_crudgo/infrastructure/kit/enum"
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

// NewUserHandler Should initialize the dependencies for this service.
func NewUserHandler(db *database.DataDB) *UserRouter {
	return &UserRouter{
		Service: service.NewUserService(persistence.NewUserRepository(db)),
	}
}

// CreateUserHandler Created initialize handler user.
func (prod *UserRouter) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user model.User
	var ctx = r.Context()

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		writeJSONResponseWithMarshalling(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := prod.Service.CreateUser(ctx, &user)
	if err != nil {
		writeJSONResponseWithMarshalling(w, http.StatusConflict, err.Error())
		return
	}

	w.Header().Add(enum.Location, fmt.Sprintf("%s%s", r.URL.String(), result))
	writeJSONResponseWithMarshalling(w, http.StatusCreated, result)
}

// LoginUserHandler is the HTTP handler for user login. It receives an HTTP request with a JSON body containing user credentials.
// It verifies the user's authenticity through the user service and returns a JSON response containing user information and an authentication token upon success.
// If there is an error processing the request, it returns an appropriate HTTP error response.
func (prod *UserRouter) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var user model.User
	var ctx = r.Context()

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		writeJSONResponseWithMarshalling(w, http.StatusBadRequest, err.Error())
		return
	}

	userResponse, err := prod.Service.LoginUser(ctx, &user)

	if err != nil {
		writeJSONResponseWithMarshalling(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSONResponseWithMarshalling(w, http.StatusOK, userResponse)
}

// GetUsersHandler is the HTTP handler for retrieving users.
// It calls the user service to retrieve the list of users and returns a JSON response containing.
// the user information upon success.
// If there is an error processing the request, it returns an appropriate HTTP error response.
func (prod *UserRouter) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	userResponse, err := prod.Service.GetUsers(ctx)
	if err != nil {
		writeJSONResponseWithMarshalling(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSONResponseWithMarshalling(w, http.StatusOK, userResponse)
}

func writeJSONResponseWithMarshalling(w http.ResponseWriter, statusCode int, data interface{}) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if _, err = w.Write(jsonBytes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
