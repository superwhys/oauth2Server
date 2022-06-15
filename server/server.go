package server

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/mongo"
	sMongo "github.com/go-session/mongo"
	"github.com/go-session/session"
	oauthMongo "github.com/superwhys/oauth2Server/db/mongo"
	"github.com/superwhys/oauth2Server/db/redis"
	"github.com/superwhys/superGo/superLog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Oauth2Server struct {
	MongoClient   *oauthMongo.MongoCli
	Router        *gin.Engine
	Store         *redis.Store
	oauth2Manager *oauth2Manage
	tokenStore    *mongo.TokenStore
	clientStore   *mongo.ClientStore
}

func NewOauth2Server(
	store *redis.Store,
	oauth2Manager *oauth2Manage,
	mongoClient *oauthMongo.MongoCli,
	tokenStore *mongo.TokenStore,
	clientStore *mongo.ClientStore,
) *Oauth2Server {
	return &Oauth2Server{
		Router:        gin.Default(),
		Store:         store,
		oauth2Manager: oauth2Manager,
		MongoClient:   mongoClient,
		tokenStore:    tokenStore,
		clientStore:   clientStore,
	}
}

func InitMongoStore(mongoAddr, dbName, cName string) session.Option {
	return session.SetStore(sMongo.NewStore(mongoAddr, dbName, cName))
}

func (s *Oauth2Server) Run(port int) {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: s.Router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			superLog.Fatal("listen error: ", err)
		}
	}()

	// wait interrupt signal to close the server gracefully
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	superLog.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// delay 5s
	if err := srv.Shutdown(ctx); err != nil {
		superLog.Fatal("Server Shutdown: ", err)
	}

	superLog.Info("Server exiting")
}
