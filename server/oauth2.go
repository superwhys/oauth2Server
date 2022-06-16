package server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	ginsession "github.com/go-session/gin-session"
	"github.com/superwhys/superGo/superLog"
	"net/http"
	"net/url"
)

func parseForm(form interface{}) url.Values {
	vMar, _ := json.Marshal(form)
	values := map[string][]string{}
	json.Unmarshal(vMar, &values)
	return values
}

func (s *Oauth2Server) authorizeHandler(c *gin.Context) {
	store := ginsession.FromContext(c)

	var userForm url.Values
	if v, ok := store.Get("userForm"); ok {
		userForm = parseForm(v)
	}
	c.Request.Form = userForm

	store.Delete("userForm")
	err := store.Save()
	if err != nil {
		superLog.Error("save userForm error: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = s.oauth2Manager.server.HandleAuthorizeRequest(c.Writer, c.Request)
	if err != nil {
		superLog.Error("handle authorize request error: %v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
}

func (s *Oauth2Server) tokenHandler(c *gin.Context) {
	err := s.oauth2Manager.server.HandleTokenRequest(c.Writer, c.Request)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
}
