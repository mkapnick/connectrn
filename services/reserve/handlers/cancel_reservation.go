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
	// CancelTableErrCode code
	CancelTableErrCode = "reservation.table.reserve.error"
	// CancelBadDataCode code
	CancelBadDataCode = "reservation.table.reserve.bad_data"
)

// CancelReservation cancels a reservation
func CancelReservation(v validator.Validator, rs reserve.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		restaurantID := vars["restaurant_id"]
		tableID := vars["table_id"]
		userReservationID := vars["user_reservation_id"]

		// parse tee time reservation data
		var req reserve.CancelReserveRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Printf("%s: %v", CancelTableErrCode, err)
			resp := &je.Response{
				Code:    CancelTableErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, http.StatusBadRequest)
			return
		}

		// validate tee time reservation
		ok, fieldErrors := v.Struct(req)
		if !ok {
			log.Printf("%s: %s", CancelTableErrCode, "validator.v9 failed")
			resp := &je.Response{
				Code:       CancelBadDataCode,
				Message:    CancelBadDataCode,
				Additional: fieldErrors,
			}
			je.Error(r, w, resp, http.StatusBadRequest)
			return
		}

		// override fields
		req.RestaurantID = restaurantID
		req.TableID = tableID
		req.UserReservationID = userReservationID

		// add the `profile_id` to the request
		session := r.Context().Value("Session").(*token.Session)
		req.ProfileID = session.ProfileID

		// this `user` is an owner. Make sure they can reserve at this
		// restaurant.
		if session.RestaurantID != "" {
			// TODO check to make sure this owner can make a reservation at
			// this restaurant.
		}

		ur, err := rs.CancelReservation(req)
		if err != nil {
			log.Printf("%s: %s", CancelTableErrCode, err)
			resp := &je.Response{
				Code:    CancelTableErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, reserve.ServiceToHTTPErrorMap(err))
			return
		}

		// return created event
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(ur)
		if err != nil {
			log.Printf("%s: %v", CancelTableErrCode, err)
			resp := &je.Response{
				Code:    CancelTableErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, http.StatusBadRequest)
			return
		}

		log.Printf("successfully cancelled reservation %s", req.UserReservationID)
		return
	}
}
