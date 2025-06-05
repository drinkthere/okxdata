package message

import (
	"github.com/drinkthere/cryptodotcom/events/public"
	"math/rand"
	"okxdata/config"
	"okxdata/container"
	"okxdata/context"
	"okxdata/utils"
	"okxdata/utils/logger"
	"strconv"
	"time"
)

func convertToCryptoTickerMessage(ticker *public.Ticker) container.TickerMessage {
	instID := ticker.InstrumentName
	instType := utils.GetCryptoInstTypeFromInstID(instID)
	if instType == config.UnknownInstrument {
		logger.Error("[Ticker] unknown instrument type, %+v", ticker)
	}

	askPx, _ := strconv.ParseFloat(ticker.AskPx, 64)
	askSz, _ := strconv.ParseFloat(ticker.AskSz, 64)
	bidPx, _ := strconv.ParseFloat(ticker.BidPx, 64)
	bidSz, _ := strconv.ParseFloat(ticker.BidSz, 64)
	return container.TickerMessage{
		Exchange:     config.CryptoExchange,
		InstType:     instType,
		InstID:       instID,
		AskPx:        askPx,
		AskSz:        askSz,
		BidPx:        bidPx,
		BidSz:        bidSz,
		UpdateTimeMs: time.Time(ticker.UpdateTime).UnixMilli(),
	}
}

func StartGatherCryptoSwapTicker(
	globalConfig *config.Config,
	globalContext *context.GlobalContext,
	tickChan chan *public.Tickers) {

	r := rand.New(rand.NewSource(2))
	go func() {
		defer func() {
			if rc := recover(); rc != nil {
				logger.Error("[Gather%sTicker] Recovered from panic: %v", config.CryptoExchange, rc)
			}

			logger.Warn("[Gather%sTicker] Swap Ticker Gather Exited.", config.CryptoExchange)
		}()

		instIDs := globalContext.InstrumentComposite.CryptoSwapInstIDs
		for {
			s := <-tickChan
			for _, t := range s.Result.Data {
				instID := t.InstrumentName
				if !utils.InArray(instID, instIDs) {
					continue
				}
				tickerMsg := convertToCryptoTickerMessage(t)
				globalContext.PriceComposite.UpdatePriceList(tickerMsg, globalConfig)
			}
			if r.Int31n(10000) < 5 && len(s.Result.Data) > 0 {
				logger.Info("[Gather%sTicker] Receive Crypto Swap Ticker %+v", config.CryptoExchange, s.Result.Data[0])
			}
		}
	}()

	logger.Info("[Gather%sTicker] Start Gather Crypto Swap Ticker", globalConfig.Exchange)
}

func StartGatherCryptoSpotTicker(
	globalConfig *config.Config,
	globalContext *context.GlobalContext,
	tickChan chan *public.Tickers) {

	r := rand.New(rand.NewSource(2))
	go func() {
		defer func() {
			if rc := recover(); rc != nil {
				logger.Error("[Gather%sTicker] Recovered from panic: %v", globalConfig.Exchange, rc)
			}

			logger.Warn("[Gather%sTicker] Ticker Gather Exited.", globalConfig.Exchange)
		}()

		instIDs := globalContext.InstrumentComposite.GetInstIDs(config.CryptoExchange, config.SpotInstrument)
		wrapper := globalContext.TickerComposite.GetTickerWrapper(config.CryptoExchange, config.SpotInstrument)
		if wrapper == nil {
			logger.Error("[Gather%sTicker] Tickers Wrapper is Nil", globalConfig.Exchange)
			return
		}
		for {
			s := <-tickChan
			for _, t := range s.Result.Data {
				instID := t.InstrumentName
				if !utils.InArray(instID, instIDs) {
					continue
				}

				tickerMsg := convertToCryptoTickerMessage(t)
				wrapper.UpdateTicker(tickerMsg)

				cancelInstID := utils.FormatCompareSourceToCryptoInstID(config.CryptoExchange, config.SpotInstrument, instID, globalConfig.InstType)
				globalContext.TickerToCancelCh <- cancelInstID
			}
			if r.Int31n(10000) < 5 && len(s.Result.Data) > 0 {
				logger.Info("[Gather%sTicker] Receive Crypto Spot Ticker %+v", globalConfig.Exchange, s.Result.Data[0])
			}
		}
	}()

	logger.Info("[Gather%sTicker] Start Gather Crypto Spot Ticker", globalConfig.Exchange)
}
