package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/pprof"
)

//Go言語でhttpサーバーを立ち上げてHello Worldをする
//http://qiita.com/taizo/items/bf1ec35a65ad5f608d45

//net/http の動きを少しだけ追ってみた - Golang
//http://m0t0k1ch1st0ry.com/blog/2014/06/09/golang/

//Goでnet/httpを使う時のこまごまとした注意
//Goでnet/httpを使う時のこまごまとした注意

//GoでHTTPサーバを立てる
//http://qiita.com/kkyouhei/items/8ce72bf997fa353b7646

func AttachProfiler(router *mux.Router) {
	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)

	// Manually add support for paths linked to by index page at /debug/pprof/
	router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	router.Handle("/debug/pprof/block", pprof.Handler("block"))
}
func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World")
}

func main() {
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

	r := mux.NewRouter()
	AttachProfiler(r)
	r.HandleFunc("/", Handler)

	http.ListenAndServe(":8080", r)
	//handler が nil のときは、DefaultServeMux が handler として使われる
}
