package container

import (
	"okxdata/config"
	"okxdata/utils"
)

type InstrumentComposite struct {
	BinanceDeliveryInstIDs []string // 币安币本位合约对应的symbol，BTCUSD_PERP, BTCUSD_240927
	BinanceFuturesInstIDs  []string // 币安U本位合约对应的symbol，BTCUSDT
	BinanceSpotInstIDs     []string // 币安U本位合约对应的symbol，BTCUSDT
	OkxSwapInstIDs         []string // Okx永续交易对，支持多个交易对，如：BTC-USDT-SWAP, ETH-USDT-SWAP
	OkxSpotInstIDs         []string // Okx现货交易对，支持多个交易对，如：BTC-USDT, ETH-USDT
	BybitLinearInstIDs     []string // Bybit永续交易对，支持多个交易对，如：BTCUSDT, ETHUSDT
	BybitSpotInstIDs       []string // Bybit现货交易对，支持多个交易对，如：BTCUSDT, ETHUSDT
	CryptoSwapInstIDs      []string // crypto.com永续交易对，支持多个交易对，如：BTCUSD-PERP， ETHUSD-PERP
}

func NewInstrumentComposite(globalConfig *config.Config) *InstrumentComposite {
	composite := &InstrumentComposite{
		BinanceDeliveryInstIDs: globalConfig.BinanceDeliveryInstIDs,
		BinanceFuturesInstIDs:  globalConfig.BinanceFuturesInstIDs,
		BinanceSpotInstIDs:     []string{},
		OkxSwapInstIDs:         globalConfig.OkxSwapInstIDs,
		OkxSpotInstIDs:         []string{},
		BybitLinearInstIDs:     globalConfig.BybitLinearInstIDs,
		BybitSpotInstIDs:       []string{},
		CryptoSwapInstIDs:      globalConfig.CryptoSwapInstIDs,
	}

	if len(globalConfig.BinanceFuturesInstIDs) > 0 {
		for _, instID := range globalConfig.BinanceFuturesInstIDs {
			composite.BinanceSpotInstIDs = append(composite.BinanceSpotInstIDs, utils.ConvertBinanceFuturesInstIDToBinanceSpotInstID(instID))
		}
	}

	// 通过swap id 初始化 spot id
	if len(globalConfig.OkxSpotInstIDs) == 0 {
		if len(globalConfig.OkxSwapInstIDs) > 0 {
			for _, instID := range globalConfig.OkxSwapInstIDs {
				composite.OkxSpotInstIDs = append(composite.OkxSpotInstIDs, utils.ConvertOkxSwapInstIDToOkxSpotInstID(instID))
			}
		}
	} else {
		composite.OkxSpotInstIDs = globalConfig.OkxSpotInstIDs
	}

	if len(globalConfig.BybitLinearInstIDs) > 0 {
		for _, instID := range globalConfig.BybitLinearInstIDs {
			composite.BybitSpotInstIDs = append(composite.BybitSpotInstIDs, utils.ConvertBybitLinearInstIDToBybitSpotInstID(instID))
		}
	}

	return composite
}
