package reserve

import (
	"fmt"
	"net/http"
)

// ErrReserveNotFound  not found
type ErrReserveNotFound struct {
	msg error
}

func (e ErrReserveNotFound) Error() string {
	if len(e.msg.Error()) != 0 {
		return fmt.Sprintf(e.msg.Error())
	}
	return "reserve not found"
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

// ServiceToHTTPErrorMap maps the s service's errors to http
func ServiceToHTTPErrorMap(err error) (code int) {
	switch err.(type) {
	case ErrReserveNotFound:
		return http.StatusNotFound
	case ErrInternal:
		return http.StatusBadRequest
	default:
		return http.StatusBadRequest
	}
}
