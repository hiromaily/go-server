package controller

import (
	"fmt"
	"net/http"
)

func Login(res http.ResponseWriter, req *http.Request) {
	fmt.Print("User Login")
	fmt.Fprintf(res, "User Login")
}
