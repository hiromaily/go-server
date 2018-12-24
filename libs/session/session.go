package session

import (
	"net/http"

	"github.com/alexedwards/scs"
	lg "github.com/hiromaily/golibs/log"
)

var sessionManager = scs.NewCookieManager("u46IpCV9y5Vlur8YvODJEhgOY8m9JVE4")

func GetSessionMgr() *scs.Manager {
	return sessionManager
}

// Generate is to generate session
func Generate(w http.ResponseWriter, r *http.Request, userID int) error {
	//session
	session := sessionManager.Load(r)

	// Then make the privilege-level change.
	err := session.PutInt(w, "userID", userID)
	if err != nil {
		return err
	}
	return nil
}

// Check is to check session
func Check(r *http.Request) (int, error) {
	// check session
	session := sessionManager.Load(r)
	userID, err := session.GetInt("userID")
	if err != nil {
		return 0, err
	}

	lg.Debug("[userID]", userID)
	return userID, nil
}
