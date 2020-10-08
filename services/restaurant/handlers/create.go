package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	je "gitlab.com/michaelk99/connectrn/internal/jsonerr"
	"gitlab.com/michaelk99/connectrn/internal/validator"
	"gitlab.com/michaelk99/connectrn/services/restaurant"
)

const (
	// CreateErrCode error code
	CreateErrCode = "restaurant.create.error"
	// CreateExistsCode error code exists
	CreateExistsCode = "restaurant.create.exists"
	// CreateBadDataCode bad data
	CreateBadDataCode = "restaurant.create.bad_data"
)

// Create sign up handler
func Create(v validator.Validator, rs restaurant.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// only support POST
		if r.Method != http.MethodPost {
			err := errors.New("method not supported [restaurant create]")
			log.Printf(err.Error())
			resp := &je.Response{
				Code:    CreateErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, restaurant.ServiceToHTTPErrorMap(err))
			return
		}

		var rest restaurant.RestaurantCreateRequest
		err := json.NewDecoder(r.Body).Decode(&rest)
		if err != nil {
			log.Printf("%s: %v", CreateErrCode, err)
			resp := &je.Response{
				Code:    CreateErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, http.StatusBadRequest)
			return
		}

		// validate restaurant
		ok, fieldErrors := v.Struct(rest)
		if !ok {
			resp := &je.Response{
				Code:       CreateBadDataCode,
				Message:    CreateBadDataCode,
				Additional: fieldErrors,
			}
			je.Error(r, w, resp, http.StatusBadRequest)
			return
		}

		rr, err := rs.CreateRestaurant(rest)
		if err != nil {
			resp := &je.Response{
				Code:    CreateErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, restaurant.ServiceToHTTPErrorMap(err))
			return
		}

		// return created restaurant
		w.WriteHeader(http.StatusCreated) // must write status header before NewEcoder closes body
		err = json.NewEncoder(w).Encode(rr)
		if err != nil {
			log.Printf("%s: %v", CreateErrCode, err)
			resp := &je.Response{
				Code:    CreateErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, http.StatusInternalServerError)
			return
		}

		log.Printf("successfully created restaurant for ID %s", rr.ID)
		return
	}
}
