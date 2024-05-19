package main

import (
	"backend_crudgo/infrastructure"
	"backend_crudgo/infrastructure/kit/enum"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("Starting API CMD")
	infrastructure.InitLogger()

	port := os.Getenv(enum.APIPort)
	infrastructure.Start(port)
}
