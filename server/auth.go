package server

import (
	"github.com/gin-gonic/gin"
	ginsession "github.com/go-session/gin-session"
	"net/http"
)

func (s *Oauth2Server) authHandler(c *gin.Context) {
	store := ginsession.FromContext(c)
	if _, ok := store.Get("LoggedInUserID"); !ok {
		c.Redirect(http.StatusMovedPermanently, "/oauth2/login")
		return
	}
	c.HTML(http.StatusOK, "auth.html", gin.H{})
}
