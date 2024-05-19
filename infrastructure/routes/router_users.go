package infrastructure

import (
	"net/http"

	v1 "backend_crudgo/domain/users/domain/handler/v1"
	"backend_crudgo/infrastructure/database"
	"backend_crudgo/infrastructure/kit/enum"
	"backend_crudgo/infrastructure/middlewares"

	"github.com/go-chi/chi"
)

// RoutesUsers creates a new router for handling user related requests.
// The function takes a database connection as an argument and returns an HTTP handler.
func RoutesUsers(conn *database.DataDB) http.Handler {
	router := chi.NewRouter()
	users := v1.NewUserHandler(conn)
	router.Mount("/", routesUser(users))
	return router
}

// routesUser creates a new router for handling user related requests.
// The function takes a user handler as an argument and returns an HTTP handler.
func routesUser(handler *v1.UserRouter) http.Handler {
	router := chi.NewRouter()
	router.Post(enum.LoginUserPath, handler.LoginUserHandler)
	router.Post(enum.RegisterPath, handler.CreateUserHandler)
	router.Post("/generatetoken", handler.AuthHandler)

	router.Group(func(protected chi.Router) {
		protected.Use(middlewares.AuthMiddleware)
		protected.Get("/protected", handler.ProtectedHandler)
	})

	return router
}
