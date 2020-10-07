package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	je "gitlab.com/michaelk99/birrdi/api-soa/internal/jsonerr"
	"gitlab.com/michaelk99/birrdi/api-soa/services/account"
	"gitlab.com/michaelk99/birrdi/api-soa/services/account/handlers"
	"github.com/stretchr/testify/assert"
)

var TestUpdateTT = []struct {
	// name for test case
	name string
	// name of the http method to use
	method string
	// path used to in request
	path string
	// the account object being updated
	acc account.Account
	// do we expect an invalid ID error
	expectInvalidID bool
	// do we expect an ID mismatch between acc and url path
	expectIDMisMatch bool
	// do we expect the fetch method to be called.
	expectFetchCall bool
	// do we expect service.Fetch to fail
	expectFetchFail bool
	// returned values from mock fetch call
	fetchReturn []interface{}
	// do we expect update method to be called
	expectUpdateCall bool
	// do we expect service.Update to fail
	expectUpdateFail bool
	// returned values from mock update call
	updateReturn []interface{}
}{
	{
		name:   "successful update",
		method: "PUT",
		path:   "/api/v1/account/test-user-id",
		acc: account.Account{
			Email:   "newemail@newemail.com",
			Enabled: true,
			ID:      "test-user-id",
		},
		expectInvalidID:  false,
		expectIDMisMatch: false,
		expectFetchCall:  true,
		expectFetchFail:  false,
		fetchReturn: []interface{}{
			&account.Account{
				Email: "oldemail@oldemail.com",
				ID:    "test-user-id",
			},
			nil,
		},
		expectUpdateCall: true,
		expectUpdateFail: false,
		updateReturn: []interface{}{
			&account.Account{
				Email:   "newemail@newemail.com",
				Enabled: true,
				ID:      "test-user-id",
			},
			nil,
		},
	},
	{
		name:   "no id provided in path",
		method: "PUT",
		path:   "/api/v1/account/",
		acc: account.Account{
			Email:   "newemail@newemail.com",
			Enabled: true,
			ID:      "test-user-id",
		},
		expectInvalidID:  true,
		expectIDMisMatch: false,
		expectFetchCall:  false,
		expectFetchFail:  false,
		fetchReturn: []interface{}{
			&account.Account{
				Email: "oldemail@oldemail.com",
				ID:    "test-user-id",
			},
			nil,
		},
		expectUpdateCall: false,
		expectUpdateFail: false,
		updateReturn: []interface{}{
			&account.Account{
				Email:   "newemail@newemail.com",
				Enabled: true,
				ID:      "test-user-id",
			},
			nil,
		},
	},
	{
		name:   "id mismatch",
		method: "PUT",
		path:   "/api/v1/account/original-id",
		acc: account.Account{
			Email:   "newemail@newemail.com",
			Enabled: true,
			ID:      "test-user-id",
		},
		expectInvalidID:  false,
		expectIDMisMatch: true,
		expectFetchCall:  false,
		expectFetchFail:  false,
		fetchReturn: []interface{}{
			&account.Account{
				Email: "oldemail@oldemail.com",
				ID:    "new-id",
			},
			nil,
		},
		expectUpdateCall: false,
		expectUpdateFail: false,
		updateReturn: []interface{}{
			&account.Account{
				Email:   "newemail@newemail.com",
				Enabled: true,
				ID:      "test-user-id",
			},
			nil,
		},
	},
	{
		name:   "fetch fail",
		method: "PUT",
		path:   "/api/v1/account/test-user-id",
		acc: account.Account{
			Email:   "newemail@newemail.com",
			Enabled: true,
			ID:      "test-user-id",
		},
		expectInvalidID:  false,
		expectIDMisMatch: false,
		expectFetchCall:  true,
		expectFetchFail:  true,
		fetchReturn: []interface{}{
			nil,
			fmt.Errorf("failed to fetch user"),
		},
		expectUpdateCall: false,
		expectUpdateFail: false,
		updateReturn: []interface{}{
			nil,
			nil,
		},
	},
	{
		name:   "update fail",
		method: "PUT",
		path:   "/api/v1/account/test-user-id",
		acc: account.Account{
			Email:   "newemail@newemail.com",
			Enabled: true,
			ID:      "test-user-id",
		},
		expectInvalidID:  false,
		expectIDMisMatch: false,
		expectFetchCall:  true,
		expectFetchFail:  false,
		fetchReturn: []interface{}{
			&account.Account{
				Email: "oldemail@oldemail.com",
				ID:    "new-id",
			},
			nil,
		},
		expectUpdateCall: true,
		expectUpdateFail: true,
		updateReturn: []interface{}{
			nil,
			fmt.Errorf("failed to update"),
		},
	},
}

func extractError(rr *httptest.ResponseRecorder) (je.Response, error) {
	var res je.Response
	err := json.Unmarshal(rr.Body.Bytes(), &res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func TestUpdate(t *testing.T) {
	// create mock account service
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range TestUpdateTT {
		t.Logf("test table: %v", tt.name)

		// create our mock service to provide handler
		s := account.NewMockService(ctrl)

		if tt.expectFetchCall {
			idQuery := account.IDQuery{
				Type:  account.ID,
				Value: tt.acc.ID,
			}
			s.EXPECT().Fetch(idQuery).Return(tt.fetchReturn...)
		}

		if tt.expectUpdateCall {
			s.EXPECT().Update(tt.acc).Return(tt.updateReturn...)
		}

		// create handler and call
		h := handlers.Update(s)

		b, err := json.Marshal(tt.acc)
		if err != nil {
			t.Fatalf("failed to serialize account to be updated: %v", err)
		}
		req := httptest.NewRequest(tt.method, tt.path, bytes.NewBuffer(b))

		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, req)

		if tt.expectInvalidID {
			jsonErr, err := extractError(rr)
			if err != nil {
				t.Fatalf("failed to deserialize json error type: %v", err)
			}

			assert.Equal(t, handlers.UpdateErrInvalidIDCode, jsonErr.Code)
			assert.Equal(t, http.StatusBadRequest, rr.Code)
			continue
		}

		if tt.expectIDMisMatch {
			jsonErr, err := extractError(rr)
			if err != nil {
				t.Fatalf("failed to deserialize json error type")
				t.Fatalf("failed to deserialize json error type: %v", err)
			}

			assert.Equal(t, handlers.UpdateErrCode, jsonErr.Code)
			assert.Equal(t, http.StatusBadRequest, rr.Code)
			continue
		}

		if tt.expectFetchFail {
			jsonErr, err := extractError(rr)
			if err != nil {
				t.Fatalf("failed to deserialize json error type: %v", err)
			}

			assert.Equal(t, handlers.UpdateErrAccountNotFoundCode, jsonErr.Code)
			assert.Equal(t, http.StatusBadRequest, rr.Code)
			continue
		}

		if tt.expectUpdateFail {
			jsonErr, err := extractError(rr)
			if err != nil {
				t.Fatalf("failed to deserialize json error type: %v", err)
			}

			assert.Equal(t, handlers.UpdateErrCode, jsonErr.Code)
			assert.Equal(t, http.StatusInternalServerError, rr.Code)
			continue
		}

		// all errors accounted for. test valid response
		var a account.Account
		err = json.Unmarshal(rr.Body.Bytes(), &a)
		if err != nil {
			t.Fatalf("failed to deserialize json error type: %v", err)
		}

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, tt.acc, a)
	}
}
