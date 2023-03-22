package database

import (
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
	DbHost := os.Getenv("DB_HOST") //"127.0.0.1"
	DbDriver := os.Getenv("DB_DRIVER")
	DbUser := os.Getenv("DB_USER")
	DbPassword := os.Getenv("DB_PASSWORD")
	DbName := os.Getenv("DB_NAME")
	DbPort := os.Getenv("DB_PORT")
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
