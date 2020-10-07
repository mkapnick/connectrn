package handlers

import (
	"errors"
	"log"
	"net/http"

	je "gitlab.com/michaelk99/birrdi/api-soa/internal/jsonerr"
	"gitlab.com/michaelk99/birrdi/api-soa/internal/token"
	"gitlab.com/michaelk99/birrdi/api-soa/services/profile"
)

const (
	// DeleteErrCode code
	DeleteErrCode = "profile.delete.error"
)

// Delete checks email against password and assigns a token if valid
func Delete(ps profile.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// only support DELETE
		if r.Method != http.MethodDelete {
			err := errors.New("method not supported [profile delete]")
			log.Printf(err.Error())
			resp := &je.Response{
				Code:    DeleteErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, profile.ServiceToHTTPErrorMap(err))
			return
		}

		// look for the id in the session
		session := r.Context().Value("Session").(*token.Session)
		id := session.ProfileID

		if id == "" {
			err := errors.New(DeleteErrCode)
			log.Printf(err.Error())
			resp := &je.Response{
				Code:    DeleteErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, profile.ServiceToHTTPErrorMap(err))
			return
		}

		query := profile.IDQuery{
			Type:  profile.ID,
			Value: id,
		}

		err := ps.Delete(query)
		if err != nil {
			resp := &je.Response{
				Code:    DeleteErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, profile.ServiceToHTTPErrorMap(err))
			return
		}

		// return deleted profile
		w.WriteHeader(http.StatusOK)
		log.Printf("successfully deleted profile id %s", id)
		return
	}
}
