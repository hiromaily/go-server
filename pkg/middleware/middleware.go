package middleware

import (
	"net/http"
)

type (
	//Handler is http handler
	Handler func(http.ResponseWriter, *http.Request) (http.ResponseWriter, *http.Request)
	key     int
)
