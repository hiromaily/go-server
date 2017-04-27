package controller

import (
	"fmt"
	lg "github.com/hiromaily/golibs/log"
	"net/http"
)

//GET
func GetLogin(res http.ResponseWriter, req *http.Request) {
	fmt.Println("[User Login]")
	lg.Debugf("[req]%+v\n", req)

	fmt.Fprintf(res, "User Login")

	//ctx := req.Context()
}
