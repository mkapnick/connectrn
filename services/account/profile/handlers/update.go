package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	je "gitlab.com/michaelk99/birrdi/api-soa/internal/jsonerr"
	"gitlab.com/michaelk99/birrdi/api-soa/internal/token"
	"gitlab.com/michaelk99/birrdi/api-soa/internal/validator"
	"gitlab.com/michaelk99/birrdi/api-soa/services/profile"
)

const (
	// UpdateErrCode code
	UpdateErrCode = "profile.update.error"
	// UpdateExistsCode code
	UpdateExistsCode = "profile.update.exists"
	// UpdateBadDataCode code
	UpdateBadDataCode = "profile.update.bad_data"
)

// Update checks email against password and assigns a token if valid
func Update(v validator.Validator, ps profile.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// only support PUT
		if r.Method != http.MethodPut {
			err := errors.New("method not supported [profile update]")
			log.Printf(err.Error())
			resp := &je.Response{
				Code:    UpdateErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, profile.ServiceToHTTPErrorMap(err))
			return
		}

		session := r.Context().Value("Session").(*token.Session)

		if session.ProfileID == "" {
			err := errors.New("invalid id")
			log.Printf(err.Error())
			resp := &je.Response{
				Code:    UpdateErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, profile.ServiceToHTTPErrorMap(err))
			return
		}

		var prof profile.Profile
		err := json.NewDecoder(r.Body).Decode(&prof)

		if err != nil {
			log.Printf("%s: %v", UpdateErrCode, err)
			resp := &je.Response{
				Code:    UpdateErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, http.StatusBadRequest)
			return
		}

		// validate profile
		ok, fieldErrors := v.Struct(prof)
		if !ok {
			resp := &je.Response{
				Code:       UpdateBadDataCode,
				Message:    UpdateBadDataCode,
				Additional: fieldErrors,
			}
			je.Error(r, w, resp, http.StatusBadRequest)
			return
		}

		// override the prof ID from the session
		prof.ID = session.ProfileID
		// override the account ID from the session
		prof.AccountID = session.AccountID
		// run the update: Note, the entire profile must be sent in the
		// PUT request otherwise zero value fields will override
		dbProfile, err := ps.Update(prof)

		if err != nil {
			resp := &je.Response{
				Code:    UpdateErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, profile.ServiceToHTTPErrorMap(err))
			return
		}

		// return created profile
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(dbProfile)
		if err != nil {
			log.Printf("%s: %v", UpdateErrCode, err)
			resp := &je.Response{
				Code:    UpdateErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, http.StatusBadRequest)
			return
		}

		log.Printf("successfully updated profile id %s", session.ProfileID)
		return
	}
}
