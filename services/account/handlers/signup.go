package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	je "gitlab.com/michaelk99/connectrn/internal/jsonerr"
	"gitlab.com/michaelk99/connectrn/internal/validator"
	"gitlab.com/michaelk99/connectrn/services/account"
)

const (
	// SignupErrCode error code
	SignupErrCode = "account.signup.error"
	// SignupExistsCode error code exists
	SignupExistsCode = "account.signup.exists"
	// SignupBadCredentialsCode bad creds
	SignupBadCredentialsCode = "account.signup.bad_credentials"
)

// SignUp sign up handler
func SignUp(v validator.Validator, s account.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// only support POST
		if r.Method != http.MethodPost {
			err := errors.New("method not supported [signup]")
			log.Printf(err.Error())
			resp := &je.Response{
				Code:    SignupErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, account.ServiceToHTTPErrorMap(err))
			return
		}

		var signUp account.SignupCredentials
		err := json.NewDecoder(r.Body).Decode(&signUp)
		if err != nil {
			log.Printf("%s: %v", SignupErrCode, err)
			resp := &je.Response{
				Code:    SignupErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, http.StatusBadRequest)
			return
		}

		// validate account credentials
		ok, fieldErrors := v.Struct(signUp)
		if !ok {
			resp := &je.Response{
				Code:       SignupBadCredentialsCode,
				Message:    SignupBadCredentialsCode,
				Additional: fieldErrors,
			}
			je.Error(r, w, resp, account.ServiceToHTTPErrorMap(err))
			return

		}

		a, err := s.SignUp(signUp)
		if err != nil {
			// do not track failed sign in attempts
			resp := &je.Response{
				Code:      SignupErrCode,
				Message:   err.Error(),
				SkipTrack: true,
			}
			je.Error(r, w, resp, account.ServiceToHTTPErrorMap(err))
			return
		}

		// return created account
		w.WriteHeader(http.StatusCreated) // must write status header before NewEcoder closes body
		err = json.NewEncoder(w).Encode(a)
		if err != nil {
			log.Printf("%s: %v", SignupErrCode, err)
			resp := &je.Response{
				Code:    SignupErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, http.StatusInternalServerError)
			return
		}

		log.Printf("successfully created account for email %s", signUp.Email)
		return
	}
}
