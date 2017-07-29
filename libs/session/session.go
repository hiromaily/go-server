package session

import (
	"net/http"

	"github.com/alexedwards/scs/session"
	lg "github.com/hiromaily/golibs/log"
)

func Generate(r *http.Request, userID int) error {
	err := session.RegenerateToken(r)
	if err != nil {
		return err
	}

	// Then make the privilege-level change.
	err = session.PutInt(r, "userID", userID)
	if err != nil {
		return err
	}
	return nil
}

func Check(r *http.Request) (int, error) {
	// check session
	userID, err := session.GetInt(r, "userID")
	if err != nil {
		return 0, err
	}

	lg.Debug("[userID]", userID)
	return userID, nil
}
