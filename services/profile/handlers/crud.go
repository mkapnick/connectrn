package handlers

import (
	"net/http"

	"gitlab.com/michaelk99/connectrn/internal/validator"
	"gitlab.com/michaelk99/connectrn/services/profile"
)

// CRUD forwards request based on http method
func CRUD(v validator.Validator, ps profile.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case http.MethodGet:
			Fetch(ps).ServeHTTP(w, r)
			return
		case http.MethodPost:
			Create(v, ps).ServeHTTP(w, r)
			return
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			http.Error(w, "method not supported", http.StatusMethodNotAllowed)
			return
		}
	}
}
