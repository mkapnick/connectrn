package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	je "gitlab.com/michaelk99/birrdi/api-soa/internal/jsonerr"
	"gitlab.com/michaelk99/birrdi/api-soa/internal/token"
	"gitlab.com/michaelk99/birrdi/api-soa/services/account"
)

const (
	// UpdateErrCode code
	UpdateErrCode = "account.update.error"
	// UpdateErrAccountNotFoundCode code
	UpdateErrAccountNotFoundCode = "account.update.notfound"
	// UpdateErrInvalidIDCode code
	UpdateErrInvalidIDCode = "account.update.invalid_id"
	// UpdateMethodNotSupportedCode code
	UpdateMethodNotSupportedCode = "account.update.method_not_supported"
	// UpdatePath is the path we expect the request specify
	UpdatePath = "/api/v1/account/"
)

// Update receives an account and updates the specified public feilds. Update will not be used
// for password reset.
func Update(s account.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// only support PUT
		if r.Method != http.MethodPut {
			err := errors.New("method not supported [signup update]")
			log.Printf(err.Error())
			resp := &je.Response{
				Code:    UpdateErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, account.ServiceToHTTPErrorMap(err))
			return
		}

		session := r.Context().Value("Session").(*token.Session)

		var acc account.Account
		err := json.NewDecoder(r.Body).Decode(&acc)
		if err != nil {
			log.Printf(UpdateErrCode)
			resp := &je.Response{
				Code:    UpdateErrCode,
				Message: "bad request provided",
			}
			je.Error(r, w, resp, http.StatusBadRequest)
			return
		}

		// override the ID
		acc.ID = session.AccountID

		// update account with details
		_, err = s.Update(acc)
		if err != nil {
			log.Printf("%s: %v", UpdateErrCode, err)
			resp := &je.Response{
				Code:    UpdateErrCode,
				Message: "failed to update account",
			}
			je.Error(r, w, resp, http.StatusInternalServerError)
			return
		}

		// return updated account
		err = json.NewEncoder(w).Encode(acc)
		if err != nil {
			log.Printf("%s: %v", UpdateErrCode, err)
			resp := &je.Response{
				Code:    UpdateErrCode,
				Message: "failed to update account",
			}
			je.Error(r, w, resp, http.StatusInternalServerError)
			return
		}

		return
	}
}
