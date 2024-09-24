package message

import (
	"github.com/hirokisan/bybit/v2"
	"okxdata/config"
	"okxdata/container"
	"okxdata/context"
	"okxdata/utils"
	"okxdata/utils/logger"
	"strconv"
)

func StartGatherBybitLinearTicker(linearTickerChan chan *bybit.V5WebsocketPublicTickerResponse, globalConfig *config.Config, globalContext *context.GlobalContext) {
	go func() {
		defer func() {
			logger.Error("[GatherBybitLTicker] Bybit Linear Ticker Gather Exited.")
		}()
		for t := range linearTickerChan {
			tickerInfo := t.Data.LinearInverse
			if !utils.InArray(string(tickerInfo.Symbol), globalContext.InstrumentComposite.BybitLinearInstIDs) {
				continue
			}
			tickerMsg := convertBybitLinearTickerEventToTickerMessage(tickerInfo)
			globalContext.PriceComposite.UpdatePriceList(tickerMsg, globalConfig)

		}
	}()

	logger.Info("[GatherBybitLTicker] Start Gather Bybit Spot Tickers")
}

func convertBybitLinearTickerEventToTickerMessage(ticker *bybit.V5WebsocketPublicTickerLinearInverseResult) container.TickerMessage {
	instID := string(ticker.Symbol)

	bidPx, _ := strconv.ParseFloat(ticker.Bid1Price, 64)
	bidSz, _ := strconv.ParseFloat(ticker.Bid1Size, 64)
	askPx, _ := strconv.ParseFloat(ticker.Ask1Price, 64)
	askSz, _ := strconv.ParseFloat(ticker.Ask1Size, 64)

	return container.TickerMessage{
		Exchange: config.BybitExchange,
		InstType: config.LinearInstrument,
		InstID:   instID,
		AskPx:    askPx,
		AskSz:    askSz,
		BidPx:    bidPx,
		BidSz:    bidSz,
	}
}
