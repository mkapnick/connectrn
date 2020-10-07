package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	je "gitlab.com/michaelk99/birrdi/api-soa/internal/jsonerr"
	"gitlab.com/michaelk99/birrdi/api-soa/services/profile"
)

const (
	// SearchErrCode code
	SearchErrCode = "profile.fetch.error"
	// SearchExistsCode code
	SearchExistsCode = "profile.fetch.exists"
)

// Search checks email against password and assigns a token if valid
func Search(ps profile.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// only support POST
		if r.Method != http.MethodGet {
			err := errors.New("method not supported [profile search]")
			log.Printf(err.Error())
			resp := &je.Response{
				Code:    SearchErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, profile.ServiceToHTTPErrorMap(err))
			return
		}

		if r.URL.Query().Get("value") == "" {
			err := errors.New("value is a required field [profile search]")
			log.Printf(err.Error())
			resp := &je.Response{
				Code:    SearchErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, profile.ServiceToHTTPErrorMap(err))
			return
		}

		v := r.URL.Query().Get("value")
		query := profile.IDQuery{
			Type:  profile.AdminCheckout,
			Value: v,
		}

		golfCourseID := r.URL.Query().Get("golf_course_id")
		if golfCourseID == "" {
			err := errors.New("golf_course_id is a required field [profile search]")
			log.Printf(err.Error())
			resp := &je.Response{
				Code:    SearchErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, profile.ServiceToHTTPErrorMap(err))
			return
		}

		// the `from_loyalty` flag tells us this is a loyalty search
		fromLoyalty := r.URL.Query().Get("from_loyalty")
		if fromLoyalty != "" {
			query = profile.IDQuery{
				Type:  profile.FromLoyalty,
				Value: v,
			}
		}

		jwt := r.Context().Value("Token").(string)

		profs, err := ps.Search(jwt, query, golfCourseID)
		if err != nil {
			resp := &je.Response{
				Code:    SearchErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, profile.ServiceToHTTPErrorMap(err))
			return
		}

		// return created profile
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(profs)
		if err != nil {
			log.Printf("%s: %v", SearchErrCode, err)
			resp := &je.Response{
				Code:    SearchErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, http.StatusBadRequest)
			return
		}

		log.Printf("successfully searched profiles %s", v)
		return
	}
}
