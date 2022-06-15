package server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	ginsession "github.com/go-session/gin-session"
	"github.com/superwhys/superGo/superLog"
	"net/http"
	"net/url"
)

func (s *Oauth2Server) authorizeHandler(c *gin.Context) {
	store := ginsession.FromContext(c)

	var userForm url.Values
	if v, ok := store.Get("userForm"); ok {
		vMar, _ := json.Marshal(v)
		values := map[string][]string{}
		json.Unmarshal(vMar, &values)
		userForm = values
	}
	c.Request.Form = userForm
	superLog.Info("1")
	store.Delete("userForm")
	err := store.Save()
	if err != nil {
		superLog.Error("save userForm error: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	superLog.Info(2)
	err = s.oauth2Manager.server.HandleAuthorizeRequest(c.Writer, c.Request)
	if err != nil {
		superLog.Error("handle authorize request error: %v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
}
