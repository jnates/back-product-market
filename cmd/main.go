// Command main starts the product-market API.
//
// @title Product Market API
// @version 1.0
// @description REST API for managing products and users with JWT-based authentication.
// @BasePath /api/market
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the JWT token.
package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"backend_crudgo/configs"
	"backend_crudgo/configs/generals/injector"
	"backend_crudgo/configs/generals/router"
	"backend_crudgo/configs/storage"
	"backend_crudgo/pkg/kit/enums"

	"github.com/jnates/go-toolkit/tools/logger/zerolog"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

const shutdownTimeout = 30 * time.Second

func main() {
	if err := godotenv.Load(); err != nil {
		log.Error().Err(err).Msg("Error loading file .env")
	}

	config := configs.LoadConfiguration()

	zerolog.InitLogger(enums.App, config.LoggerDebug)
	log.Info().Msg("Starting API")

	container := injector.BuildContainer(config)

	err := container.Invoke(func(server *echo.Echo, route *router.Router) {
		route.Init()
		start(server, config.APIPort)
	})

	storage.PostgresCloseConnection()

	if err != nil {
		log.Fatal().Err(err).Msg("could not start application")
	}
}

// start listens for HTTP requests in the background and blocks until a graceful shutdown completes.
func start(server *echo.Echo, port string) {
	go func() {
		if err := server.Start(":" + port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msgf("could not listen on port %s", port)
		}
	}()

	log.Info().Msgf("API is ready to handle requests on port %s", port)
	gracefulShutdown(server)
}

// gracefulShutdown blocks until an interrupt signal is received, then shuts the server down.
func gracefulShutdown(server *echo.Echo) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	sig := <-quit
	log.Info().Msgf("API is shutting down %s", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	err := server.Shutdown(ctx)
	cancel()

	if err != nil {
		log.Fatal().Err(err).Msg("could not gracefully shutdown the API")
	}

	log.Info().Msg("API stopped")
}
