package infrastructure

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"backend_crudgo/infrastructure/database"
	"backend_crudgo/infrastructure/kit/enum"
	"backend_crudgo/infrastructure/middlewares"
	"backend_crudgo/infrastructure/routes"

	customserverEcho "github.com/jnates/go-toolkit/tools/customserver/echo"
	pgxtool "github.com/jnates/go-toolkit/tools/sqlconnection/pgx"

	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

const shutdownTimeout = 30 * time.Second

// newServer builds an Echo server with the product and user routes mounted.
func newServer(pool pgxtool.DBPool) (*echo.Echo, error) {
	authMiddleware, err := middlewares.NewAuthMiddleware()
	if err != nil {
		return nil, err
	}

	server := customserverEcho.NewServer()
	server.GET(enum.BasePath+"/docs/*", echoSwagger.WrapHandler)
	routes.RegisterProductRoutes(server, pool, authMiddleware)
	routes.RegisterUserRoutes(server, pool, authMiddleware)

	return server, nil
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

// Start opens the database connection pool and starts the HTTP server.
func Start(port string) {
	pool, err := database.New()
	if err != nil {
		log.Fatal().Err(err).Msg("could not connect to database")
	}

	server, err := newServer(pool)
	if err != nil {
		pool.Close()
		log.Fatal().Err(err).Msg("could not build server")
	}
	defer pool.Close()

	go func() {
		if err := server.Start(":" + port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msgf("could not listen on port %s", port)
		}
	}()

	log.Info().Msgf("API is ready to handle requests on port %s", port)
	gracefulShutdown(server)
}
