// Command api starts the product-market HTTP API.
package main

import (
	"backend_crudgo/infrastructure"
	"backend_crudgo/infrastructure/kit/enum"

	"github.com/jnates/go-toolkit/tools/env"
	"github.com/rs/zerolog/log"
)

// main starts the product-market API.
//
// @title Product Market API
// @version 1.0
// @description REST API for managing products and users with JWT-based authentication.
// @BasePath /api/market
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the JWT token.
func main() {
	log.Info().Msg("Starting API")
	infrastructure.InitLogger()

	port := env.GetString(enum.APIPort, "8080")
	infrastructure.Start(port)
}
