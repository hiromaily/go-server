package controller

import (
	"fmt"
	tm "github.com/hiromaily/go-server/libs/template"
	//lg "github.com/hiromaily/golibs/log"
	"net/http"
)

//GET
func GetGlobal(res http.ResponseWriter, req *http.Request) {
	fmt.Println("[Global]")
	//lg.Debugf("[req]%#v\n", req)

	//index
	tm.Execute(res, "global", nil)
}
