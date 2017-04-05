package main

import (
	"github.com/gorilla/mux"
	"net/http"
	h "github.com/hiromaily/go-server/http"
)

func main() {
	//For profiling of test
	h.SetProfile()

	//Setting Server
	r := mux.NewRouter()
	setRoute(r)

	http.ListenAndServe(":9999", r)
}
