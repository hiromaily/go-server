package middleware

import (
	"context"
	"errors"
	"github.com/hiromaily/golibs/auth/jwt"
	lg "github.com/hiromaily/golibs/log"
	"net/http"
	"strings"
	"time"
)

func GenerateJWT(userName string) (string, error) {
	ti := time.Now().Add(time.Minute * 60).Unix()
	token, err := jwt.CreateBasicToken(ti, "", userName)
	if err != nil {
		return "", err
	}
	lg.Debugf("jwt token: %s", token)

	return token, nil
}

func CheckJWT() Handler {
	return func(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, *http.Request) {
		lg.Info("[CheckJWT]")

		var err error

		IsAuth := r.Header.Get("Authorization")
		if IsAuth != "" {
			aAry := strings.Split(IsAuth, " ")
			if len(aAry) != 2 {
				err = errors.New("Authorization header is invalid")
			} else {
				if aAry[0] != "Bearer" {
					err = errors.New("Authorization header is invalid")
				} else {
					token := aAry[1]
					err = jwt.JudgeJWT(token)
				}
			}
		} else {
			err = errors.New("Authorization header was missed.")
		}

		//TODO: how to handle when error occur??
		if err != nil {
			lg.Error(err)
		}

		return w, r
	}
}

//
//func ValidationJWT() error {
//
//}
