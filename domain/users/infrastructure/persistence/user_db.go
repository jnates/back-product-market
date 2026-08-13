package persistence

import "backend_crudgo/domain/users/domain/model"

// UserDB is the persistence entity for the public.users table.
// It carries only db tags; API-facing fields live on model.User.
type UserDB struct {
	UserID             int64  `db:"user_id"`
	UserName           string `db:"user_name"`
	UserIdentifier     int64  `db:"user_identifier"`
	UserEmail          string `db:"user_email"`
	UserPassword       string `db:"user_password"`
	UserTypeIdentifier int64  `db:"user_type_identifier"`
}

// ToModel converts the persistence entity into the domain/API model.
// UserPassword is intentionally not copied so it never leaks through API responses.
func (u *UserDB) ToModel() *model.User {
	return &model.User{
		UserID:             u.UserID,
		Name:               u.UserName,
		Email:              u.UserEmail,
		UserIdentifier:     u.UserIdentifier,
		UserTypeIdentifier: u.UserTypeIdentifier,
	}
}
