package controller

import (
	"fmt"
	lg "github.com/hiromaily/golibs/log"
	"net/http"
)

//GET
func GetIndex(res http.ResponseWriter, req *http.Request) {
	fmt.Println("[Index]")
	lg.Debugf("[req]%#v\n", req)
	fmt.Fprintf(res, "Index")
	//index
	//if err := tmpl.ExecuteTemplate(w, "layout", nil); err != nil {
	//	log.Println(err.Error())
	//	http.Error(w, http.StatusText(500), 500)
	//}

	//ctx := req.Context()
}
