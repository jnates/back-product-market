// Package routes registers the Echo routes for each domain.
package routes

import (
	"backend_crudgo/domain/products/constants"
	v1 "backend_crudgo/domain/products/domain/handler/v1"
	"backend_crudgo/infrastructure/kit/enum"

	pgxtool "github.com/jnates/go-toolkit/tools/sqlconnection/pgx"
	"github.com/labstack/echo/v4"
)

// idParamPath is the route path for endpoints scoped to a single product ID.
const idParamPath = "/:" + constants.ID

// RegisterProductRoutes mounts the product endpoints, protected by authMiddleware.
func RegisterProductRoutes(server *echo.Echo, pool pgxtool.DBPool, authMiddleware echo.MiddlewareFunc) {
	handler := v1.NewProductHandler(pool)

	group := server.Group(enum.BasePath+"/products", authMiddleware)
	group.POST("", handler.CreateProduct)
	group.GET("", handler.GetProducts)
	group.GET(idParamPath, handler.GetProduct)
	group.PUT(idParamPath, handler.UpdateProduct)
	group.DELETE(idParamPath, handler.DeleteProduct)
}
