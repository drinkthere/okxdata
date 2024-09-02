package message

import (
	binanceDelivery "github.com/dictxwang/go-binance/delivery"
	binanceFutures "github.com/dictxwang/go-binance/futures"
	"okxdata/config"
	"okxdata/container"
	"okxdata/context"
	"okxdata/utils/logger"
	"strconv"
)

func StartGatherBinanceDeliveryBookTicker(tickChan chan *binanceDelivery.WsBookTickerEvent, globalConfig *config.Config, globalContext *context.GlobalContext) {
	go func() {
		defer func() {
			logger.Warn("[GatherBDTickerErr] Binance Delivery Ticker Gather Exited.")
		}()
		for t := range tickChan {
			tickerMsg := convertDeliveryEventToBinanceTickerMessage(t)
			globalContext.PriceComposite.UpdatePriceList(tickerMsg, globalConfig)
		}
	}()

	logger.Info("[GatherBFTickerErr] Start Gather Binance Delivery Ticker")
}

func StartGatherBinanceFuturesBookTicker(tickChan chan *binanceFutures.WsBookTickerEvent, globalConfig *config.Config, globalContext *context.GlobalContext) {
	go func() {
		defer func() {
			logger.Warn("[GatherBFTickerErr] Binance Futures Ticker Gather Exited.")
		}()
		for t := range tickChan {
			tickerMsg := convertFuturesEventToBinanceTickerMessage(t)
			globalContext.PriceComposite.UpdatePriceList(tickerMsg, globalConfig)
		}
	}()

	logger.Info("[GatherBFTickerErr] Start Gather Binance Futures Flash Ticker")
}

func convertDeliveryEventToBinanceTickerMessage(ticker *binanceDelivery.WsBookTickerEvent) container.TickerMessage {
	bestAskPrice, _ := strconv.ParseFloat(ticker.BestAskPrice, 64)
	bestAskQty, _ := strconv.ParseFloat(ticker.BestAskQty, 64)
	bestBidPrice, _ := strconv.ParseFloat(ticker.BestBidPrice, 64)
	bestBidQty, _ := strconv.ParseFloat(ticker.BestBidQty, 64)
	return container.TickerMessage{
		Exchange: config.BinanceExchange,
		InstType: config.DeliveryInstrument,
		InstID:   ticker.Symbol,
		AskPx:    bestAskPrice,
		AskSz:    bestAskQty,
		BidPx:    bestBidPrice,
		BidSz:    bestBidQty,
	}
}

func convertFuturesEventToBinanceTickerMessage(ticker *binanceFutures.WsBookTickerEvent) container.TickerMessage {
	bestAskPrice, _ := strconv.ParseFloat(ticker.BestAskPrice, 64)
	bestAskQty, _ := strconv.ParseFloat(ticker.BestAskQty, 64)
	bestBidPrice, _ := strconv.ParseFloat(ticker.BestBidPrice, 64)
	bestBidQty, _ := strconv.ParseFloat(ticker.BestBidQty, 64)
	return container.TickerMessage{
		Exchange: config.BinanceExchange,
		InstType: config.FuturesInstrument,
		InstID:   ticker.Symbol,
		AskPx:    bestAskPrice,
		AskSz:    bestAskQty,
		BidPx:    bestBidPrice,
		BidSz:    bestBidQty,
	}
}
