package infrastructure

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"backend_crudgo/infrastructure/database"
	"backend_crudgo/infrastructure/kit/enum"

	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog/log"
)

// Server is a base Server configuration.
type Server struct {
	*http.Server
}

// ServeHTTP implements the http.Handler interface for the server type.
func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.Handler.ServeHTTP(w, r)
}

// newServer initialized a Routes Server with configuration.
func newServer(port string, conn *database.DataDB) *Server {

	router := chi.NewRouter()

	router.Use(chiMiddleware.RequestID)
	router.Use(chiMiddleware.RealIP)
	router.Use(chiMiddleware.Logger)
	router.Use(chiMiddleware.Recoverer)
	// router.Use(middleware.CORSMiddleware)

	//default path to be used in the health checker
	router.Mount(enum.BasePath, RoutesProducts(conn))

	s := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &Server{s}
}

func (srv *Server) gracefulShutdown() {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	sig := <-quit
	log.Info().Msgf("CMD is shutting down %s", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	srv.SetKeepAlivesEnabled(false)
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Msgf("could not gracefully shutdown the cmd %s", err.Error())
	}

	log.Info().Msg("CMD Stopped")
}

// Start initialize server
func (srv *Server) Start() {
	log.Info().Msg("Starting API cmd")

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Msgf("Could not listen on %s rv due to %s rv", srv.Addr, err.Error())
		}
	}()

	log.Info().Msgf("CMD is ready to handle requests %s", srv.Addr)
	srv.gracefulShutdown()
}

// Start aa
func Start(port string) {
	// connection to the database.
	db, err := database.New()
	if err != nil {
		return
	}

	defer func() {
		err = db.DB.Close()
		if err != nil {
			log.Fatal().Msgf("Could not close BD : [error] %s", err.Error())
		}
	}()

	server := newServer(port, db)
	// start the server.
	server.Start()
}
