package container

import (
	"okxdata/config"
	"okxdata/utils"
	"okxdata/utils/logger"
	"sort"
	"sync"
	"time"
)

type PriceItem struct {
	Value      float64
	UpdateTime time.Time
}

type PriceListMap struct {
	instIDPriceListMap map[string][]PriceItem
	rwLock             *sync.RWMutex
}

func newPriceListMap() *PriceListMap {
	return &PriceListMap{
		instIDPriceListMap: map[string][]PriceItem{},
		rwLock:             new(sync.RWMutex),
	}
}

type PriceComposite struct {
	BinanceDeliveryPriceMap *PriceListMap
	BinanceFuturesPriceMap  *PriceListMap
	BinanceSpotPriceMap     *PriceListMap
	OkxSwapPriceMap         *PriceListMap
	OkxSpotPriceMap         *PriceListMap
	bybitLinearPriceMap     *PriceListMap
	bybitSpotPriceMap       *PriceListMap
}

func NewPriceComposite(globalConfig *config.Config) *PriceComposite {
	composite := &PriceComposite{
		BinanceDeliveryPriceMap: newPriceListMap(),
		BinanceFuturesPriceMap:  newPriceListMap(),
		BinanceSpotPriceMap:     newPriceListMap(),
		OkxSwapPriceMap:         newPriceListMap(),
		OkxSpotPriceMap:         newPriceListMap(),
		bybitLinearPriceMap:     newPriceListMap(),
		bybitSpotPriceMap:       newPriceListMap(),
	}

	for _, instID := range globalConfig.BinanceDeliveryInstIDs {
		// 初始化，无需加锁
		composite.BinanceDeliveryPriceMap.instIDPriceListMap[instID] = []PriceItem{}
	}

	for _, instID := range globalConfig.BinanceFuturesInstIDs {
		composite.BinanceFuturesPriceMap.instIDPriceListMap[instID] = []PriceItem{}
		composite.BinanceSpotPriceMap.instIDPriceListMap[instID] = []PriceItem{}
	}

	for _, instID := range globalConfig.OkxSwapInstIDs {
		composite.OkxSwapPriceMap.instIDPriceListMap[instID] = []PriceItem{}

		okxSpotID := utils.ConvertOkxSwapInstIDToOkxSpotInstID(instID)
		composite.OkxSpotPriceMap.instIDPriceListMap[okxSpotID] = []PriceItem{}
	}

	for _, instID := range globalConfig.BybitLinearInstIDs {
		composite.bybitLinearPriceMap.instIDPriceListMap[instID] = []PriceItem{}
		bybitSpotInstID := utils.ConvertBybitLinearInstIDToBybitSpotInstID(instID)
		composite.bybitSpotPriceMap.instIDPriceListMap[bybitSpotInstID] = []PriceItem{}
	}
	return composite
}

func (c *PriceComposite) getPriceListMap(exchange config.Exchange, instType config.InstrumentType) *PriceListMap {
	var p *PriceListMap
	if exchange == config.BinanceExchange {
		if instType == config.FuturesInstrument {
			p = c.BinanceFuturesPriceMap
		} else if instType == config.SpotInstrument {
			p = c.BinanceSpotPriceMap
		} else if instType == config.DeliveryInstrument {
			p = c.BinanceDeliveryPriceMap
		}
	} else if exchange == config.OkxExchange {
		if instType == config.SwapInstrument {
			p = c.OkxSwapPriceMap
		} else if instType == config.SpotInstrument {
			p = c.OkxSpotPriceMap
		}
	} else if exchange == config.BybitExchange {
		if instType == config.LinearInstrument {
			p = c.bybitLinearPriceMap
		} else if instType == config.SpotInstrument {
			p = c.bybitSpotPriceMap
		}
	}
	return p
}
func (c *PriceComposite) GetPriceList(exchange config.Exchange, instType config.InstrumentType, instID string) *[]PriceItem {
	p := c.getPriceListMap(exchange, instType)
	if p == nil {
		return nil
	}

	p.rwLock.RLock()
	defer p.rwLock.RUnlock()
	wrapper, has := p.instIDPriceListMap[instID]

	if has {
		return &wrapper
	} else {
		return nil
	}
}

func (c *PriceComposite) UpdatePriceList(tickerMsg TickerMessage, globalConfig *config.Config) bool {
	p := c.getPriceListMap(tickerMsg.Exchange, tickerMsg.InstType)
	if p == nil {
		return false
	}

	instID := tickerMsg.InstID
	p.rwLock.Lock()
	defer p.rwLock.Unlock()
	priceList, has := p.instIDPriceListMap[instID]

	if !has {
		logger.Error("[Price] Update Price List Failed, instID=%s does not exist")
		return false
	} else {
		if tickerMsg.BidPx < globalConfig.MinAccuracy || tickerMsg.AskPx < globalConfig.MinAccuracy {
			return false
		}
		price := (tickerMsg.BidPx + tickerMsg.AskPx) / 2
		priceItem := PriceItem{
			Value:      price,
			UpdateTime: time.Now(),
		}
		priceList = append(priceList, priceItem)

		// 清除过期数据
		if len(priceList) > 0 {
			cutoffTime := time.Now().Add(-time.Duration(globalConfig.KeepPricesMs) * time.Millisecond)

			// 使用sort.Search找到第一个不满足要求的时间点的索引
			index := sort.Search(len(priceList), func(i int) bool {
				return !priceList[i].UpdateTime.Before(cutoffTime)
			})

			// 保留该索引之后的所有元素
			priceList = priceList[index:]
		}

		p.instIDPriceListMap[instID] = priceList
	}
	return true
}
