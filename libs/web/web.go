package web

import (
	"context"
	"crypto/tls"
	"fmt"
	c "github.com/hiromaily/go-server/libs/controller"
	mw "github.com/hiromaily/go-server/libs/middleware"
	lg "github.com/hiromaily/golibs/log"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/alexedwards/scs/engine/memstore"
	"github.com/alexedwards/scs/session"
)

type (
	//Web is web object including Mux, Router, Middleware
	Web struct {
		Mux        *http.ServeMux
		Router     map[string][]BHandler
		Middleware []Middleware
	}
	//BHandler is handler path and func
	BHandler struct {
		Path string
		Func http.HandlerFunc
	}
	//Middleware is handler func for middleware
	Middleware func() mw.Handler
)

func setProfile() {
	//For profiling
	//runtime.MemProfileRate = 1
	runtime.SetBlockProfileRate(1)
}

// New is to create Web object
func New() *Web {
	web := new(Web)
	web.Mux = http.NewServeMux()

	web.Router = map[string][]BHandler{}

	web.Router["GET"] = nil
	web.Router["POST"] = nil
	web.Router["PUT"] = nil
	web.Router["PATCH"] = nil
	web.Router["DELETE"] = nil
	web.Router["OPTIONS"] = nil
	web.Router["HEAD"] = nil

	return web
}

// Use is to set middleware you want
func (w *Web) Use(m Middleware) {
	w.Middleware = append(w.Middleware, m)
}

// AttachProfiler is for profiler pages
func (w *Web) AttachProfiler() {
	setProfile()

	w.Mux.HandleFunc("/debug/pprof/", pprof.Index)
	w.Mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	w.Mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	w.Mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)

	// Manually add support for paths linked to by index page at /debug/pprof/
	w.Mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	w.Mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	w.Mux.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	w.Mux.Handle("/debug/pprof/block", pprof.Handler("block"))
}

// SetStaticFiles is to set static files
func (w *Web) SetStaticFiles() {
	//static files
	w.Mux.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./statics/img"))))
	w.Mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./statics/css"))))
	w.Mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./statics/js"))))

	w.Mux.HandleFunc("/favicon.ico", faviconHandler)
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	lg.Debug("[]favicon.ico was requested")
	http.ServeFile(w, r, "./statics/favicon.ico")
}

// Get is to set get request
func (w *Web) Get(path string, f http.HandlerFunc) {
	w.Router["GET"] = append(w.Router["GET"], BHandler{Path: path, Func: f})
}

// Post is to set post request
func (w *Web) Post(path string, f http.HandlerFunc) {
	w.Router["POST"] = append(w.Router["POST"], BHandler{Path: path, Func: f})
}

// Put is to set put request
func (w *Web) Put(path string, f http.HandlerFunc) {
	w.Router["PUT"] = append(w.Router["PUT"], BHandler{Path: path, Func: f})
}

// Patch is to set patch request
func (w *Web) Patch(path string, f http.HandlerFunc) {
	w.Router["PATCH"] = append(w.Router["PATCH"], BHandler{Path: path, Func: f})
}

// Delete is to set delete request
func (w *Web) Delete(path string, f http.HandlerFunc) {
	w.Router["DELETE"] = append(w.Router["DELETE"], BHandler{Path: path, Func: f})
}

// Options is to set options request
func (w *Web) Options(path string, f http.HandlerFunc) {
	w.Router["OPTIONS"] = append(w.Router["OPTIONS"], BHandler{Path: path, Func: f})
}

// Head is to set head request
func (w *Web) Head(path string, f http.HandlerFunc) {
	w.Router["HEAD"] = append(w.Router["HEAD"], BHandler{Path: path, Func: f})
}

