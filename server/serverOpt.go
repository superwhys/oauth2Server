package server

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/go-oauth2/mongo"
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/generates"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/server"
	"net/http"
)

type oauth2Manage struct {
	manager *manage.Manager
	server  *server.Server
}

type OauthOption func(om *oauth2Manage)

func NewManager(opts ...OauthOption) *oauth2Manage {
	manager := manage.NewDefaultManager()
	srv := server.NewDefaultServer(manager)

	oauth2Manager := &oauth2Manage{
		manager: manager,
		server:  srv,
	}
	for _, opt := range opts {
		opt(oauth2Manager)
	}
	return oauth2Manager
}

func WithAuthCodeTokenConf(conf *manage.Config) OauthOption {
	return func(m *oauth2Manage) {
		m.manager.SetAuthorizeCodeTokenCfg(conf)
	}
}

func WithTokenStorage(store *mongo.TokenStore) OauthOption {
	return func(m *oauth2Manage) {
		m.manager.MapTokenStorage(store)
	}
}

func WithClientStorage(store *mongo.ClientStore) OauthOption {
	return func(m *oauth2Manage) {
		m.manager.MapClientStorage(store)
	}
}

func WithBasicAccess() OauthOption {
	return func(m *oauth2Manage) {
		m.manager.MapAccessGenerate(generates.NewAccessGenerate())
	}
}

func WithJWTAccess(key []byte, method jwt.SigningMethod) OauthOption {
	return func(m *oauth2Manage) {
		m.manager.MapAccessGenerate(generates.NewJWTAccessGenerate(key, method))
	}
}

func WithPasswordAuthHandler(handler func(username, password string) (userID string, err error)) OauthOption {
	return func(m *oauth2Manage) {
		m.server.SetPasswordAuthorizationHandler(handler)
	}
}

func WithUserAuthHandler(handler func(w http.ResponseWriter, r *http.Request) (userID string, err error)) OauthOption {
	return func(m *oauth2Manage) {
		m.server.SetUserAuthorizationHandler(handler)
	}
}

func WithInternalErrorHandler(handler func(err error) (re *errors.Response)) OauthOption {
	return func(m *oauth2Manage) {
		m.server.SetInternalErrorHandler(handler)
	}
}

func WithResponseErrorHandler(handler func(re *errors.Response)) OauthOption {
	return func(m *oauth2Manage) {
		m.server.SetResponseErrorHandler(handler)
	}
}
