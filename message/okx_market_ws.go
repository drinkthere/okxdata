package message

import (
	"github.com/drinkthere/okx/events"
	"github.com/drinkthere/okx/events/public"
	wsRequestPublic "github.com/drinkthere/okx/requests/ws/public"
	"okxdata/client"
	"okxdata/config"
	"okxdata/context"
	"okxdata/utils/logger"
	"time"
)

func StartOkxMarketWs(cfg *config.Config, globalContext *context.GlobalContext,
	okxFuturesTickerChan chan *public.Tickers) {

	startOkxFuturesTickers(cfg, globalContext, okxFuturesTickerChan)
	logger.Info("[WebSocket] Start Listen Okx Futures Tickers")
}

func startOkxFuturesTickers(cfg *config.Config, globalContext *context.GlobalContext, tickerChan chan *public.Tickers) {

	go func() {
		for {
		ReConnect:
			errChan := make(chan *events.Error)
			subChan := make(chan *events.Subscribe)
			uSubChan := make(chan *events.Unsubscribe)
			loginCh := make(chan *events.Login)
			successCh := make(chan *events.Success)

			var okxClient = client.OkxClient{}
			okxClient.Init(cfg)

			okxClient.Client.Ws.SetChannels(errChan, subChan, uSubChan, loginCh, successCh)
			for _, instID := range globalContext.InstrumentComposite.OkxSwapInstIDs {
				err := okxClient.Client.Ws.Public.Tickers(wsRequestPublic.Tickers{
					InstID: instID,
				}, tickerChan)

				if err != nil {
					logger.Fatal("[WebSocket] Fail To Listen Futures Ticker For %s, %s", instID, err.Error())
				}
				logger.Info("[WebSocket] Futures Ticker WebSocket Has Established For %s", instID)
			}

			for {
				select {
				case sub := <-subChan:
					channel, _ := sub.Arg.Get("channel")
					logger.Info("[WebSocket] Futures Subscribe \t%s", channel)
				case err := <-errChan:
					logger.Error("[WebSocket] Futures Occur Some Error \t%+v", err)
					for _, datum := range err.Data {
						logger.Error("[WebSocket] Futures Error Data \t\t%+v", datum)
					}
				case s := <-successCh:
					logger.Info("[WebSocket] Futures Receive Success: %+v", s)
				case b := <-okxClient.Client.Ws.DoneChan:
					logger.Info("[WebSocket] Futures End\t%v", b)
					// 暂停一秒再跳出，避免异常时频繁发起重连
					logger.Warn("[WebSocket] Will Reconnect Futures-WebSocket After 1 Second")
					time.Sleep(time.Second * 1)
					goto ReConnect
				}
			}
		}
	}()
}
