package customErrors

import (
	"errors"
	"net/http"

	"github.com/jnates/go-toolkit/tools/customserver"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// HandleError maps a domain error to the standard HTTP error response.
// Unknown errors are logged with full context and answered with a generic 500 message.
func HandleError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, ErrNotFound):
		return c.JSON(http.StatusNotFound, customserver.GenerateErrorGenericResponse(
			http.StatusNotFound, "resource not found", nil))
	case errors.Is(err, ErrInvalidInput):
		return c.JSON(http.StatusBadRequest, customserver.GenerateErrorGenericResponse(
			http.StatusBadRequest, "invalid input", nil))
	case errors.Is(err, ErrUnauthorized):
		return c.JSON(http.StatusUnauthorized, customserver.GenerateErrorGenericResponse(
			http.StatusUnauthorized, "unauthorized", nil))
	case errors.Is(err, ErrConflict):
		return c.JSON(http.StatusConflict, customserver.GenerateErrorGenericResponse(
			http.StatusConflict, "resource already exists", nil))
	default:
		log.Error().Err(err).Msg("internal server error")
		return c.JSON(http.StatusInternalServerError, customserver.GenerateErrorGenericResponse(
			http.StatusInternalServerError, "internal server error", nil))
	}
}
