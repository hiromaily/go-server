package controller

import (
	"net/http"

	ss "github.com/hiromaily/go-server/pkg/session"
	tm "github.com/hiromaily/go-server/pkg/template"
	lg "github.com/hiromaily/golibs/log"
	u "github.com/hiromaily/golibs/utils"
)

// GetLogin is for /login page
func GetLogin(res http.ResponseWriter, req *http.Request) {
	lg.Info("[GetLogin]")
	//lg.Debugf("[req]%+v\n", req)

	//fmt.Fprintf(res, "User Login")
	//ctx := req.Context()

	// TODO:check sssion for login
	userID, err := ss.Check(req)
	//if err != nil {
	//	lg.Error(err)
	//	http.Error(res, err.Error(), 500)
	//}
	if err == nil && userID != 0 {
		//redirect
		http.Redirect(res, req, "/global", http.StatusTemporaryRedirect) //307
		return
	}

	//index
	tm.Execute(res, "login", nil)
}

// PostLogin is for login submit from login page
func PostLogin(res http.ResponseWriter, req *http.Request) {
	lg.Info("[PostLogin]")

	//TODO: check user is stored in database or not

	//login NG
	//tm.Execute(res, "login", nil)

	//login OK

	//session
	var dummyUserID = u.GenerateRandom(1, 99999)
	ss.Generate(res, req, dummyUserID)

	//redirect (TODO:needed Get to Post)
	//FIXME:Browser request cache data when redirecting at status code 301
	//301 Moved Permanently   (Do cache,   it's possible to change from POST to GET)
	//302 Found               (Not cache,  it's possible to change from POST to GET)
	//307 Temporary Redirect  (Not cache,  it's not possible to change from POST to GET)
	//308 Moved Permanently   (Do cache,   it's not possible to change from POST to GET)
	http.Redirect(res, req, "/global", http.StatusFound) //302
}
