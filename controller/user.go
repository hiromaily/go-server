package controller

import (
	"fmt"
	"net/http"
)

func Login(res http.ResponseWriter, req *http.Request) {
	fmt.Println("[User Login]")
	fmt.Fprintf(res, "User Login")

	//ctx := req.Context()
}
