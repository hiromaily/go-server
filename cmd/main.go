package main

import (
	c "github.com/hiromaily/go-server/controller"
	"github.com/hiromaily/go-server/libs/web"
	lg "github.com/hiromaily/golibs/log"
	"os"
	"path/filepath"
)

var (
	certFile string
	keyFile  string
)

func init() {
	lg.InitializeLog(1, lg.LogOff, 99, "[GO-SERVER]",
		"/var/log/go/go-server.log")

	//get path executed command
	pwd, _ := os.Getwd()
	//fmt.Println(pwd)
	// /Users/hy/work/go/src/github.com/hiromaily/go-server

	certFile, _ = filepath.Abs(pwd + "/ssl/cert.pem")
	keyFile, _ = filepath.Abs(pwd + "/ssl/key.pem")
}

func setRoute(w *web.Web) {
	w.AttachProfiler()
	w.Get("/", c.Login)
}

func main() {

	//For profiling of test
	w := web.New()
	setRoute(w)

	w.StartServer(8080, certFile, keyFile)
}
