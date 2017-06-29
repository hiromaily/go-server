package controller

import (
	//"fmt"
	"net/http"
)

func BadRequest(res http.ResponseWriter, req *http.Request) {
	//fmt.Print("Bad Request: 404")
	//fmt.Fprintf(res, "404")
	http.NotFound(res, req)
}
