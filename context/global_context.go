package context

import (
	"okxdata/client"
	"okxdata/config"
	"okxdata/container"
	"sync"
)

type GlobalContext struct {
	InstrumentComposite   container.InstrumentComposite
	OkxPriceComposite     container.PriceComposite
	BinancePriceComposite container.PriceComposite
	rwLock                sync.RWMutex
}

func (context *GlobalContext) Init(globalConfig *config.Config, globalOkxClient *client.OkxClient) {
	// 初始化交易对数据
	context.initInstrumentComposite(globalConfig)

	// 初始化volatility数据
	context.initPriceComposite(globalConfig)

	context.rwLock = sync.RWMutex{}
}

func (context *GlobalContext) initInstrumentComposite(globalConfig *config.Config) {
	instrumentComposite := container.InstrumentComposite{}
	instrumentComposite.Init(globalConfig.InstIDs)
	context.InstrumentComposite = instrumentComposite
}

func (context *GlobalContext) initPriceComposite(globalConfig *config.Config) {
	okxPriceComposite := container.PriceComposite{}
	okxPriceComposite.InitOkx(globalConfig)
	context.OkxPriceComposite = okxPriceComposite

	binancePriceComposite := container.PriceComposite{}
	binancePriceComposite.InitBinance(globalConfig)
	context.BinancePriceComposite = binancePriceComposite
}
