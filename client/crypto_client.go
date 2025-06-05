package client

import (
	"context"
	"github.com/drinkthere/cryptodotcom"
	"github.com/drinkthere/cryptodotcom/api"
	"golang.org/x/time/rate"
	"log"
	"okxdata/config"
)

type CryptoClient struct {
	Client       *api.Client
	limiter      *rate.Limiter
	limitProcess int
}

func (c *CryptoClient) Init(cfg *config.Config) bool {
	ctx := context.Background()

	client, err := api.NewClient(ctx, cfg.CryptoAPIKey, cfg.CryptoSecretKey, cryptodotcom.NormalServer, cfg.CryptoLocalIP)
	if err != nil {
		log.Fatal(err)
		return false
	}

	c.Client = client
	return true
}
