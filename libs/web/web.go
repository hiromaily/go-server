package web

import (
	"context"
	"fmt"
	c "github.com/hiromaily/go-server/controller"
	mw "github.com/hiromaily/go-server/libs/middleware"
	lg "github.com/hiromaily/golibs/log"
	"net/http"
	"net/http/pprof"
	"runtime"
	"time"
)

type (
	Web struct {
		Mux *http.ServeMux
		//Router []Router
		Router map[string][]BHandler
		//Middleware []http.HandlerFunc
		Middleware []Middleware
	}
	BHandler struct {
		Path string
		Func http.HandlerFunc
	}
	Middleware func() mw.Handler
	//Handler    func(http.ResponseWriter, *http.Request) (http.ResponseWriter, *http.Request)
	//Handler func(context.Context, http.ResponseWriter, *http.Request) error

	//Router struct {
	//	Method string
	//	Path   []string
	//}
)

func setProfile() {
	//For profiling
	//runtime.MemProfileRate = 1
	runtime.SetBlockProfileRate(1)

}

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

	//web.Router = append(web.Router, Router{Method: "GET"})
	//web.Router = append(web.Router, Router{Method: "POST"})
	//web.Router = append(web.Router, Router{Method: "PUT"})
	//web.Router = append(web.Router, Router{Method: "PATCH"})
	//web.Router = append(web.Router, Router{Method: "DELETE"})
	//web.Router = append(web.Router, Router{Method: "OPTIONS"})
	//web.Router = append(web.Router, Router{Method: "HEAD"})

	return web
}

//func (w *Web) Use(f http.HandlerFunc) {
func (w *Web) Use(m Middleware) {
	w.Middleware = append(w.Middleware, m)
}

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

func (w *Web) Get(path string, f http.HandlerFunc) {
	//type HandlerFunc func(ResponseWriter, *Request)
	w.Router["GET"] = append(w.Router["GET"], BHandler{Path: path, Func: f})
}

func (w *Web) Post(path string, f http.HandlerFunc) {
	//type HandlerFunc func(ResponseWriter, *Request)
	w.Router["POST"] = append(w.Router["POST"], BHandler{Path: path, Func: f})
}

func (w *Web) Put(path string, f http.HandlerFunc) {
	//type HandlerFunc func(ResponseWriter, *Request)
	w.Router["PUT"] = append(w.Router["PUT"], BHandler{Path: path, Func: f})
}

func (w *Web) Patch(path string, f http.HandlerFunc) {
	//type HandlerFunc func(ResponseWriter, *Request)
	w.Router["PATCH"] = append(w.Router["PATCH"], BHandler{Path: path, Func: f})
}

func (w *Web) Delete(path string, f http.HandlerFunc) {
	//type HandlerFunc func(ResponseWriter, *Request)
	w.Router["DELETE"] = append(w.Router["DELETE"], BHandler{Path: path, Func: f})
}

func (w *Web) Options(path string, f http.HandlerFunc) {
	//type HandlerFunc func(ResponseWriter, *Request)
	w.Router["OPTIONS"] = append(w.Router["OPTIONS"], BHandler{Path: path, Func: f})
}

func (w *Web) Head(path string, f http.HandlerFunc) {
	//type HandlerFunc func(ResponseWriter, *Request)
	w.Router["HEAD"] = append(w.Router["HEAD"], BHandler{Path: path, Func: f})
}

func (w *Web) handler(res http.ResponseWriter, req *http.Request) {
	lg.Debugf("Method:%s, Path:%s", req.Method, req.URL.Path)

	//execute middleware
	rw, r := w.execMiddleware(res, req)

	//After middleware
	ctx := r.Context()
	id, err := mw.GetRequestID(ctx)
	fmt.Printf("[%d] handler() started\n", id)

	//set timeout
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer func() {
		id, err = mw.GetRequestID(ctx)
		fmt.Printf("[%d] handler() ended\n", id)
		cancel()
	}()

	//common work for form
	r.ParseForm()

	//execute main function
	ch := make(chan int, 1)
	go w.execMainFunc(rw, r.WithContext(ctx), ch)

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

func (w *Web) execMainFunc(res http.ResponseWriter, req *http.Request, ch chan<- int) {
	var flg bool
	//test
	time.Sleep(5 * time.Second)
	for _, el := range w.Router[req.Method] {
		if el.Path == req.URL.Path {
			el.Func(res, req)
			flg = true
			//TODO:cancel timeout
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

//func contextHandler(ctx context.Context, res http.ResponseWriter) {
//	//
//	select {
//	case <-ctx.Done():
//		err := ctx.Err()
//		if err == context.Canceled {
//			fmt.Println("context.Canceled")
//		} else if err == context.DeadlineExceeded {
//			fmt.Println("context.DeadlineExceeded")
//		} else {
//			fmt.Println("else")
//		}
//		id, err := mw.GetRequestID(ctx)
//		http.Error(res, fmt.Sprintf("[%d] %s", id, err.Error()), http.StatusInternalServerError)
//	//default:
//	//	fmt.Println("default")
//	}
//}

func (w *Web) displayPaths() {
	var search = func(method string) {
		for _, el := range w.Router[method] {
			fmt.Printf("[%s] %s\n", method, el.Path)
		}
	}
	search("GET")
	search("POST")
}

func (w *Web) listen(port int, cert, key string) {

	//
	w.displayPaths()

	lg.Infof("Server start with port %d ...", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), w.Mux)

	//lg.Info("TSL Server start with port 443 ...")
	//err := http.ListenAndServeTLS(":443", cert, key, w.Mux)
	//if err != nil {
	//	lg.Warn("ListenAndServeTLS:", err)
	//
	//	lg.Infof("Server start with port %d ...", port)
	//	http.ListenAndServe(fmt.Sprintf(":%d", port), w.Mux)
	//}
}

func (w *Web) StartServer(port int, cert, key string) {
	//http.HandleFunc("/", w.handler)
	w.Mux.Handle("/", http.HandlerFunc(w.handler))
	//w.Mux.HandleFunc("/", mw.Decorate(w.handler))

	w.listen(port, cert, key)
}
