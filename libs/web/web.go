package web

import (
	"fmt"
	c "github.com/hiromaily/go-server/controller"
	lg "github.com/hiromaily/golibs/log"
	"net/http"
	"net/http/pprof"
	"runtime"
)

type (
	Web struct {
		Mux *http.ServeMux
		//Router []Router
		Router map[string][]Handler
	}
	Handler struct {
		Path string
		Func http.HandlerFunc
	}

	//Router struct {
	//	Method string
	//	Path   []string
	//}
)

func New() *Web {
	web := new(Web)
	web.Mux = http.NewServeMux()

	web.Router = map[string][]Handler{}

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

func setProfile() {
	//For profiling
	//runtime.MemProfileRate = 1
	runtime.SetBlockProfileRate(1)

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
	w.Router["GET"] = append(w.Router["GET"], Handler{Path: path, Func: f})
}

func (w *Web) Post(path string, f http.HandlerFunc) {
	//type HandlerFunc func(ResponseWriter, *Request)
	w.Router["POST"] = append(w.Router["POST"], Handler{Path: path, Func: f})
}

func (w *Web) Put(path string, f http.HandlerFunc) {
	//type HandlerFunc func(ResponseWriter, *Request)
	w.Router["PUT"] = append(w.Router["PUT"], Handler{Path: path, Func: f})
}

func (w *Web) Patch(path string, f http.HandlerFunc) {
	//type HandlerFunc func(ResponseWriter, *Request)
	w.Router["PATCH"] = append(w.Router["PATCH"], Handler{Path: path, Func: f})
}

func (w *Web) Delete(path string, f http.HandlerFunc) {
	//type HandlerFunc func(ResponseWriter, *Request)
	w.Router["DELETE"] = append(w.Router["DELETE"], Handler{Path: path, Func: f})
}

func (w *Web) Options(path string, f http.HandlerFunc) {
	//type HandlerFunc func(ResponseWriter, *Request)
	w.Router["OPTIONS"] = append(w.Router["OPTIONS"], Handler{Path: path, Func: f})
}

func (w *Web) Head(path string, f http.HandlerFunc) {
	//type HandlerFunc func(ResponseWriter, *Request)
	w.Router["HEAD"] = append(w.Router["HEAD"], Handler{Path: path, Func: f})
}

func (w *Web) handler(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	lg.Debugf("Method:%s, Path:%s", req.Method, req.URL.Path)

	var flg bool
	for _, el := range w.Router[req.Method] {
		if el.Path == req.URL.Path {
			el.Func(res, req)
			flg = true
			break
		}
	}

	if !flg {
		c.BadRequest(res, req)
	}
}

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

	lg.Info("TSL Server start with port 443 ...")
	err := http.ListenAndServeTLS(":443", cert, key, w.Mux)
	if err != nil {
		lg.Warn("ListenAndServeTLS:", err)

		lg.Infof("Server start with port %d ...", port)
		http.ListenAndServe(fmt.Sprintf(":%d", port), w.Mux)
	}
}

func (w *Web) StartServer(port int, cert, key string) {
	//http.HandleFunc("/", w.handler)
	w.Mux.Handle("/", http.HandlerFunc(w.handler))

	w.listen(port, cert, key)
}
