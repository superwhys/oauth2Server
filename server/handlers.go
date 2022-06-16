package server

import (
	"github.com/go-session/session"
	"github.com/superwhys/superGo/superLog"
	"gopkg.in/oauth2.v3/errors"
	"log"
	"net/http"
)

func UserAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		superLog.Error("get store error: %v", err)
		return
	}

	uid, ok := store.Get("LoggedInUserID")
	if !ok {
		superLog.Info("user not login, redirect to login page")
		if r.Form == nil {
			r.ParseForm()
		}

		store.Set("userForm", r.Form)
		err = store.Save()
		if err != nil {
			superLog.Error("store save form error: %v", err)
			return "", err
		}
		w.Header().Set("Location", "/oauth2/login")
		w.WriteHeader(http.StatusFound)
		return
	}
	userID = uid.(string)
	store.Delete("LoggedInUserID")
	store.Save()
	superLog.Infof("userId: %v", userID)
	superLog.Infof("err: %v", err)
	return
}

func InterErrorHandler(err error) (re *errors.Response) {
	log.Println("Internal Error:", err.Error())
	return
}

func ResponseErrorHandler(re *errors.Response) {
	log.Println("Response Error:", re.Error.Error())
}
