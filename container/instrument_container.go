package container

import (
	"okxdata/config"
	"okxdata/utils"
)

type InstrumentComposite struct {
	BinanceDeliveryInstIDs []string // 币安币本位合约对应的symbol，BTCUSD_PERP, BTCUSD_240927
	BinanceFuturesInstIDs  []string // 币安U本位合约对应的symbol，BTCUSDT
	BinanceSpotInstIDs     []string // 币安现货对应的symbol，BTCUSDT。U本位永续和现货的symbol不一定相同，比如1000LUNC和LUNC
	OkxSwapInstIDs         []string // Okx永续交易对，支持多个交易对，如：BTC-USDT-SWAP, ETH-USDT-SWAP
	OkxSpotInstIDs         []string // Okx现货交易对，支持多个交易对，如：BTC-USDT, ETH-USDT
	BybitLinearInstIDs     []string // Bybit永续交易对，支持多个交易对，如：BTCUSDT, ETHUSDT
}

func NewInstrumentComposite(globalConfig *config.Config) *InstrumentComposite {
	composite := &InstrumentComposite{
		BinanceDeliveryInstIDs: globalConfig.BinanceDeliveryInstIDs,
		BinanceFuturesInstIDs:  []string{},
		BinanceSpotInstIDs:     []string{},
		OkxSwapInstIDs:         globalConfig.OkxSwapInstIDs,
		OkxSpotInstIDs:         []string{},
		BybitLinearInstIDs:     globalConfig.BybitLinearInstIDs,
	}

	// 通过delivery id 初始化 futures和spot id
	for _, instID := range globalConfig.BinanceDeliveryInstIDs {
		composite.BinanceFuturesInstIDs = append(composite.BinanceFuturesInstIDs, utils.ConvertBinanceDeliveryInstIDToFuturesInstID(instID))
		composite.BinanceSpotInstIDs = append(composite.BinanceSpotInstIDs, utils.ConvertBinanceDeliveryInstIDToSpotInstID(instID))
	}
	// 因为可能存在BTCUSD_PERPETUAL和BTCUSD_240927 映射到相同的instID的情况，这里进行去重
	composite.BinanceFuturesInstIDs = utils.RemoveDuplicateInstIDs(composite.BinanceFuturesInstIDs)
	composite.BinanceSpotInstIDs = utils.RemoveDuplicateInstIDs(composite.BinanceSpotInstIDs)

	// 通过swap id 初始化 spot id
	for _, instID := range globalConfig.OkxSwapInstIDs {
		composite.OkxSpotInstIDs = append(composite.OkxSpotInstIDs, utils.ConvertOkxSwapInstIDToOkxSpotInstID(instID))
	}
	return composite
}
