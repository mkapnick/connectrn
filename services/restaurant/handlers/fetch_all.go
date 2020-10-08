package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	je "gitlab.com/michaelk99/birrdi/api-soa/internal/jsonerr"
	// "gitlab.com/michaelk99/birrdi/api-soa/internal/token"
	"gitlab.com/michaelk99/birrdi/api-soa/services/restaurant"
)

const (
	// FetchAllErrCode code
	FetchAllErrCode = "golfCourseRestaurant.fetchall.error"
	// FetchAllBadDataCode code
	FetchAllBadDataCode = "golfCourseRestaurant.fetchall.bad_data"
)

// FetchAll restaurants
func FetchAll(bs restaurant.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// only support GET
		if r.Method != http.MethodGet {
			err := errors.New("method not supported")
			log.Printf(err.Error())
			resp := &je.Response{
				Code:    FetchAllErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, restaurant.ServiceToHTTPErrorMap(err))
			return
		}

		ids := r.URL.Query().Get("golf_course_ids")

		if ids == "" {
			err := errors.New("golf_course_ids is a required field")
			log.Printf(err.Error())
			resp := &je.Response{
				Code:    FetchAllErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, restaurant.ServiceToHTTPErrorMap(err))
			return
		}

		date := r.URL.Query().Get("date")
		if date == "" {
			err := errors.New("date is a required field")
			log.Printf(err.Error())
			resp := &je.Response{
				Code:    FetchAllErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, restaurant.ServiceToHTTPErrorMap(err))
			return
		}

		// req can look like the following:
		// session := r.Context().Value("Session").(*token.Session)

		// TODO fetch golf course and make sure this user can view them
		/*
			ok := false
			for _, ar := range session.AccountRoles {
				if ar.CompanyID.String == gc.CompanyID {
					ok = true
				}
			}

			// no access, no entry ya boi
			if !ok {
				log.Printf("%s: %s", FetchAllErrCode, "Invalid admin requesting dates")
				resp := &je.Response{
					Code:    FetchAllErrCode,
					Message: "Invalid admin requesting dates",
				}
				je.Error(r, w, resp, restaurant.ServiceToHTTPErrorMap(err))
				return
			}
		*/

		b, err := bs.FetchAll(strings.Split(ids, ","), date)
		if err != nil {
			resp := &je.Response{
				Code:    FetchAllErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, restaurant.ServiceToHTTPErrorMap(err))
			return
		}

		// return created golfCourseRestaurant
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(b)
		if err != nil {
			log.Printf("%s: %v", FetchAllErrCode, err)
			resp := &je.Response{
				Code:    FetchAllErrCode,
				Message: err.Error(),
			}
			je.Error(r, w, resp, http.StatusBadRequest)
			return
		}

		log.Printf("successfully fetched all restaurants id in %s", date)
		return
	}
}
