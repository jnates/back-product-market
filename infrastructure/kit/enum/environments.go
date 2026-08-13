// Package enum centralizes the literal string constants used across the application
// (environment variable names, route paths, HTTP header names).
package enum

// Environment variable names read at bootstrap.
const (
	SecretKey   string = "SECRET_KEY"
	APIPort     string = "API_PORT"
	DBHost      string = "DB_HOST"
	DBUser      string = "DB_USER"
	DBPassword  string = "DB_PASSWORD"
	DBName      string = "DB_NAME"
	DBPort      string = "DB_PORT"
	DBSSLMode   string = "DB_SSL_MODE"
	LoggerDebug string = "LOGGER_DEBUG"
)
