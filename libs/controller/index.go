package controller

import (
	"fmt"
	tm "github.com/hiromaily/go-server/libs/template"
	//lg "github.com/hiromaily/golibs/log"
	"net/http"
)

type Params struct {
	Title     string
	LinkNames []Link
}

type Link struct {
	Name string
	Url  string
}

//GET
func GetIndex(res http.ResponseWriter, req *http.Request) {
	fmt.Println("[Index]")
	//lg.Debugf("[req]%#v\n", req)

	data := Params{
		Title: "TitleIndex",
		LinkNames: []Link{
			Link{"index", "/"},
			Link{"login", "/login"},
			Link{"webpush", "/webpush"},
			Link{"global", "/global"},
		},
	}

	//index
	tm.Execute(res, "index", &data)
}
