package models

// UserDB is the persistence entity for the public.users table.
// It carries only db tags; API-facing fields live on User.
type UserDB struct {
	UserID             int64  `db:"user_id"`
	UserName           string `db:"user_name"`
	UserIdentifier     int64  `db:"user_identifier"`
	UserEmail          string `db:"user_email"`
	UserPassword       string `db:"user_password"`
	UserTypeIdentifier int64  `db:"user_type_identifier"`
}
