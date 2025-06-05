package client

import (
	"context"
	"github.com/drinkthere/cryptodotcom"
	"github.com/drinkthere/cryptodotcom/api"
	"golang.org/x/time/rate"
	"log"
	"okxdata/config"
	"okxdata/utils/logger"
	"time"
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
	limit := rate.Every(1 * time.Second / time.Duration(cfg.APILimit))
	c.limiter = rate.NewLimiter(limit, 60)
	c.limitProcess = cfg.LimitProcess
	return true
}

func (c *CryptoClient) CheckLimit(n int) bool {
	if c.limitProcess == 1 {
		err := c.limiter.WaitN(context.Background(), n)
		if err != nil {
			logger.Error("[CryptoClient] reach to limit, error:%s", err.Error())
		}
		return true
	}
	ret := c.limiter.AllowN(time.Now(), n)
	if !ret {
		logger.Warn("[CryptoClient] reach to limit")
	}
	return ret
}
