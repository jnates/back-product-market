// Package apperrors defines the sentinel business errors shared across domains.
package apperrors

import "errors"

var (
	// ErrNotFound indicates the requested resource does not exist.
	ErrNotFound = errors.New("resource not found")

	// ErrInvalidInput indicates the request body or parameters are invalid.
	ErrInvalidInput = errors.New("invalid input")

	// ErrUnauthorized indicates the request lacks valid authentication.
	ErrUnauthorized = errors.New("unauthorized")

	// ErrConflict indicates the resource already exists or violates a uniqueness constraint.
	ErrConflict = errors.New("resource already exists")
)
