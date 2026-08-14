// Package users registers the Echo routes for the users domain.
package users

import (
	"backend_crudgo/pkg/kit/enums"
	"backend_crudgo/users/interfaces"

	"github.com/labstack/echo/v4"
)

// RouteUsers mounts the user endpoints on a given route group.
type RouteUsers interface {
	Resource(g *echo.Group)
}

type routeUsers struct {
	handler        interfaces.UserHandler
	authMiddleware echo.MiddlewareFunc
}

// NewRouteUsers builds a RouteUsers backed by the given handler.
// Register/login are public; listing users requires a valid JWT.
func NewRouteUsers(handler interfaces.UserHandler, authMiddleware echo.MiddlewareFunc) RouteUsers {
	return &routeUsers{handler: handler, authMiddleware: authMiddleware}
}

// Resource registers the user routes on g.
func (r *routeUsers) Resource(g *echo.Group) {
	g.POST(enums.RegisterPath, r.handler.CreateUser)
	g.POST(enums.LoginUserPath, r.handler.LoginUser)
	g.GET("", r.handler.GetUsers, r.authMiddleware)
}
