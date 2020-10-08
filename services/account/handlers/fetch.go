package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"gitlab.com/michaelk99/connectrn/internal/token"

	je "gitlab.com/michaelk99/connectrn/internal/jsonerr"
	"gitlab.com/michaelk99/connectrn/services/account"
)

const (
	// FetchErrCode code
	FetchErrCode = "account.fetch.error"
	// FetchExistsCode code
	FetchExistsCode = "account.fetch.exists"
)

// Fetch checks email against password and assigns a token if valid
func Fetch(s account.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := r.Context().Value("Session").(*token.Session)

		idQuery := account.IDQuery{
			Type:  account.ID,
			Value: session.AccountID,
		}

		acc, err := s.Fetch(idQuery)
		if err != nil {
			resp := &je.Response{
				Code:    FetchErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, account.ServiceToHTTPErrorMap(err))
			return
		}

		// return fetched account
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(acc)
		if err != nil {
			log.Printf("%s: %v", FetchErrCode, err)
			resp := &je.Response{
				Code:    FetchErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, http.StatusBadRequest)
			return
		}

		log.Printf("successfully fetched account id %s", session.AccountID)
		return
	}
}
