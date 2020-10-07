package handlers

import (
	"errors"
	"log"
	"net/http"

	je "gitlab.com/michaelk99/birrdi/api-soa/internal/jsonerr"
	"gitlab.com/michaelk99/birrdi/api-soa/services/account"
)

// CRUD forwards request based on http method
func CRUD(s account.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			Fetch(s).ServeHTTP(w, r)
			return
		case http.MethodDelete:
			Delete(s).ServeHTTP(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			err := errors.New("method not allowed")
			log.Printf(err.Error())
			resp := &je.Response{
				Code:    err.Error(),
				Message: "method not allowed",
			}
			je.Error(r, w, resp, account.ServiceToHTTPErrorMap(err))
			return
		}
	}
}
