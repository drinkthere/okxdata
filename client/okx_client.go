package client

import (
	"context"
	"github.com/drinkthere/okx"
	"github.com/drinkthere/okx/api"
	"golang.org/x/time/rate"
	"log"
	"okxdata/config"
)

type OkxClient struct {
	Client       *api.Client
	limiter      *rate.Limiter
	limitProcess int
}

func (okxClient *OkxClient) Init(cfg *config.Config) bool {
	dest := okx.NormalServer
	ctx := context.Background()
	client, err := api.NewClient(ctx, cfg.OkxAPIKey, cfg.OkxSecretKey, cfg.OkxPassword, dest)

	if err != nil {
		log.Fatal(err)
		return false
	}

	okxClient.Client = client
	return true
}
