package middleware

import (
	"net/http"
)

// SetHTTPHeader is to set http header for response
func SetHTTPHeader() Handler {
	return func(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, *http.Request) {
		//r.Header.Set("GoServer", "something value")
		w.Header().Set("GoServer", "something value")
		return w, r
	}
}
