// Package middleware contains cross-cutting Echo middlewares.
package middleware

import (
	"backend_crudgo/pkg/kit/enums"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jnates/go-toolkit/tools/env"
	"github.com/jnates/go-toolkit/tools/jwttools"
	jwtEcho "github.com/jnates/go-toolkit/tools/jwttools/echo"

	"github.com/labstack/echo/v4"
)

// NewAuthMiddleware builds an Echo middleware that validates the Bearer JWT
// on protected routes using jwttools, signed with the configured SECRET_KEY.
func NewAuthMiddleware() (echo.MiddlewareFunc, error) {
	secretKey := env.GetString(enums.SecretKey, "")

	validator, err := jwttools.NewValidator(jwttools.Config{
		ValidationKeyFunc: func(_ *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		},
	})
	if err != nil {
		return nil, err
	}

	authenticator, err := jwttools.NewAuthenticator(validator)
	if err != nil {
		return nil, err
	}

	return jwtEcho.JWTAuth(jwtEcho.Config{Authenticator: authenticator}), nil
}
