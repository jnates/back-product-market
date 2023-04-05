package v1

import (
	"backend_crudgo/infrastructure/kit/enum"
	"backend_crudgo/infrastructure/middlewares"
	"encoding/json"
	"fmt"
	"net/http"

	"backend_crudgo/domain/users/domain/model"
	"backend_crudgo/domain/users/domain/service"
	"backend_crudgo/domain/users/infrastructure/persistence"
	"backend_crudgo/infrastructure/database"
)

// UserRouter router
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

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		_ = middlewares.HTTPError(w, r, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	result, err := prod.Service.CreateUserHandler(ctx, &user)
	if err != nil {
		_ = middlewares.HTTPError(w, r, http.StatusConflict, "Conflict", err.Error())
		return
	}

	w.Header().Add(enum.Location, fmt.Sprintf("%s%s", r.URL.String(), result))
	_ = middlewares.JSON(w, r, http.StatusCreated, result)
}

// LoginUserHandler login a user.
func (prod *UserRouter) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var user model.User
	var ctx = r.Context()

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		_ = middlewares.HTTPError(w, r, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	userResponse, err := prod.Service.LoginUserHandler(ctx, &user)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonBytes, err := json.Marshal(userResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonBytes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (prod *UserRouter) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	userResponse, err := prod.Service.GetUsersHandler(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonBytes, err := json.Marshal(userResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonBytes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
