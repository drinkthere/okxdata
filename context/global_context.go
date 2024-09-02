package context

import (
	"okxdata/config"
	"okxdata/container"
	"sync"
)

type GlobalContext struct {
	InstrumentComposite *container.InstrumentComposite
	PriceComposite      *container.PriceComposite
	rwLock              sync.RWMutex
}

func (context *GlobalContext) Init(globalConfig *config.Config) {
	// 初始化交易对数据
	context.initInstrumentComposite(globalConfig)

	// 初始化volatility数据
	context.initPriceComposite(globalConfig)

	context.rwLock = sync.RWMutex{}
}

func (context *GlobalContext) initInstrumentComposite(globalConfig *config.Config) {
	instrumentComposite := container.NewInstrumentComposite(globalConfig.DeliveryInstIDs)
	context.InstrumentComposite = instrumentComposite
}

func (context *GlobalContext) initPriceComposite(globalConfig *config.Config) {
	priceComposite := container.NewPriceComposite(globalConfig.DeliveryInstIDs)
	context.PriceComposite = priceComposite
}
