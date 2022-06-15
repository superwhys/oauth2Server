package server

import (
	"github.com/gin-gonic/gin"
	ginsession "github.com/go-session/gin-session"
	oauthModel "github.com/superwhys/oauth2Server/models"
	"github.com/superwhys/superGo/superLog"
	"net/http"
)

func (s *Oauth2Server) loginHandler(c *gin.Context) {
	store := ginsession.FromContext(c)

	if c.Request.Method == "POST" {
		// post: store userid
		user := &oauthModel.LoginInfo{}
		err := c.ShouldBind(user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		superLog.Debug(user)

		store.Set("LoggedInUserID", user.Username)
		err = store.Save()
		if err != nil {
			superLog.Error("store userid error: %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		if err != nil {
			superLog.Error("json unmarshal error: ", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// redirect to /auth
		c.Redirect(http.StatusMovedPermanently, "/oauth2/auth")
		return
	}
	// get: output html
	c.HTML(http.StatusOK, "login.html", gin.H{})
}
