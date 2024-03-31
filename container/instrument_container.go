package container

type InstrumentComposite struct {
	InstIDs []string // 永续合约支持多个交易对，如：BTC-USDT-SWAP, ETH-USDT-SWAP
}

func (composite *InstrumentComposite) Init(instIDs []string) {
	for _, instID := range instIDs {
		composite.InstIDs = append(composite.InstIDs, instID)
	}
}
