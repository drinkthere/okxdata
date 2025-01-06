package message

import (
	"context"
	"github.com/hirokisan/bybit/v2"
	"okxdata/config"
	mmContext "okxdata/context"
	"okxdata/utils/logger"
	"time"
)

func StartBybitMarketWs(
	cfg *config.Config,
	globalContext *mmContext.GlobalContext,
	linearTickerChan chan *bybit.V5WebsocketPublicTickerResponse,
	spotTickerChan chan *bybit.V5WebsocketPublicTickerResponse) {

	bybitMarketWs := newBybitMarketWebsocket(linearTickerChan, spotTickerChan)
	bybitMarketWs.Start(cfg, globalContext)
}

type BybitMarketWebSocket struct {
	linearTickerChan     chan *bybit.V5WebsocketPublicTickerResponse
	spotTickerChan       chan *bybit.V5WebsocketPublicTickerResponse
	isLinerTickerStopped bool
	isSpotTickerStopped  bool
}

func newBybitMarketWebsocket(
	linearTickerChan chan *bybit.V5WebsocketPublicTickerResponse,
	spotTickerChan chan *bybit.V5WebsocketPublicTickerResponse) *BybitMarketWebSocket {

	return &BybitMarketWebSocket{
		linearTickerChan:     linearTickerChan,
		spotTickerChan:       spotTickerChan,
		isLinerTickerStopped: true,
		isSpotTickerStopped:  true,
	}
}

func (ws *BybitMarketWebSocket) Start(cfg *config.Config, globalContext *mmContext.GlobalContext) {
	ws.startBybitLinearTickers(cfg, globalContext)
	logger.Info("[BybitLTickerWebSocket] Start Listen Bybit Linear Tickers")

	ws.startBybitSpotTickers(cfg, globalContext)
	logger.Info("[BybitSTickerWebSocket] Start Listen Bybit Spot Tickers")
}

func (ws *BybitMarketWebSocket) handleLinerTickerEvent(event bybit.V5WebsocketPublicTickerResponse) error {
	ws.linearTickerChan <- &event
	return nil
}

func (ws *BybitMarketWebSocket) handleLinerTickerError(isWebsocketClosed bool, err error) {
	if err != nil {
		logger.Error("[BybitLinerTickerWebSocket] Bybit Liner Ticker Error: %+v", err)
	}
	logger.Warn("[BybitLinerTickerWebSocket] Bybit Liner Tickers Ws Will Reconnect In 1 Second")
	time.Sleep(time.Second * 1)
	ws.isLinerTickerStopped = true
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
				time.Sleep(time.Millisecond * 100)
			}

			go svc.Start(context.Background(), ws.handleLinerTickerError)

			logger.Info("[BybitLinearTickerWebSocket] Subscribe Bybit Linear Tickers: %+v", globalContext.InstrumentComposite.BybitLinearInstIDs)
			ws.isLinerTickerStopped = false
		}
	}()
}

func (ws *BybitMarketWebSocket) handleSpotTickerEvent(event bybit.V5WebsocketPublicTickerResponse) error {
	ws.spotTickerChan <- &event
	return nil
}

func (ws *BybitMarketWebSocket) handleSpotTickerError(isWebsocketClosed bool, err error) {
	if err != nil {
		logger.Error("[BybitSpotTickerWebSocket] Bybit Spot Ticker Error: %+v", err)
	}
	logger.Warn("[BybitMarketWs] Bybit Spot Tickers Ws Will Reconnect In 1 Second")
	time.Sleep(time.Second * 1)
	ws.isSpotTickerStopped = true
}

func (ws *BybitMarketWebSocket) startBybitSpotTickers(cfg *config.Config, globalContext *mmContext.GlobalContext) {

	go func() {
		defer func() {
			logger.Warn("[BybitSpotTickerWebSocket] Bybit Spot Tickers Listening Exited.")
		}()
		for {
			if !ws.isSpotTickerStopped {
				time.Sleep(time.Second * 1)
				continue
			}

			wsClient := bybit.NewWebsocketClient()
			svc, err := wsClient.V5().Public(bybit.CategoryV5Spot)
			if err != nil {
				logger.Error("[BybitSpotTickerWebSocket] Start Bybit Spot Websocket Failed, error: %+v", err)
				time.Sleep(time.Second * 1)
				continue
			}

			for _, instID := range globalContext.InstrumentComposite.BybitSpotInstIDs {
				_, err = svc.SubscribeTicker(bybit.V5WebsocketPublicTickerParamKey{
					Symbol: bybit.SymbolV5(instID),
				}, ws.handleSpotTickerEvent)

				if err != nil {
					logger.Error("[BybitSpotTickerWebSocket] Subscribe %s Ticker Failed, error: %+v", instID, err)
					time.Sleep(time.Second * 10)
					continue
				}
				time.Sleep(time.Millisecond * 100)
			}

			go svc.Start(context.Background(), ws.handleSpotTickerError)

			logger.Info("[BybitSpotTickerWebSocket] Subscribe Bybit Spot Tickers: %+v", globalContext.InstrumentComposite.BybitSpotInstIDs)
			ws.isLinerTickerStopped = false
		}
	}()
}
