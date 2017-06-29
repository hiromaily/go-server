package main

import (
	"flag"
	conf "github.com/hiromaily/go-server/libs/config"
	c "github.com/hiromaily/go-server/libs/controller"
	mw "github.com/hiromaily/go-server/libs/middleware"
	tm "github.com/hiromaily/go-server/libs/template"
	"github.com/hiromaily/go-server/libs/web"
	"github.com/hiromaily/golibs/auth/jwt"
	enc "github.com/hiromaily/golibs/cipher/encryption"
	lg "github.com/hiromaily/golibs/log"
	"os"
	"path/filepath"
)

//TODO:LIST
//[initialization]
//* template files
//* connect to Database (postgreSQL)
//* connect to Redis

//[middleware]
//* cookies
//* jwt
//* check http header and store on context

//[tomlFiles]

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

	//cipher
	enc.NewCryptDefault()

	// config
	cnf := conf.New(*tomlPath, true)

	//log
	lg.InitializeLog(cnf.Server.Log.Level, lg.LogOff, 99, "[GO-SERVER]",
		cnf.Server.Log.Path)

	//For TSL
	if *tsl {
		//get path executed command
		pwd, _ := os.Getwd()

		certFile, _ = filepath.Abs(pwd + "/data/ssl/cert.pem")
		keyFile, _ = filepath.Abs(pwd + "/data/ssl/key.pem")
	}
}

func initAuth() {
	auth := conf.GetConf().API.JWT
	if auth.Mode == jwt.HMAC && auth.Secret != "" {
		jwt.InitSecretKey(auth.Secret)
	} else if auth.Mode == jwt.RSA && auth.PrivateKey != "" && auth.PublicKey != "" {
		err := jwt.InitKeys(auth.PrivateKey, auth.PublicKey)
		if err != nil {
			lg.Error(err)
			panic(err)
		}
	} else {
		jwt.InitEncrypted(jwt.HMAC)
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
