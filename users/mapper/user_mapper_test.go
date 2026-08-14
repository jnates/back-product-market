package mapper_test

import (
	"testing"

	"backend_crudgo/users/mapper"
	"backend_crudgo/users/models"

	"github.com/stretchr/testify/assert"
)

func TestToUserModel(t *testing.T) {
	dbUser := &models.UserDB{
		UserID:             1,
		UserName:           "Test User",
		UserIdentifier:     1,
		UserEmail:          "test@example.com",
		UserPassword:       "hashed",
		UserTypeIdentifier: 1,
	}

	user := mapper.ToUserModel(dbUser)

	assert.Equal(t, dbUser.UserID, user.UserID)
	assert.Equal(t, dbUser.UserName, user.Name)
	assert.Equal(t, dbUser.UserEmail, user.Email)
	assert.Equal(t, dbUser.UserIdentifier, user.UserIdentifier)
	assert.Equal(t, dbUser.UserTypeIdentifier, user.UserTypeIdentifier)
	assert.Empty(t, user.UserPassword)
}
