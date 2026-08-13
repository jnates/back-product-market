package database_test

import (
	"testing"

	"backend_crudgo/infrastructure/database"
	"backend_crudgo/infrastructure/kit/enum"

	"github.com/stretchr/testify/assert"
)

func TestNewMissingConfig(t *testing.T) {
	// Force an incomplete configuration so Validate() fails before any network I/O is attempted.
	t.Setenv(enum.DBHost, "")
	t.Setenv(enum.DBPort, "")
	t.Setenv(enum.DBName, "")
	t.Setenv(enum.DBUser, "")

	pool, err := database.New()

	assert.Error(t, err)
	assert.Nil(t, pool)
}
