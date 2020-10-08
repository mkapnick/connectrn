package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	je "gitlab.com/michaelk99/connectrn/internal/jsonerr"
	"gitlab.com/michaelk99/connectrn/internal/token"
	"gitlab.com/michaelk99/connectrn/services/profile"
)

const (
	// FetchErrCode code
	FetchErrCode = "profile.fetch.error"
	// FetchExistsCode code
	FetchExistsCode = "profile.fetch.exists"
)

// Fetch checks email against password and assigns a token if valid
func Fetch(ps profile.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// req can look like the following:
		session := r.Context().Value("Session").(*token.Session)
		id := session.ProfileID
		accountID := session.AccountID

		// !!! IMPORTANT !!!
		// ?account_id query param overrides the session. This allows clients
		// to fetch other profiles on demand. Ideally this gets separated out
		// into a /search endpoint. For now it lives here
		if r.URL.Query().Get("account_id") != "" {
			accountID = r.URL.Query().Get("account_id")
		}

		if id == "" && accountID == "" {
			err := errors.New("id or accound id are required")
			log.Printf(err.Error())
			resp := &je.Response{
				Code:    FetchErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, profile.ServiceToHTTPErrorMap(err))
			return
		}

		// default query to profile id
		query := profile.IDQuery{
			Type:  profile.ID,
			Value: id,
		}

		// accountID overrides profile ID
		if accountID != "" {
			query = profile.IDQuery{
				Type:  profile.AccountID,
				Value: accountID,
			}
		}

		prof, err := ps.Fetch(query)
		if err != nil {
			resp := &je.Response{
				Code:    FetchErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, profile.ServiceToHTTPErrorMap(err))
			return
		}

		// return created profile
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(prof)
		if err != nil {
			log.Printf("%s: %v", FetchErrCode, err)
			resp := &je.Response{
				Code:    FetchErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, http.StatusBadRequest)
			return
		}

		log.Printf("successfully fetched profile id %s", id)
		return
	}
}
