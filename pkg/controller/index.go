package controller

import (
	"fmt"

	tm "github.com/hiromaily/go-server/pkg/template"

	//lg "github.com/hiromaily/golibs/log"
	"net/http"
)

// Params is parameter for response html
type Params struct {
	Title     string
	LinkNames []Link
}

// Link is parameter for response html
type Link struct {
	Name string
	URL  string
}

// GetIndex is for / or /index page
func GetIndex(res http.ResponseWriter, req *http.Request) {
	fmt.Println("[Index]")
	//lg.Debugf("[req]%#v\n", req)
	fmt.Println(" ip address:", req.RemoteAddr)

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
