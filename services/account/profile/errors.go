package profile

import (
	"net/http"
)

// ErrProfileCreate is returned when profile cannot be created
type ErrProfileCreate struct {
	msg error
}

func (e ErrProfileCreate) Error() string {
	return "could not create profile"
}

// ErrProfileNotFound profile not found
type ErrProfileNotFound struct {
	msg error
}

func (e ErrProfileNotFound) Error() string {
	return "profile not found"
}

// ErrProfileExists profile exists in db
type ErrProfileExists struct {
	msg error
}

func (e ErrProfileExists) Error() string {
	return "profile exists"
}

// ServiceToHTTPErrorMap maps the profiles service's errors to http
func ServiceToHTTPErrorMap(err error) (code int) {
	switch err.(type) {
	case ErrProfileCreate:
		return http.StatusConflict
	case ErrProfileNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
