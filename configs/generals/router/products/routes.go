// Package products registers the Echo routes for the products domain.
package products

import (
	"backend_crudgo/products/constants"
	"backend_crudgo/products/interfaces"

	"github.com/labstack/echo/v4"
)

// RouteProducts mounts the product endpoints on a given route group.
type RouteProducts interface {
	Resource(g *echo.Group)
}

type routeProducts struct {
	handler        interfaces.ProductHandler
	authMiddleware echo.MiddlewareFunc
}

// NewRouteProducts builds a RouteProducts backed by the given handler, protected by authMiddleware.
func NewRouteProducts(handler interfaces.ProductHandler, authMiddleware echo.MiddlewareFunc) RouteProducts {
	return &routeProducts{handler: handler, authMiddleware: authMiddleware}
}

// Resource registers the product routes (all protected by JWT) on g.
func (r *routeProducts) Resource(g *echo.Group) {
	g.Use(r.authMiddleware)

	g.POST("", r.handler.CreateProduct)
	g.GET("", r.handler.GetProducts)
	g.GET("/:"+constants.ID, r.handler.GetProduct)
	g.PUT("/:"+constants.ID, r.handler.UpdateProduct)
	g.DELETE("/:"+constants.ID, r.handler.DeleteProduct)
}
