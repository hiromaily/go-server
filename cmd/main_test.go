package main

import (
	"fmt"
	//"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	//"reflect"
	"github.com/hiromaily/go-server/libs/web"
	"testing"
	"time"
)

var ts *httptest.Server
var getTests = []struct {
	url    string
	code   int
	method string
	err    bool //when err is true, set wrong response to occur error intentionally
}{
	{"/", http.StatusOK, "GET", false},
	{"/index", http.StatusOK, "GET", false},
	{"/login", http.StatusInternalServerError, "GET", false},
	{"/global", http.StatusOK, "GET", false},
	{"/webpush", http.StatusOK, "GET", false},
}

//w.Get("/", c.GetIndex)
//w.Get("/index", c.GetIndex)
//w.Get("/login", c.GetLogin)
//w.Post("/login", c.PostLogin)
//w.Get("/global", c.GetGlobal)
//w.Get("/webpush", c.GetWebPush)
//w.Post("/webpush", c.PostWebPush)

//-----------------------------------------------------------------------------
// Test Framework
//-----------------------------------------------------------------------------
// Initialize
func init() {
	//flag.Parse()
}

func setup() {

	w := web.New()
	setMiddleware(w)
	setRoute(w)

	//mock server
	//type Handler interface {
	//	ServeHTTP(ResponseWriter, *Request)
	//}
	w.Mux.Handle("/", http.HandlerFunc(w.Handler))
	ts = httptest.NewServer(w.Mux)
}

func teardown() {
	ts.Close()
}

func TestMain(m *testing.M) {
	setup()

	code := m.Run()

	teardown()

	os.Exit(code)
}

//-----------------------------------------------------------------------------
// functions
//-----------------------------------------------------------------------------
func setParams(req *http.Request, params []string) {
	if params == nil || len(params) == 0 {
		return
	}

	q := req.URL.Query()
	for _, v := range params {
		q.Add("u", v)
	}
	req.URL.RawQuery = q.Encode()
}

//-----------------------------------------------------------------------------
// Test
//-----------------------------------------------------------------------------
func TestGetRequestOnTable(t *testing.T) {
	//request
	client := &http.Client{
		Timeout: time.Duration(500) * time.Millisecond,
	}

	for i, tt := range getTests {
		fmt.Printf("%d [%s: %s]\n", i+1, tt.method, tt.url)

		req, _ := http.NewRequest(tt.method, ts.URL+tt.url, nil)

		//setParams(req, tt.params)

		resp, _ := client.Do(req)
		//body, _ := ioutil.ReadAll(resp.Body)
		//res, _ := getNumbers(body)
		//fmt.Println(body)

		//compare result
		if tt.code != resp.StatusCode {
			t.Errorf("%d [%s:%s] response is not correct. \n return code is %d / expected %d", i,
				tt.method, tt.url, resp.StatusCode, tt.code)
		}

		resp.Body.Close()
	}
}
