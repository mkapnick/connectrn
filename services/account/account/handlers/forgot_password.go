package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	je "gitlab.com/michaelk99/birrdi/api-soa/internal/jsonerr"
	"gitlab.com/michaelk99/birrdi/api-soa/internal/validator"
	"gitlab.com/michaelk99/birrdi/api-soa/services/account"
)

const (
	// ForgotPasswordErrCode error code
	ForgotPasswordErrCode = "account.forgot_password.error"
	// ForgotPasswordExistsCode error code exists
	ForgotPasswordExistsCode = "account.forgot_password.exists"
	// ForgotPasswordBadCode bad creds
	ForgotPasswordBadCode = "account.forgot_password.bad_credentials"
)

// ForgotPassword sign up handler
func ForgotPassword(v validator.Validator, s account.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// only support POST
		if r.Method != http.MethodPost {
			err := errors.New("method not supported [forgot_password]")
			log.Printf(err.Error())
			resp := &je.Response{
				Code:    ForgotPasswordErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, account.ServiceToHTTPErrorMap(err))
			return
		}

		var req account.ForgotPasswordRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Printf("%s: %v", ForgotPasswordErrCode, err)
			resp := &je.Response{
				Code:    ForgotPasswordErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, http.StatusBadRequest)
			return
		}

		// validate account credentials
		ok, fieldErrors := v.Struct(req)
		if !ok {
			resp := &je.Response{
				Code:       ForgotPasswordBadCode,
				Message:    ForgotPasswordBadCode,
				Additional: fieldErrors,
			}
			je.Error(r, w, resp, account.ServiceToHTTPErrorMap(err))
			return

		}

		err = s.CreatePasswordResetToken(&req)
		if err != nil {
			// no need to log this in sentry because could be a human error
			// and gives us false positive error reporting
			// see https://sentry.io/organizations/birrdi/issues/1629398365/?project=4953326&query=is%3Aunresolved
			resp := &je.Response{
				Code:      ForgotPasswordErrCode,
				Message:   err.Error(),
				SkipTrack: true,
			}
			je.Error(r, w, resp, account.ServiceToHTTPErrorMap(err))
			return
		}

		// return status 201
		w.WriteHeader(http.StatusCreated)

		log.Printf("successfully created password reset token for email %s", req.Email)
		return
	}
}
