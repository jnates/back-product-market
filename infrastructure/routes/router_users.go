package routes

import (
	v1 "backend_crudgo/domain/users/domain/handler/v1"
	"backend_crudgo/infrastructure/kit/enum"

	pgxtool "github.com/jnates/go-toolkit/tools/sqlconnection/pgx"
	"github.com/labstack/echo/v4"
)

// RegisterUserRoutes mounts the user endpoints. Register/login are public;
// listing users requires a valid JWT via authMiddleware.
func RegisterUserRoutes(server *echo.Echo, pool pgxtool.DBPool, authMiddleware echo.MiddlewareFunc) {
	handler := v1.NewUserHandler(pool)

	public := server.Group(enum.BasePathUser)
	public.POST(enum.RegisterPath, handler.CreateUser)
	public.POST(enum.LoginUserPath, handler.LoginUser)

	protected := server.Group(enum.BasePathUser, authMiddleware)
	protected.GET("", handler.GetUsers)
}
