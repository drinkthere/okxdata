package container

import "okxdata/utils"

type InstrumentComposite struct {
	InstIDs        []string // 永续合约支持多个交易对，如：BTC-USDT-SWAP, ETH-USDT-SWAP
	BinanceInstIDs []string // 币安现货和永续交易对一样（这里复用），支持多个交易对，如：BTCUSDT, ETHUSDT
}

func (composite *InstrumentComposite) Init(instIDs []string) {
	for _, instID := range instIDs {
		composite.InstIDs = append(composite.InstIDs, instID)
		composite.BinanceInstIDs = append(composite.BinanceInstIDs, utils.ConvertToBinanceInstID(instID))
	}
}