// endpoint of router
func (w *Web) Handler(res http.ResponseWriter, req *http.Request) {
	lg.Debugf("Method:%s, Path:%s", req.Method, req.URL.Path)

	//TODO:check file path first
	fn := w.findFunc(req)
	if fn == nil {
		c.BadRequest(res, req)
		return
	}

	ch := make(chan int, 1)

	//execute middleware
	rw, r := w.execMiddleware(res, req)

	//After middleware
	ctx := r.Context()
	//id, err := mw.GetRequestID(ctx)
	//fmt.Printf("[%d] handler() started\n", id)

	//set timeout
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer func() {
		//id, err = mw.GetRequestID(ctx)
		//fmt.Printf("[%d] handler() ended\n", id)
		cancel()
	}()

	//parse form
	r.ParseForm()

	//execute main function
	//go w.execMainFunc(rw, r.WithContext(ctx), ch)
	go func() {
		fn(rw, r.WithContext(ctx))
		ch <- 1
	}()

	//contextHandler
	select {
	case <-ctx.Done():
		err := ctx.Err()
		if err == context.Canceled {
			fmt.Println("context.Canceled")
		} else if err == context.DeadlineExceeded {
			fmt.Println("context.DeadlineExceeded")
		} else {
			fmt.Println("else")
		}
	case <-ch:
		fmt.Println("finished correctly")
	}
}

func (w *Web) execMiddleware(res http.ResponseWriter, req *http.Request) (http.ResponseWriter, *http.Request) {
	rw := res
	r := req
	for _, mw := range w.Middleware {
		f := mw() //return handler
		rw, r = f(rw, r)
	}
	return rw, r
}

func (w *Web) findFunc(req *http.Request) http.HandlerFunc {
	for _, el := range w.Router[req.Method] {
		if el.Path == req.URL.Path {
			return el.Func
		}
	}
	return nil
}

func (w *Web) execMainFunc(res http.ResponseWriter, req *http.Request, ch chan<- int) {
	var flg bool
	//test
	//time.Sleep(5 * time.Second)
	for _, el := range w.Router[req.Method] {
		if el.Path == req.URL.Path {
			el.Func(res, req)
			flg = true
			//send done
			ch <- 1
			break
		}
	}

	if !flg {
		c.BadRequest(res, req)
		//send done
		ch <- 1
	}
}

// StartServer is to start server with setting handler
func (w *Web) StartServer(port int, cert, key string) {
	w.Mux.Handle("/", http.HandlerFunc(w.Handler))

	w.listen2(port, cert, key)
}

func getTLSConf(cert, key string) (*tls.Config, error) {
	if cert == "" || key == "" {
		return nil, fmt.Errorf("%s", "parameters are invalid.")
	}

	cer, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}
	tlsConf := tls.Config{Certificates: []tls.Certificate{cer}}
	return &tlsConf, nil
}

func (w *Web) listen(port int, cert, key string) {
	//
	w.displayPaths()

	if cert != "" && key != "" {
		lg.Info("TSL Server start with port 443 ...")
		err := http.ListenAndServeTLS(":443", cert, key, w.Mux)
		if err != nil {
			lg.Warn("ListenAndServeTLS:", err)
		}
	}
	lg.Infof("Server start with port %d ...", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), w.Mux)
}

func (w *Web) listen2(port int, cert, key string) {
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	//
	w.displayPaths()

	// session
	engine := memstore.New(0)
	sessionManager := session.Manage(engine)

	//server object
	var srv *http.Server
	conf, err := getTLSConf(cert, key)
	if err != nil {
		//srv = &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: w.Mux}
		srv = &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: sessionManager(w.Mux)}
		lg.Infof("Server start with port %d ...", port)
		go func() {
			// service connections
			if err := srv.ListenAndServe(); err != nil {
				lg.Errorf("listen: %s\n", err)
			}
		}()
	} else {
		//srv = &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: w.Mux, TLSConfig: conf}
		srv = &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: sessionManager(w.Mux), TLSConfig: conf}
		lg.Info("TSL Server start with port 443 ...")
		go func() {
			err := srv.ListenAndServeTLS("", "")
			if err != nil {
				lg.Warn("ListenAndServeTLS:", err)
			}
		}()
	}

	<-stopChan // wait for SIGINT
	lg.Info("Shutting down server...")

	// shut down gracefully, but wait no longer than 5 seconds before halting
	//TODO:the cancel function returned by context.WithTimeout should be called, not discarded, to avoid a context leak
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	//defer func() {
	//	cancel()
	//}()

	srv.Shutdown(ctx)
	lg.Info("Server gracefully stopped")
}

func (w *Web) displayPaths() {
	var search = func(method string) {
		for _, el := range w.Router[method] {
			fmt.Printf("[%s] %s\n", method, el.Path)
		}
	}
	search("GET")
	search("POST")
	search("PUT")
	search("PATCH")
	search("DELETE")
	search("OPTIONS")
	search("HEAD")
}
