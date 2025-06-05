package message

import (
	"github.com/drinkthere/cryptodotcom/events"
	"github.com/drinkthere/cryptodotcom/events/public"
	wsRequestPublic "github.com/drinkthere/cryptodotcom/requests/ws/public"
	"okxdata/client"
	"okxdata/config"
	"okxdata/context"
	"okxdata/utils/logger"
	"strings"
	"time"
)

func StartCryptoMarketWs(globalConfig *config.Config, globalContext *context.GlobalContext,
	instType config.InstrumentType, tickerChan chan *public.Tickers) {

	logger.Info("[%s%sTickerWs] Start Listen Tickers", config.CryptoExchange, instType)
	go func() {
		defer func() {
			if rc := recover(); rc != nil {
				logger.Error("[%s%sTickerWs] Recovered from panic: %v", config.CryptoExchange, instType, rc)
			}

			logger.Warn("[%s%sTickerWs] Ticker Listening Exited.", config.CryptoExchange, instType)
		}()
		for {
		ReConnect:
			errChan := make(chan *events.Basic)
			loginCh := make(chan *events.Login)
			successCh := make(chan *events.Basic)

			var cli = client.CryptoClient{}
			cli.Init(globalConfig)

			cli.Client.Ws.SetChannels(loginCh, errChan, successCh)
			err := cli.Client.Ws.Connect(false)
			if err != nil {
				logger.Error("Connect Error is %+v", err)
				time.Sleep(time.Second * 30)
				goto ReConnect
			}
			// crypto.com 建议建立连接之后，sleep 1s再发送订阅请求，避免计算限制的时候，触发LIMIT
			time.Sleep(time.Second)
			instIDs := globalContext.InstrumentComposite.CryptoSwapInstIDs
			err = cli.Client.Ws.Public.Tickers(wsRequestPublic.Tickers{
				InstrumentNames: instIDs,
			}, tickerChan)
			if err != nil {
				logger.Fatal("[%s%sTickerWs] Fail To Listen Ticker For %+v, %s", config.CryptoExchange, instType, instIDs, err.Error())
			} else {
				logger.Info("[%s%sTickerWs] Ticker WebSocket Has Established For %+v", config.CryptoExchange, instType, instIDs)
			}

			for {
				select {
				case e := <-errChan:
					if strings.Contains(e.Message, "i/o timeout") {
						logger.Warn("[%s%sTickerWs] Error occurred %s, Will reconnect after 1 second.", config.CryptoExchange, instType, e.Message)
						time.Sleep(time.Second * 1)
						goto ReConnect
					}
					logger.Error("[%s%sTickerWs] Occur Some Error %+v", config.CryptoExchange, instType, err)
				case s := <-successCh:
					logger.Info("[%s%sTickerWs] Receive Success: %+v", config.CryptoExchange, instType, s)
				case b := <-cli.Client.Ws.DoneChan:
					logger.Info("[%s%sTickerWs] End %v", config.CryptoExchange, instType, b)
					// 暂停一秒再跳出，避免异常时频繁发起重连
					logger.Warn("[%s%sTickerWs] Will Reconnect WebSocket After 1 Second", config.CryptoExchange, instType)
					time.Sleep(time.Second * 1)
					goto ReConnect
				}
			}
		}
	}()
}
