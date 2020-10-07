package validator

import (
	"log"

	v9 "gopkg.in/go-playground/validator.v9"
)

// Simple wrapper around the official validator package to keep usage in HTTP
// handler simple. We swallow up any errors other then FieldErrors and log them
// out to stdout. We return a bool signifying if input was validated or not.

// FieldError is a useful struct for returning to clients. expresses the errors
// that were encountered with their json requests
type fieldError struct {
	// The field that failed validation
	Field string `json:"field"`
	// The validation tag that failed
	Reason string `json:"reason"`
}

// FieldErrors is a list of FieldError
type FieldErrors struct {
	Errors []fieldError `json:"validation_errors"`
}

// Validator A public interface for mocking and testing inside handlers (tho
// mocking is most likely unnecessary)
type Validator interface {
	Struct(v interface{}) (bool, *FieldErrors)
}

// private struct which simply wraps the official validator package
type validator struct {
	v *v9.Validate
}

// NewValidator A constructor for a Validator
func NewValidator(v *v9.Validate) Validator {
	return &validator{
		v: v,
	}
}

// Struct implemented our Validator interface. Swallows unlikely errors and
// returns a JSON serializable object for returning to clients
func (v *validator) Struct(V interface{}) (bool, *FieldErrors) {
	err := v.v.Struct(V)

	// type switch on returned error
	switch e := err.(type) {
	case *v9.InvalidValidationError:
		log.Printf("Validator: received type that cannot be handled: %v", e)
		return false, nil
	case v9.ValidationErrors:
		fieldErrs := make([]fieldError, 0)

		for _, v9FieldError := range e {
			fieldErr := fieldError{
				Field:  v9FieldError.Namespace(),
				Reason: v9FieldError.Tag(),
			}
			fieldErrs = append(fieldErrs, fieldErr)
		}

		return false, &FieldErrors{Errors: fieldErrs}
	}

	return true, nil
}
