// Package infrastructure wires the application's cross-cutting concerns (logger, HTTP server, DB).
package infrastructure

import (
	"backend_crudgo/infrastructure/kit/enum"

	"github.com/jnates/go-toolkit/tools/env"
	toolkitLogger "github.com/jnates/go-toolkit/tools/logger/zerolog"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

// InitLogger configures the global zerolog logger via the go-toolkit logger tool.
func InitLogger() {
	toolkitLogger.InitLogger(enum.App, env.GetBoolean(enum.LoggerDebug, false))
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Error().Err(err).Msg("Error loading file .env")
	}
}
