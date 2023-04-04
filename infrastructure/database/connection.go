package database

import (
	"backend_crudgo/infrastructure/kit/enum"
	"database/sql"
	"fmt"
	"os"

	// registering database driver
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

// New returns a new instance of Data with the database connection ready.
func New() (*DataDB, error) {
	db, err := getConnection()
	if err != nil {
		return nil, err
	}

	return &DataDB{DB: db}, nil
}

// DataDB is struct for library database/sql
type DataDB struct {
	DB *sql.DB
}

func getConnection() (*sql.DB, error) {
	DbHost := os.Getenv(enum.DBHost) //"127.0.0.1"
	DbDriver := os.Getenv(enum.DBDriver)
	DbUser := os.Getenv(enum.DBUser)
	DbPassword := os.Getenv(enum.DBPassword)
	DbName := os.Getenv(enum.DBName)
	DbPort := os.Getenv(enum.DBPort)
	uri := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)

	db, err := sql.Open(DbDriver, uri)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	log.Info().Msg("Connected to database")
	return db, nil
}
