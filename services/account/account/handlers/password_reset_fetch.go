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
	// FetchPasswordResetTokenErrCode code
	FetchPasswordResetTokenErrCode = "account.fetch_reset_token.error"
	// FetchPasswordResetTokenExistsCode code
	FetchPasswordResetTokenExistsCode = "account.fetch_reset_token.exists"
)

// GetResetToken checks email against password and assigns a token if valid
func GetResetToken(s account.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// only support GET
		if r.Method != http.MethodGet {
			err := errors.New("method not supported [fetch_reset_token]")
			log.Printf(err.Error())
			resp := &je.Response{
				Code:    FetchPasswordResetTokenErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, account.ServiceToHTTPErrorMap(err))
			return
		}

		token := r.URL.Query().Get("token")

		if token == "" {
			err := errors.New("token not found [fetch_reset_token]")
			log.Printf(err.Error())
			resp := &je.Response{
				Code:    FetchPasswordResetTokenErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, account.ServiceToHTTPErrorMap(err))
			return
		}

		prt, err := s.FetchPasswordResetToken(token)

		if err != nil {
			resp := &je.Response{
				Code:    FetchPasswordResetTokenErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, account.ServiceToHTTPErrorMap(err))
			return
		}

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(prt)
		if err != nil {
			log.Printf("%s: %v", FetchPasswordResetTokenErrCode, err)
			resp := &je.Response{
				Code:    FetchPasswordResetTokenErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, http.StatusBadRequest)
			return
		}

		log.Printf("successfully fetched reset token %s", token)
		return
	}
}
