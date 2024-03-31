package container

import (
	"okxdata/config"
	"sync"
	"time"
)

type TickerMessage struct {
	Exchange config.Exchange
	InstType config.InstrumentType
	InstID   string
	AskPx    float64
	AskSz    float64
	BidPx    float64
	BidSz    float64
}

type TickerWrapper struct {
	Exchange config.Exchange
	InstType config.InstrumentType
	InstID   string

	BidPrice        float64 // 买1价
	BidVolume       float64 // 买1量
	AskPrice        float64 // 卖1价
	AskVolume       float64 // 卖1量
	UpdateTimeMicro int64   //更新时间（微秒）
}

func (wrapper *TickerWrapper) init(exchange config.Exchange, instType config.InstrumentType, instID string) {
	wrapper.Exchange = exchange
	wrapper.InstType = instType
	wrapper.InstID = instID
	wrapper.BidPrice = 0.0
	wrapper.BidVolume = 0.0
	wrapper.AskPrice = 0.0
	wrapper.AskVolume = 0.0
	wrapper.UpdateTimeMicro = 0
}

func (wrapper *TickerWrapper) IsExpired(micro int64) bool {
	return time.Now().UnixMicro()-wrapper.UpdateTimeMicro > micro
}

func (wrapper *TickerWrapper) updateTicker(message TickerMessage) bool {
	if message.AskPx > 0.0 && message.AskSz > 0.0 {
		wrapper.AskPrice = message.AskPx
		wrapper.AskVolume = message.AskSz
	}

	if message.BidPx > 0.0 && message.BidSz > 0.0 {
		wrapper.BidPrice = message.BidPx
		wrapper.BidVolume = message.BidSz
	}

	wrapper.UpdateTimeMicro = time.Now().UnixMicro()
	return true
}

type TickerComposite struct {
	Exchange       config.Exchange
	InstType       config.InstrumentType
	tickerWrappers map[string]TickerWrapper
	rwLock         *sync.RWMutex
}

func (composite *TickerComposite) Init(exchange config.Exchange, instType config.InstrumentType) {
	composite.Exchange = exchange
	composite.InstType = instType
	composite.tickerWrappers = map[string]TickerWrapper{}
	composite.rwLock = new(sync.RWMutex)
}

func (composite *TickerComposite) GetTicker(instID string) *TickerWrapper {
	composite.rwLock.RLock()
	wrapper, has := composite.tickerWrappers[instID]
	composite.rwLock.RUnlock()

	if has {
		return &wrapper
	} else {
		return nil
	}
}

func (composite *TickerComposite) UpdateTicker(message TickerMessage) bool {
	if composite.Exchange != message.Exchange || composite.InstType != message.InstType {
		return false
	}

	composite.rwLock.RLock()
	wrapper, has := composite.tickerWrappers[message.InstID]
	composite.rwLock.RUnlock()

	if !has {
		wrapper = TickerWrapper{}
		wrapper.init(composite.Exchange, composite.InstType, message.InstID)
		wrapper.updateTicker(message)
		composite.rwLock.Lock()
		composite.tickerWrappers[message.InstID] = wrapper
		composite.rwLock.Unlock()
	} else {
		wrapper.updateTicker(message)

		composite.rwLock.Lock()
		composite.tickerWrappers[message.InstID] = wrapper
		composite.rwLock.Unlock()
	}
	return true
}
