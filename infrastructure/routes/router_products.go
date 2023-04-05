package infrastructure

import (
	v1 "backend_crudgo/domain/products/domain/handler/v1"
	"backend_crudgo/infrastructure/database"
	"backend_crudgo/infrastructure/middlewares"
	"net/http"

	"github.com/go-chi/chi"
)

// RoutesProducts creates a new router for handling product related requests.
// The function takes a database connection as an argument and returns an HTTP handler.
func RoutesProducts(conn *database.DataDB) http.Handler {
	router := chi.NewRouter()
	products := v1.NewProductHandler(conn) //domain
	router.Mount("/products", routesProduct(products))
	return router
}

// routesProduct creates a new router for handling product related requests.
// The function takes a product handler as an argument and returns an HTTP handler.
func routesProduct(handler *v1.ProductRouter) http.Handler {
	router := chi.NewRouter()
	router.Use(middlewares.AuthMiddleware)
	router.Post("/", handler.CreateProductHandler)
	router.Get("/", handler.GetProductsHandler)
	router.Get("/{id}", handler.GetProductHandler)
	router.Put("/{id}", handler.UpdateProductHandler)
	router.Delete("/{id}", handler.DeleteProductHandler)
	return router
}
