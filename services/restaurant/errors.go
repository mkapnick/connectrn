package restaurant

import (
	"fmt"
	"net/http"
)

// ErrRestaurantNotFound golf course not found
type ErrRestaurantNotFound struct {
	msg error
}

func (e ErrRestaurantNotFound) Error() string {
	if len(e.msg.Error()) != 0 {
		return fmt.Sprintf(e.msg.Error())
	}
	return "restaurant not found"
}

// ErrInternal internal error
type ErrInternal struct {
	msg error
}

func (e ErrInternal) Error() string {
	if len(e.msg.Error()) != 0 {
		return fmt.Sprintf(e.msg.Error())
	}
	return "internal error"
}

// ServiceToHTTPErrorMap maps the golf courses service's errors to http
func ServiceToHTTPErrorMap(err error) (code int) {
	switch err.(type) {
	case ErrRestaurantNotFound:
		return http.StatusNotFound
	case ErrInternal:
		return http.StatusBadRequest
	default:
		return http.StatusBadRequest
	}
}
