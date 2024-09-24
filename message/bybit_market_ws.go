package message

import (
	"context"
	"github.com/hirokisan/bybit/v2"
	"okxdata/config"
	mmContext "okxdata/context"
	"okxdata/utils/logger"
	"time"
)

func StartBybitMarketWs(cfg *config.Config, globalContext *mmContext.GlobalContext, linearTickerChan chan *bybit.V5WebsocketPublicTickerResponse) {

	bybitMarketWs := newBybitMarketWebsocket(linearTickerChan)
	bybitMarketWs.Start(cfg, globalContext)
}

type BybitMarketWebSocket struct {
	linearTickerChan     chan *bybit.V5WebsocketPublicTickerResponse
	isLinerTickerStopped bool
}

func newBybitMarketWebsocket(linearTickerChan chan *bybit.V5WebsocketPublicTickerResponse) *BybitMarketWebSocket {

	return &BybitMarketWebSocket{
		linearTickerChan:     linearTickerChan,
		isLinerTickerStopped: true,
	}
}

func (ws *BybitMarketWebSocket) Start(cfg *config.Config, globalContext *mmContext.GlobalContext) {
	ws.startBybitLinearTickers(cfg, globalContext)
	logger.Info("[BybitLTickerWebSocket] Start Listen Bybit Linear Tickers")
}

func (ws *BybitMarketWebSocket) handleLinerTickerEvent(event bybit.V5WebsocketPublicTickerResponse) error {
	ws.linearTickerChan <- &event
	return nil
}

func (ws *BybitMarketWebSocket) handleLinerTickerError(isWebsocketClosed bool, err error) {
	if err != nil {
		logger.Error("[BybitLinerTickerWebSocket] Bybit Liner Ticker Error: %+v", err)
	}
	if isWebsocketClosed {
		ws.isLinerTickerStopped = false
	}
}

func (ws *BybitMarketWebSocket) startBybitLinearTickers(cfg *config.Config, globalContext *mmContext.GlobalContext) {

	go func() {
		defer func() {
			logger.Warn("[BybitLinearTickerWebSocket] Bybit Linear Tickers Listening Exited.")
		}()
		for {
			if !ws.isLinerTickerStopped {
				time.Sleep(time.Second * 1)
				continue
			}

			wsClient := bybit.NewWebsocketClient()
			svc, err := wsClient.V5().Public(bybit.CategoryV5Linear)
			if err != nil {
				logger.Error("[BybitLinearTickerWebSocket] Start Bybit Linear Websocket Failed, error: %+v", err)
				time.Sleep(time.Second * 1)
				continue
			}

			for _, instID := range globalContext.InstrumentComposite.BybitLinearInstIDs {
				_, err = svc.SubscribeTicker(bybit.V5WebsocketPublicTickerParamKey{
					Symbol: bybit.SymbolV5(instID),
				}, ws.handleLinerTickerEvent)

				if err != nil {
					logger.Error("[BybitLinearTickerWebSocket] Subscribe %s Ticker Failed, error: %+v", instID, err)
					time.Sleep(time.Second * 10)
					continue
				}
			}

			go svc.Start(context.Background(), ws.handleLinerTickerError)

			logger.Info("[BybitLinearTickerWebSocket] Subscribe Bybit Linear Tickers: %+v", globalContext.InstrumentComposite.BybitLinearInstIDs)
			ws.isLinerTickerStopped = false
		}
	}()
}
