package handlers

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"log"
	"net/http"

	je "gitlab.com/michaelk99/connectrn/internal/jsonerr"
	"gitlab.com/michaelk99/connectrn/services/restaurant"
)

const (
	// FetchErrCode code
	FetchErrCode = "restaurant.fetch.error"
	// FetchExistsCode code
	FetchExistsCode = "restaurant.fetch.exists"
)

// Fetch checks email against password and assigns a token if valid
func Fetch(rs restaurant.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// only support GET
		if r.Method != http.MethodGet {
			err := errors.New("method not supported [restaurant fetch]")
			log.Printf(err.Error())
			resp := &je.Response{
				Code:    FetchErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, restaurant.ServiceToHTTPErrorMap(err))
			return
		}

		vars := mux.Vars(r)
		restaurantID := vars["restauraunt_id"]

		rse, err := rs.FetchRestaurant(restaurantID)
		if err != nil {
			resp := &je.Response{
				Code:    FetchErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, restaurant.ServiceToHTTPErrorMap(err))
			return
		}

		// return created restaurant
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(rse)
		if err != nil {
			log.Printf("%s: %v", FetchErrCode, err)
			resp := &je.Response{
				Code:    FetchErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, http.StatusBadRequest)
			return
		}

		log.Printf("successfully fetched restaurant id %s", restaurantID)
		return
	}
}
