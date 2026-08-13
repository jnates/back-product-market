// Package database wires the PostgreSQL connection pool used by every repository.
package database

import (
	"backend_crudgo/infrastructure/kit/enum"

	"github.com/jnates/go-toolkit/tools/env"
	"github.com/jnates/go-toolkit/tools/sqlconnection/pgx"
)

// New opens a PostgreSQL connection pool configured from environment variables
// and validates connectivity before returning.
func New() (pgx.DBPool, error) {
	dbConfig := &pgx.DatabaseConfig{
		Host:     env.GetString(enum.DBHost, ""),
		Port:     int(env.GetInt64(enum.DBPort, 0)),
		Database: env.GetString(enum.DBName, ""),
		Username: env.GetString(enum.DBUser, ""),
		Password: env.GetString(enum.DBPassword, ""),
		SSLMode:  env.GetString(enum.DBSSLMode, "disable"),
	}

	return pgx.NewClient(dbConfig, pgx.DefaultPoolConfig())
}
