package main

import (
	"fmt"
	"github.com/gorilla/mux"
	h "github.com/hiromaily/go-server/libs/http"
	"math"
	"net/http"
)

func setRoute(r *mux.Router) {
	//static
	r.PathPrefix("/statics/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))

	//r.GET("/", Index)
	r.HandleFunc("/", Index)

	//For debug mode or development
	h.AttachProfiler(r)
}

//pattern 1 for handler (sample)
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

//pattern 2 for handler (sample)
func handleFunc() {
	// HandleFunc registers the handler function for the given pattern
	// in the DefaultServeMux.
	// The documentation for ServeMux explains how patterns are matched.
	http.HandleFunc("/list", Index)
	http.HandleFunc("/price", Index)
}

func Index(w http.ResponseWriter, r *http.Request) {
	for i := 0; i < 1000000; i++ {
		math.Pow(36, 89)
	}
	fmt.Fprintf(w, "Hello, World")
}
