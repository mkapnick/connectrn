package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	je "gitlab.com/michaelk99/connectrn/internal/jsonerr"
	"gitlab.com/michaelk99/connectrn/internal/token"
	"gitlab.com/michaelk99/connectrn/internal/validator"
	"gitlab.com/michaelk99/connectrn/services/reserve"
	"log"
	"net/http"
)

const (
	// ReserveTablesErrCode code
	ReserveTablesErrCode = "reservation.tables.reserve.error"
	// ReserveTablesBadDataCode code
	ReserveTablesBadDataCode = "reservation.tables.reserve.bad_data"
)

// ReserveTables reserves a table
func ReserveTables(v validator.Validator, rs reserve.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		restaurantID := vars["restaurant_id"]

		// parse tee time reservation data
		var reqs []*reserve.ReserveRequest
		err := json.NewDecoder(r.Body).Decode(&reqs)
		if err != nil {
			log.Printf("%s: %v", ReserveTablesErrCode, err)
			resp := &je.Response{
				Code:    ReserveTablesErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, http.StatusBadRequest)
			return
		}

		// add the `profile_id` to the request
		session := r.Context().Value("Session").(*token.Session)

		// validate and override fields
		for _, req := range reqs {
			// validate tee time reservation
			ok, fieldErrors := v.Struct(req)
			if !ok {
				log.Printf("%s: %s", ReserveTablesErrCode, "validator.v9 failed")
				resp := &je.Response{
					Code:       ReserveBadDataCode,
					Message:    ReserveBadDataCode,
					Additional: fieldErrors,
				}
				je.Error(r, w, resp, http.StatusBadRequest)
				return
			}
			req.RestaurantID = restaurantID
			req.ProfileID = session.ProfileID
		}

		// this `user` is an owner. Make sure they can reserve at this
		// restaurant.
		if session.RestaurantID != "" {
			// TODO check to make sure this owner can make a reservation at
			// this restaurant.
		}

		ur, err := rs.ReserveTables(reqs)
		if err != nil {
			log.Printf("%s: %s", ReserveTablesErrCode, err)
			resp := &je.Response{
				Code:    ReserveTablesErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, reserve.ServiceToHTTPErrorMap(err))
			return
		}

		// return created event
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(ur)
		if err != nil {
			log.Printf("%s: %v", ReserveTablesErrCode, err)
			resp := &je.Response{
				Code:    ReserveTablesErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, http.StatusBadRequest)
			return
		}

		log.Println("successfully reserved tables")
		return
	}
}
