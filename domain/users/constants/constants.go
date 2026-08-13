// Package constants centralizes literal values used across the users domain
// (table and column names) so they are defined once and referenced everywhere.
package constants

const (
	// UsersTable is the public.users table name.
	UsersTable = "users"

	// ColumnUserID is the users table primary key column.
	ColumnUserID = "user_id"
	// ColumnUserName is the username column, used for login lookup.
	ColumnUserName = "user_name"
	// ColumnUserIdentifier is the user's identification number column.
	ColumnUserIdentifier = "user_identifier"
	// ColumnUserEmail is the user email column.
	ColumnUserEmail = "user_email"
	// ColumnUserPassword is the hashed password column.
	ColumnUserPassword = "user_password"
	// ColumnUserTypeIdentifier is the identification-type foreign key column.
	ColumnUserTypeIdentifier = "user_type_identifier"
)

// Columns lists every users column returned by SELECT queries, in scan order.
var Columns = []string{
	ColumnUserID,
	ColumnUserName,
	ColumnUserIdentifier,
	ColumnUserEmail,
	ColumnUserPassword,
	ColumnUserTypeIdentifier,
}
