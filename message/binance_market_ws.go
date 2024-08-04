package message

import (
	binanceSpot "github.com/dictxwang/go-binance"
	binanceFutures "github.com/dictxwang/go-binance/futures"
	"math/rand"
	"okxdata/config"
	"okxdata/context"
	"okxdata/utils/logger"
	"time"
)

func StartBinanceMarketWs(cfg *config.Config, globalContext *context.GlobalContext,
	binanceFuturesTickerChan chan *binanceFutures.WsBookTickerEvent) {

	binanceWs := BinanceMarketWebSocket{}
	binanceWs.Init(binanceFuturesTickerChan)
	binanceWs.StartBinanceMarketWs(cfg, globalContext)
}

type BinanceMarketWebSocket struct {
	spotTickerChan    chan *binanceSpot.WsBookTickerEvent
	futuresTickerChan chan *binanceFutures.WsBookTickerEvent
}

func (ws *BinanceMarketWebSocket) Init(futuresTickerChan chan *binanceFutures.WsBookTickerEvent) {
	ws.futuresTickerChan = futuresTickerChan
}

func (ws *BinanceMarketWebSocket) StartBinanceMarketWs(cfg *config.Config, globalContext *context.GlobalContext) {

	batchSize := 30
	instIDs := globalContext.InstrumentComposite.BinanceInstIDs

	instIDsSize := len(instIDs)
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
	logger.Info("[BSTickerWebSocket] Start Listen Binance Spot Tickers")
	logger.Info("[BFTickerWebSocket] Start Listen Binance Futures Tickers")
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
