package handlers

import (
	"log"
	"net/http"

	"errors"
	je "gitlab.com/michaelk99/birrdi/api-soa/internal/jsonerr"
	"gitlab.com/michaelk99/birrdi/api-soa/internal/token"
	"gitlab.com/michaelk99/birrdi/api-soa/services/account"
)

const (
	// DeleteErrCode code
	DeleteErrCode = "account.delete.error"
)

// Delete checks email against password and assigns a token if valid
func Delete(s account.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// only support DELETE
		if r.Method != http.MethodDelete {
			err := errors.New("method not supported [signup delete]")
			log.Printf(err.Error())
			resp := &je.Response{
				Code:    DeleteErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, account.ServiceToHTTPErrorMap(err))
			return
		}

		session := r.Context().Value("Session").(*token.Session)
		w.WriteHeader(http.StatusOK)
		log.Printf("successfully deleted account id %s", session.AccountID)
		return

		/*
			_, err := s.Delete(session.AccountID)
			if err != nil {
				fmt.Printf("%s", err)
				resp := &je.Response{
					Code:    DeleteErrCode,
					Message: err.Error(),
				}
				je.Error(r, w, resp, account.ServiceToHTTPErrorMap(err))
				return
			}

			// return deleted account
			w.WriteHeader(http.StatusOK)
			log.Printf("successfully deleted account id %s", session.AccountID)
			return
		*/
	}
}
