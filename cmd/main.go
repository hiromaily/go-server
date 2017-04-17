package main

import (
	"flag"
	c "github.com/hiromaily/go-server/controller"
	mw "github.com/hiromaily/go-server/libs/middleware"
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
	flag.Parse()

	lg.InitializeLog(1, lg.LogOff, 99, "[GO-SERVER]",
		"/var/log/go/go-server.log")

	//For TSL
	//get path executed command
	pwd, _ := os.Getwd()

	certFile, _ = filepath.Abs(pwd + "/ssl/cert.pem")
	keyFile, _ = filepath.Abs(pwd + "/ssl/key.pem")
}

func setMiddleware(w *web.Web) {
	w.Use(mw.SetRequestID) //func() web.Middleware
	//w.Use(middleware.LogRequest(true))
	//w.Use(middleware.ErrorHandler(w, true))
	//w.Use(middleware.Recover())
}

func setRoute(w *web.Web) {
	w.AttachProfiler()
	//Add Router
	w.Get("/", c.Login)
}

func main() {

	//For profiling of test
	w := web.New()
	setMiddleware(w)
	setRoute(w)

	w.StartServer(8080, certFile, keyFile)
}
