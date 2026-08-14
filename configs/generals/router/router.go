// Package router assembles the application's Echo routes.
package router

import (
	"backend_crudgo/configs/generals/router/health"
	productsRoutes "backend_crudgo/configs/generals/router/products"
	usersRoutes "backend_crudgo/configs/generals/router/users"
	"backend_crudgo/pkg/kit/enums"

	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/labstack/echo/v4"
)

// Router mounts every domain's routes on the Echo server.
type Router struct {
	server   *echo.Echo
	products productsRoutes.RouteProducts
	users    usersRoutes.RouteUsers
	health   health.Handler
}

// NewRouter builds a Router from the Echo server and each domain's route resource.
func NewRouter(
	server *echo.Echo,
	products productsRoutes.RouteProducts,
	users usersRoutes.RouteUsers,
	health health.Handler,
) *Router {
	return &Router{server: server, products: products, users: users, health: health}
}

// Init mounts the health check, Swagger UI, and every domain's routes on the server.
func (r *Router) Init() {
	r.server.GET(enums.HealthPath, r.health.Check)
	r.server.GET(enums.BasePath+"/docs/*", echoSwagger.WrapHandler)

	productsGroup := r.server.Group(enums.BasePath + "/products")
	r.products.Resource(productsGroup)

	usersGroup := r.server.Group(enums.BasePathUser)
	r.users.Resource(usersGroup)
}
