package account

import (
	"fmt"
	"net/http"
)

// ErrUserExists is returned when creating a user fails due to the user already
// existing
type ErrUserExists struct{}

func (e ErrUserExists) Error() string {
	return "user exists"
}

// ErrUserNotFound is returned when creating a user fails due to the user
// already existing
type ErrUserNotFound struct{}

func (e ErrUserNotFound) Error() string {
	return "user not found"
}

// ErrUpdateFail is returned when an error occured updating an account's detail
type ErrUpdateFail struct {
	msg error
}

func (e ErrUpdateFail) Error() string {
	return fmt.Sprintf("%v", e.msg)
}

// ErrDeleteFail is returned when delete of an account fails
type ErrDeleteFail struct {
	msg error
}

func (e ErrDeleteFail) Error() string {
	return fmt.Sprintf("%v", e.msg)
}

// ErrPasswordHash is returned to indicate issue hashing an account's password
type ErrPasswordHash struct {
	msg string
}

func (e ErrPasswordHash) Error() string {
	return fmt.Sprintf(e.msg)
}

// ErrInvalidLogin is returned when invalid account credentials have been used
type ErrInvalidLogin struct{}

func (e ErrInvalidLogin) Error() string {
	return "username or password incorrect"
}

// ErrCreateToken is returned when a token fails to be created
type ErrCreateToken struct {
	msg error
}

func (e ErrCreateToken) Error() string {
	return fmt.Sprintf("%v", e.msg)
}

// ErrInvalidIDType is returned when an unknown ID type in an IDQuery is found.
// See queries.go for appropriate enum constants
type ErrInvalidIDType struct{}

func (e ErrInvalidIDType) Error() string {
	return "unknown ID type"
}

// ErrInternal is a catch all error for internal issues
type ErrInternal struct {
	msg string
}

func (e ErrInternal) Error() string {
	return fmt.Sprintf(e.msg)
}

// ErrProf is an error for issues contacting the Profile service
type ErrProf struct {
	msg string
}

func (e ErrProf) Error() string {
	return fmt.Sprintf(e.msg)
}

// ErrStripe is an error for issues contacting the Stripe service
type ErrStripe struct {
	msg error
}

func (e ErrStripe) Error() string {
	if len(e.msg.Error()) != 0 {
		return fmt.Sprintf("Error with stripe: %s", e.msg.Error())
	}
	return "Error with stripe"
}

// ServiceToHTTPErrorMap maps the account service's errors to http
func ServiceToHTTPErrorMap(err error) (code int) {
	switch err.(type) {
	case ErrUserExists:
		return http.StatusConflict
	case ErrUserNotFound:
		return http.StatusBadRequest
	case ErrPasswordHash:
		return http.StatusInternalServerError
	case ErrInvalidLogin:
		return http.StatusUnauthorized
	case ErrCreateToken:
		return http.StatusInternalServerError
	case ErrInternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
