package main

import (
	"flag"
	c "github.com/hiromaily/go-server/controller"
	mw "github.com/hiromaily/go-server/libs/middleware"
	tm "github.com/hiromaily/go-server/libs/template"
	"github.com/hiromaily/go-server/libs/web"
	lg "github.com/hiromaily/golibs/log"
	"os"
	"path/filepath"
)

var (
	certFile string
	keyFile  string
)

var (
	tomlPath = flag.String("f", "", "Toml file path")
	portNum  = flag.Int("p", 8080, "Port of server")
	tsl      = flag.Bool("s", false, "True means to run TSL server")
)

func init() {
	flag.Parse()

	lg.InitializeLog(1, lg.LogOff, 99, "[GO-SERVER]",
		"/var/log/go/go-server.log")

	//For TSL
	if *tsl {
		//get path executed command
		pwd, _ := os.Getwd()

		certFile, _ = filepath.Abs(pwd + "/ssl/cert.pem")
		keyFile, _ = filepath.Abs(pwd + "/ssl/key.pem")
	}
}

func setMiddleware(w *web.Web) {
	w.Use(mw.SetRequestID) //func() web.Middleware
	//w.Use(middleware.LogRequest(true))
	//w.Use(middleware.ErrorHandler(w, true))
	//w.Use(middleware.Recover())
}

func setRoute(w *web.Web) {
	w.AttachProfiler()
	w.SetStaticFiles()
	tm.LoadTemplatesFiles()
	//Add Router
	w.Get("/", c.GetIndex)
	w.Get("/login", c.GetLogin)
}

func main() {

	//For profiling of test
	w := web.New()
	setMiddleware(w)
	setRoute(w)

	w.StartServer(*portNum, certFile, keyFile)
}
