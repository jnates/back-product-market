// Package model defines the user domain entities exposed over the API.
package model

// User represents a user exposed through the API (request and response body).
// UserPassword is only populated when binding a create/login request; it is never
// populated by the repository when returning users, so it never leaks in responses.
type User struct {
	UserID             int64  `json:"user_id"`
	Name               string `json:"user_name"`
	Email              string `json:"user_email"`
	UserIdentifier     int64  `json:"user_identifier"`
	UserPassword       string `json:"user_password,omitempty"`
	UserTypeIdentifier int64  `json:"user_type_identifier"`
}

// LoginResponse carries the JWT issued after a successful login.
type LoginResponse struct {
	Token string `json:"token"`
}
