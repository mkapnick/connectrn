package handlers

import (
	"net/http"

	"gitlab.com/michaelk99/birrdi/api-soa/internal/validator"
	"gitlab.com/michaelk99/birrdi/api-soa/services/profile"
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
		case http.MethodPut:
			Update(v, ps).ServeHTTP(w, r)
			return
		case http.MethodDelete:
			Delete(ps).ServeHTTP(w, r)
			return
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			http.Error(w, "method not supported", http.StatusMethodNotAllowed)
			return
		}
	}
}
