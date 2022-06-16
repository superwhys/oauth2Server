package server

import (
	"github.com/gin-gonic/gin"
	ginSession "github.com/go-session/gin-session"
	"github.com/go-session/session"
)

func GetAndPostHandler(group *gin.RouterGroup, path string, handlerFunc ...gin.HandlerFunc) {
	group.GET(path, handlerFunc...)
	group.POST(path, handlerFunc...)
}

func (s *Oauth2Server) SetUpRouter(sOpt session.Option) {
	Router := s.Router
	Router.Use(ginSession.New(sOpt))
	Router.LoadHTMLGlob("static/*")

	oauthGroup := Router.Group("/oauth2")
	{
		oauthGroup.GET("/login", s.loginPageHandler)
		oauthGroup.POST("/login", s.loginHandler)
		oauthGroup.GET("/auth", s.authHandler)
		GetAndPostHandler(oauthGroup, "/authorize", s.authorizeHandler)
		oauthGroup.POST("/token", s.tokenHandler)
	}
}
