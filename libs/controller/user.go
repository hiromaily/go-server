package controller

import (
	tm "github.com/hiromaily/go-server/libs/template"
	lg "github.com/hiromaily/golibs/log"
	"net/http"
)

//GET
func GetLogin(res http.ResponseWriter, req *http.Request) {
	lg.Info("[User Login]")
	lg.Debugf("[req]%+v\n", req)

	//fmt.Fprintf(res, "User Login")
	//ctx := req.Context()

	//index
	tm.Execute(res, "login", nil)
}
