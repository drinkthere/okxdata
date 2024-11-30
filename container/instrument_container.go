package container

import (
	"okxdata/config"
	"okxdata/utils"
)

type InstrumentComposite struct {
	BinanceDeliveryInstIDs []string // 币安币本位合约对应的symbol，BTCUSD_PERP, BTCUSD_240927
	BinanceFuturesInstIDs  []string // 币安U本位合约对应的symbol，BTCUSDT
	OkxSwapInstIDs         []string // Okx永续交易对，支持多个交易对，如：BTC-USDT-SWAP, ETH-USDT-SWAP
	OkxSpotInstIDs         []string // Okx现货交易对，支持多个交易对，如：BTC-USDT, ETH-USDT
	BybitLinearInstIDs     []string // Bybit永续交易对，支持多个交易对，如：BTCUSDT, ETHUSDT
}

func NewInstrumentComposite(globalConfig *config.Config) *InstrumentComposite {
	composite := &InstrumentComposite{
		BinanceDeliveryInstIDs: globalConfig.BinanceDeliveryInstIDs,
		BinanceFuturesInstIDs:  globalConfig.BinanceFuturesInstIDs,
		OkxSwapInstIDs:         globalConfig.OkxSwapInstIDs,
		OkxSpotInstIDs:         []string{},
		BybitLinearInstIDs:     globalConfig.BybitLinearInstIDs,
	}

	// 通过swap id 初始化 spot id
	for _, instID := range globalConfig.OkxSwapInstIDs {
		composite.OkxSpotInstIDs = append(composite.OkxSpotInstIDs, utils.ConvertOkxSwapInstIDToOkxSpotInstID(instID))
	}
	return composite
}
