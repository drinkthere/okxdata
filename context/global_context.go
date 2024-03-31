package context

import (
	"okxdata/client"
	"okxdata/config"
	"okxdata/container"
	"sync"
)

type GlobalContext struct {
	InstrumentComposite       container.InstrumentComposite
	OkxFuturesTickerComposite container.TickerComposite
	PriceComposite            container.PriceComposite
	rwLock                    sync.RWMutex
}

func (context *GlobalContext) Init(globalConfig *config.Config, globalOkxClient *client.OkxClient) {
	// 初始化交易对数据
	context.initInstrumentComposite(globalConfig)

	// 初始化ticker数据
	context.initTickerComposite()

	// 初始化volatility数据
	context.initPirceComposite(globalConfig)

	context.rwLock = sync.RWMutex{}
}

func (context *GlobalContext) initInstrumentComposite(globalConfig *config.Config) {
	instrumentComposite := container.InstrumentComposite{}
	instrumentComposite.Init(globalConfig.InstIDs)
	context.InstrumentComposite = instrumentComposite
}

func (context *GlobalContext) initTickerComposite() {
	okxFuturesComposite := container.TickerComposite{}
	okxFuturesComposite.Init(config.OkxExchange, config.FuturesInstrument)
	context.OkxFuturesTickerComposite = okxFuturesComposite
}

func (context *GlobalContext) initPirceComposite(globalConfig *config.Config) {
	priceComposite := container.PriceComposite{}
	priceComposite.Init(globalConfig)
	context.PriceComposite = priceComposite
}
