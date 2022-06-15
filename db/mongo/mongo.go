package mongo

import (
	"context"

	"github.com/globalsign/mgo/bson"
	"github.com/superwhys/oauth2Server/models"
	"github.com/superwhys/superGo/superLog"
	"github.com/superwhys/superGo/superMongo"
)

type MongoCli struct {
	Client *superMongo.Client
}

func DialMongo(addr string) *MongoCli {
	cli := superMongo.NewClient(addr)
	return &MongoCli{Client: cli}
}

func (c *MongoCli) SetClient(ctx context.Context, info *models.TokenInfo) error {
	con := c.Client.OpenWithContext(ctx, "oauth2", "clientInfo")
	defer con.Close()
	return con.Insert(info)
}

func (c *MongoCli) VerifySecret(ctx context.Context, secret string) (bool, error) {
	con := c.Client.OpenWithContext(ctx, "oauth2", "clientSecret")
	defer con.Close()

	superLog.Info("find secret: ", secret)
	secretInfo := &models.SecretInfo{}
	err := con.Find(bson.M{"secret": secret}).One(secretInfo)
	if err != nil {
		return false, err
	}
	return secretInfo.IsBlock, nil
}

func (c *MongoCli) SetToken(ctx context.Context, clientId, token string) error {
	con := c.Client.OpenWithContext(ctx, "oauth2", "clientInfo")
	defer con.Close()

	return con.Update(bson.M{"client_id": clientId}, bson.M{"token": token})
}
