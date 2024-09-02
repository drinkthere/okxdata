package message

import (
	binanceDelivery "github.com/dictxwang/go-binance/delivery"
	binanceFutures "github.com/dictxwang/go-binance/futures"
	"math/rand"
	"okxdata/config"
	"okxdata/context"
	"okxdata/utils/logger"
	"time"
)

func StartBinanceMarketWs(cfg *config.Config, globalContext *context.GlobalContext,
	binanceFuturesTickerChan chan *binanceFutures.WsBookTickerEvent, binanceDeliveryTickerChan chan *binanceDelivery.WsBookTickerEvent) {

	binanceWs := newBinanceMarketWebSocket(binanceDeliveryTickerChan, binanceFuturesTickerChan)
	binanceWs.StartBinanceMarketWs(globalContext)
}

type BinanceMarketWebSocket struct {
	deliveryTickerChan chan *binanceDelivery.WsBookTickerEvent
	futuresTickerChan  chan *binanceFutures.WsBookTickerEvent
}

func newBinanceMarketWebSocket(
	binanceDeliveryTickerChan chan *binanceDelivery.WsBookTickerEvent,
	binanceFuturesTickerChan chan *binanceFutures.WsBookTickerEvent) *BinanceMarketWebSocket {
	return &BinanceMarketWebSocket{
		deliveryTickerChan: binanceDeliveryTickerChan,
		futuresTickerChan:  binanceFuturesTickerChan,
	}
}

func (ws *BinanceMarketWebSocket) StartBinanceMarketWs(globalContext *context.GlobalContext) {
	batchSize := 30
	instIDs := globalContext.InstrumentComposite.BinanceDeliveryInstIDs
	instIDsSize := len(instIDs)
	for i := 0; i < instIDsSize; i += batchSize {
		end := i + batchSize
		if end > instIDsSize {
			end = instIDsSize
		}
		batch := instIDs[i:end]

		// 初始化币本位行情监控
		innerDelivery := innerBinanceDeliveryWebSocket{
			instIDs:    batch,
			tickerChan: ws.deliveryTickerChan,
			isStopped:  true,
			randGen:    rand.New(rand.NewSource(2)),
		}
		innerDelivery.subscribeBookTickers("", false)
	}
	logger.Info("[BDTickerWebSocket] Start Listen Binance Delivery Tickers")

	instIDs = globalContext.InstrumentComposite.BinanceFuturesInstIDs
	instIDsSize = len(instIDs)
	for i := 0; i < instIDsSize; i += batchSize {
		end := i + batchSize
		if end > instIDsSize {
			end = instIDsSize
		}
		batch := instIDs[i:end]

		// 初始化合约行情监控
		innerFutures := innerBinanceFuturesWebSocket{
			instIDs:    batch,
			tickerChan: ws.futuresTickerChan,
			isStopped:  true,
			randGen:    rand.New(rand.NewSource(2)),
		}
		innerFutures.startBookTickers("")
	}
	logger.Info("[BFTickerWebSocket] Start Listen Binance Futures Tickers")
}

type innerBinanceDeliveryWebSocket struct {
	instIDs    []string
	tickerChan chan *binanceDelivery.WsBookTickerEvent
	isStopped  bool
	stopChan   chan struct{}
	randGen    *rand.Rand
}

func (ws *innerBinanceDeliveryWebSocket) handleTickerEvent(event *binanceDelivery.WsBookTickerEvent) {
	if ws.randGen.Int31n(10000) < 2 {
		logger.Info("[BDTickerWebSocket] Binance Delivery Event: %+v", event)
	}
	ws.tickerChan <- event
}

func (ws *innerBinanceDeliveryWebSocket) handleError(err error) {
	// 出错断开连接，再重连
	logger.Error("[BDTickerWebSocket] Binance Delivery Handle Error And Reconnect Ws: %s", err.Error())
	ws.stopChan <- struct{}{}
	ws.isStopped = true
}

func (ws *innerBinanceDeliveryWebSocket) subscribeBookTickers(ip string, colo bool) {

	go func() {
		defer func() {
			logger.Warn("[BDTickerWebSocket] Binance Delivery BookTicker Listening Exited.")
		}()
		for {
			if !ws.isStopped {
				time.Sleep(time.Second * 1)
				continue
			}

			if colo {
				binanceDelivery.UseIntranet = true
			} else {
				binanceDelivery.UseIntranet = false
			}

			var stopChan chan struct{}
			var err error
			if ip == "" {
				_, stopChan, err = binanceDelivery.WsCombinedBookTickerServe(ws.instIDs, ws.handleTickerEvent, ws.handleError)
			} else {
				_, stopChan, err = binanceDelivery.WsCombinedBookTickerServeWithIP(ip, ws.instIDs, ws.handleTickerEvent, ws.handleError)
			}
			if err != nil {
				logger.Error("[BDTickerWebSocket] Subscribe Binance Margin Book Tickers Error: %s", err.Error())
				time.Sleep(time.Second * 1)
				continue
			}
			logger.Info("[BDTickerWebSocket] Subscribe Binance Margin Book Tickers: %d", len(ws.instIDs))
			// 重置channel和时间
			ws.stopChan = stopChan
			ws.isStopped = false
		}
	}()
}

type innerBinanceFuturesWebSocket struct {
	instIDs    []string
	tickerChan chan *binanceFutures.WsBookTickerEvent
	isStopped  bool
	stopChan   chan struct{}
	randGen    *rand.Rand
}

func (ws *innerBinanceFuturesWebSocket) handleTickerEvent(event *binanceFutures.WsBookTickerEvent) {
	if ws.randGen.Int31n(10000) < 2 {
		logger.Info("[BFTickerWebSocket] Binance Futures Event: %+v", event)
	}
	ws.tickerChan <- event
}

func (ws *innerBinanceFuturesWebSocket) handleError(err error) {
	// 出错断开连接，再重连
	logger.Error("[BFTickerWebSocket] Binance Futures Handle Error And Reconnect Ws: %s", err.Error())
	ws.stopChan <- struct{}{}
	ws.isStopped = true
}

func (ws *innerBinanceFuturesWebSocket) startBookTickers(ip string) {

	go func() {
		defer func() {
			logger.Warn("[BFTickerWebSocket] Binance Futures Flash Ticker Listening Exited.")
		}()
		for {
			if !ws.isStopped {
				time.Sleep(time.Second * 1)
				continue
			}
			var stopChan chan struct{}
			var err error
			if ip == "" {
				_, stopChan, err = binanceFutures.WsCombinedBookTickerServe(ws.instIDs, ws.handleTickerEvent, ws.handleError)
			} else {
				_, stopChan, err = binanceFutures.WsCombinedBookTickerServeWithIP(ip, ws.instIDs, ws.handleTickerEvent, ws.handleError)
			}

			if err != nil {
				logger.Error("[BFTickerWebSocket] Subscribe Binance Futures Book Tickers Error: %s", err.Error())
				time.Sleep(time.Second * 1)
				continue
			}
			logger.Info("[BFTickerWebSocket] Subscribe Binance Futures Book Tickers: %d", len(ws.instIDs))
			// 重置channel和时间
			ws.stopChan = stopChan
			ws.isStopped = false
		}
	}()
}
