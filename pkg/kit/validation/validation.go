// Package validation wires go-playground/validator into Echo via the
// echo.Validator interface, per GUIDE-HTTP-009.
package validation

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

// EchoValidator adapts *validator.Validate to the echo.Validator interface.
type EchoValidator struct {
	validate *validator.Validate
}

// New builds an EchoValidator ready to be assigned to echo.Echo.Validator.
func New() *EchoValidator {
	return &EchoValidator{validate: validator.New()}
}

// Validate implements echo.Validator.
func (v *EchoValidator) Validate(i interface{}) error {
	return v.validate.Struct(i)
}

// FormatErrors turns a validation error into a field -> failed-rule map
// suitable for the API's error response payload.
func FormatErrors(err error) map[string]string {
	var fieldErrors validator.ValidationErrors
	if !errors.As(err, &fieldErrors) {
		return map[string]string{"error": err.Error()}
	}

	out := make(map[string]string, len(fieldErrors))
	for _, fe := range fieldErrors {
		out[fe.Field()] = fe.ActualTag()
	}

	return out
}
