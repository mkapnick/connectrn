package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"

	je "gitlab.com/michaelk99/connectrn/internal/jsonerr"
	"gitlab.com/michaelk99/connectrn/services/restaurant"
)

const (
	// FetchTableErrCode code
	FetchTableErrCode = "restaurant.fetch.error"
	// FetchTableExistsCode code
	FetchTableExistsCode = "restaurant.fetch.exists"
)

// FetchTable checks email against password and assigns a token if valid
func FetchTable(rs restaurant.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		restaurantID := vars["restaurant_id"]
		tableID := vars["table_id"]

		t, err := rs.FetchTable(restaurantID, tableID)
		if err != nil {
			resp := &je.Response{
				Code:    FetchTableErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, restaurant.ServiceToHTTPErrorMap(err))
			return
		}

		// return created restaurant
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(t)
		if err != nil {
			log.Printf("%s: %v", FetchTableErrCode, err)
			resp := &je.Response{
				Code:    FetchTableErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, http.StatusBadRequest)
			return
		}

		log.Printf("successfully fetched table id %s", tableID)
		return
	}
}
