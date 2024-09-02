package container

import (
	"okxdata/utils"
)

type InstrumentComposite struct {
	BinanceDeliveryInstIDs []string // 币安币本位合约对应的symbol，BTCUSD_PERP, BTCUSD_240927
	BinanceFuturesInstIDs  []string // 币安U本位合约对应的symbol，BTCUSDT
	BinanceSpotInstIDs     []string // 币安现货对应的symbol，BTCUSDT。U本位永续和现货的symbol不一定相同，比如1000LUNC和LUNC
	OkxSwapInstIDs         []string // Okx永续交易对，支持多个交易对，如：BTC-USDT-SWAP, ETH-USDT-SWAP
	OkxSpotInstIDs         []string // Okx现货交易对，支持多个交易对，如：BTC-USDT, ETH-USDT

}

func NewInstrumentComposite(deliveryInstIDs []string) *InstrumentComposite {
	composite := &InstrumentComposite{
		BinanceDeliveryInstIDs: deliveryInstIDs,
		BinanceFuturesInstIDs:  []string{},
		BinanceSpotInstIDs:     []string{},
		OkxSwapInstIDs:         []string{},
		OkxSpotInstIDs:         []string{},
	}

	for _, instID := range deliveryInstIDs {
		composite.BinanceFuturesInstIDs = append(composite.BinanceFuturesInstIDs, utils.ConvertBinanceDeliveryInstIDToFuturesInstID(instID))
		composite.BinanceSpotInstIDs = append(composite.BinanceSpotInstIDs, utils.ConvertBinanceDeliveryInstIDToSpotInstID(instID))
		composite.OkxSwapInstIDs = append(composite.OkxSwapInstIDs, utils.ConvertBinanceDeliveryInstIDToOkxSwapInstID(instID))
		composite.OkxSpotInstIDs = append(composite.OkxSpotInstIDs, utils.ConvertBinanceDeliveryInstIDToOkxSpotInstID(instID))
	}

	// 因为可能存在BTCUSD_PERPETUAL和BTCUSD_240927 映射到相同的instID的情况，这里进行去重
	composite.BinanceFuturesInstIDs = utils.RemoveDuplicateInstIDs(composite.BinanceFuturesInstIDs)
	composite.BinanceSpotInstIDs = utils.RemoveDuplicateInstIDs(composite.BinanceSpotInstIDs)
	composite.OkxSwapInstIDs = utils.RemoveDuplicateInstIDs(composite.OkxSwapInstIDs)
	composite.OkxSpotInstIDs = utils.RemoveDuplicateInstIDs(composite.OkxSpotInstIDs)

	return composite
}
