// Package mapper converts between persistence entities and API models for the users domain.
package mapper

import "backend_crudgo/users/models"

// ToUserModel converts a persistence entity into the domain/API model.
// UserPassword is intentionally not copied so it never leaks through API responses.
func ToUserModel(db *models.UserDB) *models.User {
	return &models.User{
		UserID:             db.UserID,
		Name:               db.UserName,
		Email:              db.UserEmail,
		UserIdentifier:     db.UserIdentifier,
		UserTypeIdentifier: db.UserTypeIdentifier,
	}
}
