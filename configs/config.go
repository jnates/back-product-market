// Package configs loads the application configuration at bootstrap.
package configs

import (
	"backend_crudgo/pkg/kit/enums"

	"github.com/jnates/go-toolkit/tools/env"
	"github.com/rs/zerolog/log"
)

// Config holds every environment-driven setting the application needs at startup.
type Config struct {
	APIPort     string
	DBHost      string
	DBPort      int
	DBName      string
	DBUser      string
	DBPassword  string
	DBSSLMode   string
	SecretKey   string
	LoggerDebug bool
}

// LoadConfiguration reads Config from environment variables, applying sane defaults.
// It fails fast (log.Fatal) if a security-critical variable is missing, rather than
// letting the application boot in an insecure state (e.g. JWTs signed/validated
// with an empty secret).
func LoadConfiguration() *Config {
	config := &Config{
		APIPort:     env.GetString(enums.APIPort, "8080"),
		DBHost:      env.GetString(enums.DBHost, ""),
		DBPort:      int(env.GetInt64(enums.DBPort, 0)),
		DBName:      env.GetString(enums.DBName, ""),
		DBUser:      env.GetString(enums.DBUser, ""),
		DBPassword:  env.GetString(enums.DBPassword, ""),
		DBSSLMode:   env.GetString(enums.DBSSLMode, "disable"),
		SecretKey:   env.GetString(enums.SecretKey, ""),
		LoggerDebug: env.GetBoolean(enums.LoggerDebug, false),
	}

	if config.SecretKey == "" {
		log.Fatal().Msgf("%s environment variable is required", enums.SecretKey)
	}

	return config
}
