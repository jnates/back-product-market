// Package health exposes the application's health check endpoint.
package health

import (
	"net/http"

	pgxtool "github.com/jnates/go-toolkit/tools/sqlconnection/pgx"
	"github.com/labstack/echo/v4"
)

// Response reports the status of each dependency the API relies on.
type Response struct {
	Postgres string `json:"postgres"`
}

// Handler exposes the health check endpoint.
type Handler interface {
	// Check reports whether the API's dependencies (currently: PostgreSQL) are reachable.
	Check(c echo.Context) error
}

type handler struct {
	pool pgxtool.DBPool
}

// NewHandler builds a Handler backed by the given connection pool.
func NewHandler(pool pgxtool.DBPool) Handler {
	return &handler{pool: pool}
}

// Check handles GET /health.
//
// @Description Report the health of the API and its dependencies
// @Tags Health
// @Produce json
// @ID HealthCheck
// @Success 200 {object} Response
// @Router /health [GET]
func (h *handler) Check(c echo.Context) error {
	response := Response{Postgres: "ok"}

	if err := h.pool.Ping(c.Request().Context()); err != nil {
		response.Postgres = "error"
	}

	return c.JSON(http.StatusOK, response)
}
