package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/go-oauth2/mongo"
	"github.com/go-session/session"

	oauthMongo "github.com/superwhys/oauth2Server/db/mongo"
	"github.com/superwhys/oauth2Server/db/redis"
	"github.com/superwhys/oauth2Server/server"
	"github.com/superwhys/superGo/superFlags"
	"gopkg.in/oauth2.v3/manage"
	"time"
)

var (
	port      = superFlags.Int("port", 9915, "run port of service")
	redisAddr = superFlags.String("redisAddr", "localhost:6379", "redis server address")
	mongoAddr = superFlags.String("mongoAddr", "localhost:27017", "mongo server address")
	secretKey = superFlags.String("secretKey", "fe8a711ed3fcc1ba9e56d35369ebc589bc420d6bf474eed24b878b7a09e9ed96", "key of jwt token")
)

func main() {
	sessionStoreOption := server.InitMongoStore(mongoAddr(), "oauth2", "sessions")
	session.InitManager(sessionStoreOption)

	redisPool := redis.DialRedisPoolBlocked(redisAddr(), 1, 20, time.Minute*3)
	defer redisPool.Close()

	clientMongo := mongo.NewClientStore(mongo.NewConfig(mongoAddr(), "client"))
	tokenMongo := mongo.NewTokenStore(mongo.NewConfig(mongoAddr(), "token"))

	oauth2Manger := server.NewManager(
		server.WithTokenStorage(tokenMongo),
		server.WithClientStorage(clientMongo),
		server.WithAuthCodeTokenConf(manage.DefaultAuthorizeCodeTokenCfg),
		server.WithJWTAccess([]byte(secretKey()), jwt.SigningMethodHS256),
		server.WithUserAuthHandler(server.UserAuthorizeHandler),
		server.WithInternalErrorHandler(server.InterErrorHandler),
		server.WithResponseErrorHandler(server.ResponseErrorHandler),
	)

	mongoCli := oauthMongo.DialMongo(mongoAddr())

	srv := server.NewOauth2Server(redisPool, oauth2Manger, mongoCli, tokenMongo, clientMongo)

	srv.SetUpRouter(sessionStoreOption)
	srv.Run(port())
}
