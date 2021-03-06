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
	// FetchAllErrCode code
	FetchAllErrCode = "restaurant.tables.fetchall.error"
	// FetchAllBadDataCode code
	FetchAllBadDataCode = "restaurant.tables.fetchall.bad_data"
)

// FetchAllTables fetch all restaurant tables
func FetchAllTables(rs restaurant.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		restaurantID := vars["restaurant_id"]

		date := r.URL.Query().Get("date")
		if date == "" {
			err := errors.New("date is a required field")
			log.Printf(err.Error())
			resp := &je.Response{
				Code:    FetchAllErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, restaurant.ServiceToHTTPErrorMap(err))
			return
		}

		ts, err := rs.FetchAllTables(restaurantID, date)
		if err != nil {
			resp := &je.Response{
				Code:    FetchAllErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, restaurant.ServiceToHTTPErrorMap(err))
			return
		}

		// return created restaurant.tables
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(ts)
		if err != nil {
			log.Printf("%s: %v", FetchAllErrCode, err)
			resp := &je.Response{
				Code:    FetchAllErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, http.StatusBadRequest)
			return
		}

		log.Printf("successfully fetched all restaurant tables on %s", date)
		return
	}
}
