package main

import (
	"fmt"
	"github.com/gorilla/mux"
	h "github.com/hiromaily/go-server/libs/http"
	"net/http"
	//"golang.org/x/net/http2"
	lg "github.com/hiromaily/golibs/log"
	"os"
	"path/filepath"
)

func init() {
	lg.InitializeLog(1, lg.LogOff, 99, "/var/log/go/go-server.log",
		"[GO-SERVER]")
}

func main() {

	//For profiling of test
	h.SetProfile()

	//Setting Server
	r := mux.NewRouter()
	setRoute(r)

	//http2
	http2(r)
}

func http1(r *mux.Router) {
	http.ListenAndServe(":9999", r)
	lg.Info("Server start with port 9999")
}

func http2(r *mux.Router) {
	//get path executed command
	pwd, _ := os.Getwd()
	fmt.Println(pwd)
	// /Users/hy/work/go/src/github.com/hiromaily/go-server

	certFile, _ := filepath.Abs(pwd + "/ssl/cert1.pem")
	keyFile, _ := filepath.Abs(pwd + "/ssl/key1.pem")

	err := http.ListenAndServeTLS(":443", certFile, keyFile, r)
	if err != nil {
		lg.Error("ListenAndServeTLS:", err)
		http1(r)
	} else {
		lg.Info("Server start with port 443")
	}
}
