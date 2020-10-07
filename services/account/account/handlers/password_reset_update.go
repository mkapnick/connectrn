package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	je "gitlab.com/michaelk99/birrdi/api-soa/internal/jsonerr"
	"gitlab.com/michaelk99/birrdi/api-soa/services/account"
)

const (
	// UpdatePasswordResetErrCode code
	UpdatePasswordResetErrCode = "account.update_password.error"
	// UpdatePasswordResetBadDataCode bad data
	UpdatePasswordResetBadDataCode = "account.update_password_bad_data.error"
)

// UpdatePassword updates a pswd
func UpdatePassword(s account.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// only support POST
		if r.Method != http.MethodPost {
			err := errors.New("method not supported [update_password]")
			log.Printf(err.Error())
			resp := &je.Response{
				Code:    UpdatePasswordResetErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, account.ServiceToHTTPErrorMap(err))
			return
		}

		var rpr account.ResetPasswordRequest
		err := json.NewDecoder(r.Body).Decode(&rpr)
		if err != nil {
			log.Printf(UpdatePasswordResetErrCode)
			resp := &je.Response{
				Code:    UpdatePasswordResetBadDataCode,
				Message: "invalid body provided",
			}
			je.Error(r, w, resp, http.StatusBadRequest)
			return
		}

		err = s.UpdatePassword(&rpr)
		if err != nil {
			log.Printf("%s: %v", UpdatePasswordResetErrCode, err)
			resp := &je.Response{
				Code:    UpdatePasswordResetErrCode,
				Message: "failed to update account password",
				// no need to track a human error ;)
				SkipTrack: true,
			}
			je.Error(r, w, resp, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		log.Printf("successfully updated account password for reset token id %s", rpr.ID)
		return
	}
}
