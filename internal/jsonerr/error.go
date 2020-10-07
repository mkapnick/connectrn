package jsonerr

import (
	"encoding/json"
	"net/http"
)

// Additional any additional info
type Additional interface{}

// Response json response
type Response struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	// SkipTrack skips entry tracking. Defaults to `false`
	SkipTrack bool `json:"omitempty"`
	// Additional must be json serializable or expect errors
	Additional `json:"additional,omitempty"`
}

// Error JsonError works like http.Error but uses our response
// struct as the body of the response. Like http.Error
// you will still need to call a naked return in the http handler
func Error(req *http.Request, w http.ResponseWriter, r *Response, httpcode int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	w.WriteHeader(httpcode)
	b, _ := json.Marshal(r)

	w.Write(b)
}
