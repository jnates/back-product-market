package storage_test

import (
	"testing"

	"backend_crudgo/configs"
	"backend_crudgo/configs/storage"

	"github.com/stretchr/testify/assert"
)

func TestPostgresConnectionMissingConfig(t *testing.T) {
	// An incomplete configuration must fail Validate() before any network I/O is attempted.
	pool, err := storage.PostgresConnection(&configs.Config{})

	assert.Error(t, err)
	assert.Nil(t, pool)
}
