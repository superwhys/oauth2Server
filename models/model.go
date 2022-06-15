package models

type LoginInfo struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type SecretInfo struct {
	ClientSecret string `json:"client_secret" bson:"client_secret"`
	IsBlock      bool   `json:"is_block" bson:"is_block"`
}

type TokenInfo struct {
	ClientId     string `json:"client_id" bson:"client_id"`
	ClientSecret string `json:"client_secret" bson:"client_secret"`
	Scope        string `json:"scope" bson:"scope"`
	Token        string `json:"token" bson:"token"`
}
