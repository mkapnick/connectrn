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
	// LoginErrCode code
	LoginErrCode = "account.login.error"
	// LoginExistsCode code
	LoginExistsCode = "account.login.exists"
)

type tokenResponse struct {
	Token string `json:"token"`
}

// Login checks email against password and assigns a token if valid
func Login(s account.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// only support POST
		if r.Method != http.MethodPost {
			err := errors.New("method not supported [login]")
			log.Printf(err.Error())
			resp := &je.Response{
				Code:    LoginErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, account.ServiceToHTTPErrorMap(err))
			return
		}

		var loginReq account.AccountCredentials
		err := json.NewDecoder(r.Body).Decode(&loginReq)
		if err != nil {
			log.Printf("%s: %v", LoginErrCode, err)
			resp := &je.Response{
				Code:    LoginErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, http.StatusBadRequest)
			return
		}

		token, err := s.LogIn(r.Context(), loginReq)
		if err != nil {
			// do not track failed log in attempts
			resp := &je.Response{
				Code:      LoginErrCode,
				Message:   err.Error(),
				SkipTrack: true,
			}
			je.Error(r, w, resp, account.ServiceToHTTPErrorMap(err))
			return
		}

		t := tokenResponse{
			Token: token,
		}

		// return logged in account
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(t)
		if err != nil {
			log.Printf("%s: %v", LoginErrCode, err)
			resp := &je.Response{
				Code:    LoginErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, http.StatusBadRequest)
			return
		}

		log.Printf("successfully logged in email %s", loginReq.Email)
		return
	}
}
