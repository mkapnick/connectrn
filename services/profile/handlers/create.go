package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	je "gitlab.com/michaelk99/connectrn/internal/jsonerr"
	"gitlab.com/michaelk99/connectrn/internal/token"
	"gitlab.com/michaelk99/connectrn/internal/validator"
	"gitlab.com/michaelk99/connectrn/services/profile"
)

const (
	// CreateErrCode error code
	CreateErrCode = "profile.create.error"
	// CreateExistsCode error code exists
	CreateExistsCode = "profile.create.exists"
	// CreateBadDataCode bad data
	CreateBadDataCode = "profile.create.bad_data"
)

// Create sign up handler
func Create(v validator.Validator, ps profile.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// only support POST
		if r.Method != http.MethodPost {
			err := errors.New("method not supported [profile create]")
			log.Printf(err.Error())
			resp := &je.Response{
				Code:    CreateErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, profile.ServiceToHTTPErrorMap(err))
			return
		}

		var prof profile.Profile
		err := json.NewDecoder(r.Body).Decode(&prof)
		if err != nil {
			log.Printf("%s: %v", CreateErrCode, err)
			resp := &je.Response{
				Code:    CreateErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, http.StatusBadRequest)
			return
		}

		// validate profile
		ok, fieldErrors := v.Struct(prof)
		if !ok {
			resp := &je.Response{
				Code:       CreateBadDataCode,
				Message:    CreateBadDataCode,
				Additional: fieldErrors,
			}
			je.Error(r, w, resp, http.StatusBadRequest)
			return
		}

		// override the account ID with the session accountID
		session := r.Context().Value("Session").(*token.Session)
		prof.AccountID = session.AccountID

		p, err := ps.Create(prof)
		if err != nil {
			resp := &je.Response{
				Code:    CreateErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, profile.ServiceToHTTPErrorMap(err))
			return
		}

		// return created profile
		w.WriteHeader(http.StatusCreated) // must write status header before NewEcoder closes body
		err = json.NewEncoder(w).Encode(p)
		if err != nil {
			log.Printf("%s: %v", CreateErrCode, err)
			resp := &je.Response{
				Code:    CreateErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, http.StatusInternalServerError)
			return
		}

		log.Printf("successfully created profile for ID %s", prof.ID)
		return
	}
}
