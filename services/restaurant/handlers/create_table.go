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
	// CreateTableErrCode error code
	CreateTableErrCode = "restaurant.table.create.error"
	// CreateTableExistsCode error code exists
	CreateTableExistsCode = "restaurant.table.create.exists"
	// CreateTableBadDataCode bad data
	CreateTableBadDataCode = "restaurant.table.create.bad_data"
)

// CreateTable sign up handler
func CreateTable(v validator.Validator, rs restaurant.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// only support POST
		if r.Method != http.MethodPost {
			err := errors.New("method not supported [restaurant create]")
			log.Printf(err.Error())
			resp := &je.Response{
				Code:    CreateTableErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, restaurant.ServiceToHTTPErrorMap(err))
			return
		}

		var cr restaurant.TableCreateRequest
		err := json.NewDecoder(r.Body).Decode(&cr)
		if err != nil {
			log.Printf("%s: %v", CreateTableErrCode, err)
			resp := &je.Response{
				Code:    CreateTableErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, http.StatusBadRequest)
			return
		}

		// validate restaurant
		ok, fieldErrors := v.Struct(cr)
		if !ok {
			resp := &je.Response{
				Code:       CreateTableBadDataCode,
				Message:    CreateTableBadDataCode,
				Additional: fieldErrors,
			}
			je.Error(r, w, resp, http.StatusBadRequest)
			return
		}

		t, err := rs.CreateTable(cr)
		if err != nil {
			resp := &je.Response{
				Code:    CreateTableErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, restaurant.ServiceToHTTPErrorMap(err))
			return
		}

		// return created restaurant
		w.WriteHeader(http.StatusCreated) // must write status header before NewEcoder closes body
		err = json.NewEncoder(w).Encode(t)
		if err != nil {
			log.Printf("%s: %v", CreateTableErrCode, err)
			resp := &je.Response{
				Code:    CreateTableErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, http.StatusInternalServerError)
			return
		}

		log.Printf("successfully created restaurant table for ID %s", t.ID)
		return
	}
}
