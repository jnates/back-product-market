// Package storage wires the application's singleton PostgreSQL connection pool.
package storage

import (
	"sync"

	"backend_crudgo/configs"

	"github.com/jnates/go-toolkit/tools/sqlconnection/pgx"
)

var (
	once sync.Once
	pool pgx.DBPool
	err  error
)

// PostgresConnection returns the singleton PostgreSQL connection pool, opening it
// on first call and validating connectivity.
func PostgresConnection(config *configs.Config) (pgx.DBPool, error) {
	once.Do(func() {
		dbConfig := &pgx.DatabaseConfig{
			Host:     config.DBHost,
			Port:     config.DBPort,
			Database: config.DBName,
			Username: config.DBUser,
			Password: config.DBPassword,
			SSLMode:  config.DBSSLMode,
		}

		pool, err = pgx.NewClient(dbConfig, pgx.DefaultPoolConfig())
	})

	return pool, err
}

// PostgresCloseConnection closes the singleton connection pool, if it was opened.
func PostgresCloseConnection() {
	if pool != nil {
		pool.Close()
	}
}
