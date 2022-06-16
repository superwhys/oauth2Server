package server

import (
	"github.com/gin-gonic/gin"
	ginsession "github.com/go-session/gin-session"
	oauthModel "github.com/superwhys/oauth2Server/models"
	"github.com/superwhys/superGo/superLog"
	"gopkg.in/oauth2.v3/models"
	"net/http"
	"net/url"
)

func (s *Oauth2Server) loginPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

func (s *Oauth2Server) setLoginForm(values url.Values, username string) {
	s.clientStore.Set(&models.Client{
		ID:     values.Get("client_id"),
		Secret: values.Get("client_secret"),
		Domain: values.Get("redirect_uri"),
		UserID: username,
	})
}

func (s *Oauth2Server) loginHandler(c *gin.Context) {
	store := ginsession.FromContext(c)

	var userForm url.Values
	if v, ok := store.Get("userForm"); ok {
		userForm = parseForm(v)
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user := &oauthModel.LoginInfo{}
	err := c.ShouldBind(user)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	s.setLoginForm(userForm, user.Username)

	superLog.Debugf("login user info: %v", user)
	store.Set("LoggedInUserID", user.Username)
	err = store.Save()
	if err != nil {
		superLog.Error("store userid error: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// redirect to /auth
	c.Redirect(http.StatusMovedPermanently, "/oauth2/auth")
	return
}
