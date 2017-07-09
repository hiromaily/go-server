package controller

import (
	"fmt"
	tm "github.com/hiromaily/go-server/libs/template"
	//lg "github.com/hiromaily/golibs/log"
	"net/http"
)

type Params struct {
	Str1   string
	Int1   int
	Bool1  bool
	Slice1 []string
}

//GET
func GetIndex(res http.ResponseWriter, req *http.Request) {
	fmt.Println("[Index]")
	//lg.Debugf("[req]%#v\n", req)

	data := Params{Str1: "test", Int1: 100, Bool1: false, Slice1: []string{"aaa", "bbb", "ccc"}}

	//index
	tm.Execute(res, "index", &data)
}
