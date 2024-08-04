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

type PriceComposite struct {
	instIDPriceListMap map[string][]PriceItem
	rwLock             *sync.RWMutex
}

func (c *PriceComposite) InitOkx(globalConfig *config.Config) {
	c.instIDPriceListMap = map[string][]PriceItem{}
	for _, instID := range globalConfig.InstIDs {
		// 初始化，无需加锁
		c.instIDPriceListMap[instID] = []PriceItem{}

	}
	c.rwLock = new(sync.RWMutex)
}

func (c *PriceComposite) InitBinance(globalConfig *config.Config) {
	c.instIDPriceListMap = map[string][]PriceItem{}
	for _, instID := range globalConfig.InstIDs {
		// 初始化，无需加锁
		binanceInstID := utils.ConvertToBinanceInstID(instID)
		c.instIDPriceListMap[binanceInstID] = []PriceItem{}

	}
	c.rwLock = new(sync.RWMutex)
}

func (c *PriceComposite) GetPriceList(instID string) *[]PriceItem {
	c.rwLock.RLock()
	wrapper, has := c.instIDPriceListMap[instID]
	c.rwLock.RUnlock()

	if has {
		return &wrapper
	} else {
		return nil
	}
}

func (c *PriceComposite) UpdatePriceList(tickerMsg TickerMessage, globalConfig *config.Config) bool {
	instID := tickerMsg.InstID
	c.rwLock.RLock()
	priceList, has := c.instIDPriceListMap[instID]
	c.rwLock.RUnlock()

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

		c.rwLock.Lock()
		c.instIDPriceListMap[instID] = priceList
		c.rwLock.Unlock()
	}
	return true
}
