package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"math"
	"net/http"
	"net/http/pprof"
	"runtime"
)

func setRoute(r *mux.Router) {
	//r.GET("/", Index)
	r.HandleFunc("/", Index)

	//For debug mode or development
	attachProfiler(r)
}

func attachProfiler(r *mux.Router) {
	//r.GET("/debug/pprof/", pprof.Index)
	//r.GET("/debug/pprof/cmdline", pprof.Cmdline)
	//r.GET("/debug/pprof/profile", pprof.Profile)
	//r.GET("/debug/pprof/symbol", pprof.Symbol)

	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)

	// Manually add support for paths linked to by index page at /debug/pprof/
	r.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	r.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	r.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	r.Handle("/debug/pprof/block", pprof.Handler("block"))
}

func Index(w http.ResponseWriter, r *http.Request) {
	for i := 0; i < 1000000; i++ {
		math.Pow(36, 89)
	}
	fmt.Fprintf(w, "Hello, World")
}

func setProfile() {
	//For profiling
	//runtime.MemProfileRate = 1
	runtime.SetBlockProfileRate(1)

}

//pattern 1 for handler
func createMux() *http.ServeMux {
	//http.HandleFunc("/", handler) // ハンドラを登録してウェブページを表示させる
	//->DefaultServerMuxと言うものにマッピングが登録される
	/*
		func HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
			DefaultServeMux.HandleFunc(pattern, handler)
		}
	*/
	//http.DefaultMaxIdleConnsPerHost: 2
	//http.ListenAndServe(":8080", nil)
	//handler が nil のときは、DefaultServeMux が handler として使われる

	//Default Mux
	mux := http.NewServeMux()
	mux.Handle("/list", http.HandlerFunc(Index))
	mux.Handle("/price", http.HandlerFunc(Index))

	return mux
}

//pattern 2 for handler
func handleFunc() {
	// HandleFunc registers the handler function for the given pattern
	// in the DefaultServeMux.
	// The documentation for ServeMux explains how patterns are matched.
	http.HandleFunc("/list", Index)
	http.HandleFunc("/price", Index)
}

func main() {
	//For profiling of test
	setProfile()

	//Setting Server
	r := mux.NewRouter()
	setRoute(r)
	http.ListenAndServe(":9999", r)
}
