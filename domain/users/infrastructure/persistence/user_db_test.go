package persistence

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserDBToModel(t *testing.T) {
	dbUser := &UserDB{
		UserID:             1,
		UserName:           "Test User",
		UserIdentifier:     1,
		UserEmail:          "test@example.com",
		UserPassword:       "hashed",
		UserTypeIdentifier: 1,
	}

	user := dbUser.ToModel()

	assert.Equal(t, dbUser.UserID, user.UserID)
	assert.Equal(t, dbUser.UserName, user.Name)
	assert.Equal(t, dbUser.UserEmail, user.Email)
	assert.Equal(t, dbUser.UserIdentifier, user.UserIdentifier)
	assert.Equal(t, dbUser.UserTypeIdentifier, user.UserTypeIdentifier)
	assert.Empty(t, user.UserPassword)
}
