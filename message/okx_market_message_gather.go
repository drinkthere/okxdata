package message

import (
	"github.com/drinkthere/okx/events/public"
	"github.com/drinkthere/okx/models/market"
	"math/rand"
	"okxdata/config"
	"okxdata/container"
	"okxdata/context"
	"okxdata/utils"
	"okxdata/utils/logger"
)

func convertToOkxTickerMessage(ticker *market.Ticker) container.TickerMessage {
	return container.TickerMessage{
		Exchange: config.OkxExchange,
		InstType: utils.ConvertToStdInstType(config.OkxExchange, string(ticker.InstType)),
		InstID:   ticker.InstID,
		AskPx:    float64(ticker.AskPx),
		AskSz:    float64(ticker.AskSz),
		BidPx:    float64(ticker.BidPx),
		BidSz:    float64(ticker.BidSz),
	}
}

func StartGatherOkxFuturesTicker(tickChan chan *public.Tickers, globalConfig *config.Config,
	globalContext *context.GlobalContext) {

	r := rand.New(rand.NewSource(2))
	go func() {
		defer func() {
			if rc := recover(); rc != nil {
				logger.Error("[GatherFTicker] Recovered from panic: %v", rc)
			}

			logger.Warn("[GatherFTicker] Okx Swap Ticker Gather Exited.")
		}()
		for {
			s := <-tickChan
			for _, t := range s.Tickers {
				tickerMsg := convertToOkxTickerMessage(t)
				globalContext.PriceComposite.UpdatePriceList(tickerMsg, globalConfig)
			}
			if r.Int31n(10000) < 5 && len(s.Tickers) > 0 {
				logger.Info("[Gather] Receive Okx Futures Ticker %+v", s.Tickers[0])
			}
		}
	}()

	logger.Info("[Gather] Start Gather Okx Futures Ticker")
}

func StartGatherOkxSpotTicker(tickChan chan *public.Tickers, globalConfig *config.Config,
	globalContext *context.GlobalContext) {

	r := rand.New(rand.NewSource(2))
	go func() {
		defer func() {
			if rc := recover(); rc != nil {
				logger.Error("[GatherSTicker] Recovered from panic: %v", rc)
			}

			logger.Warn("[GatherSTicker] Okx Spot Ticker Gather Exited.")
		}()
		for {
			s := <-tickChan
			for _, t := range s.Tickers {
				tickerMsg := convertToOkxTickerMessage(t)
				globalContext.PriceComposite.UpdatePriceList(tickerMsg, globalConfig)
			}
			if r.Int31n(10000) < 5 && len(s.Tickers) > 0 {
				logger.Info("[GatherSTicker] Receive Okx Spot Ticker %+v\n", s.Tickers[0])
			}
		}
	}()

	logger.Info("[GatherSTicker] Start Gather Okx Spot Ticker")
}
