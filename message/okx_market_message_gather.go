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
		for {
			s := <-tickChan
			for _, t := range s.Tickers {
				tickerMsg := convertToOkxTickerMessage(t)
				globalContext.OkxPriceComposite.UpdatePriceList(tickerMsg, globalConfig)
			}
			if r.Int31n(10000) < 5 && len(s.Tickers) > 0 {
				logger.Info("[Gather] Receive Okx Futures Ticker %+v", s.Tickers[0])
			}
		}
	}()

	logger.Info("[Gather] Start Gather Okx Futures Ticker")
}
