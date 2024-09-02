package main

import (
	"fmt"
	binanceDelivery "github.com/dictxwang/go-binance/delivery"
	binanceFutures "github.com/dictxwang/go-binance/futures"
	"github.com/drinkthere/okx/events/public"
	"okxdata/config"
	"okxdata/context"
	"okxdata/message"
	"okxdata/utils/logger"
	"os"
	"time"
)

var globalConfig config.Config
var globalContext context.GlobalContext

func startWebSocket() {
	// 监听okx行情信息
	okxFuturesTickerChan := make(chan *public.Tickers)
	message.StartOkxMarketWs(&globalConfig, &globalContext, okxFuturesTickerChan)

	// 开启Okx行情数据收集整理
	message.StartGatherOkxFuturesTicker(okxFuturesTickerChan, &globalConfig, &globalContext)

	// 监听binance行情信息并收集整理
	binanceFuturesTickerChan := make(chan *binanceFutures.WsBookTickerEvent)
	binanceDeliveryTickerChan := make(chan *binanceDelivery.WsBookTickerEvent)
	message.StartBinanceMarketWs(&globalConfig, &globalContext, binanceFuturesTickerChan, binanceDeliveryTickerChan)
	message.StartGatherBinanceDeliveryBookTicker(binanceDeliveryTickerChan, &globalConfig, &globalContext)
	message.StartGatherBinanceFuturesBookTicker(binanceFuturesTickerChan, &globalConfig, &globalContext)
}

func main() {
	// 参数判断
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s config_file\n", os.Args[0])
		os.Exit(1)
	}

	// 加载配置文件
	globalConfig = *config.LoadConfig(os.Args[1])

	// 设置日志级别, 并初始化日志
	logger.InitLogger(globalConfig.LogPath, globalConfig.LogLevel)

	// 解析config，加载杠杆和合约交易对，初始化context，账户初始化设置，拉取仓位、余额等
	globalContext.Init(&globalConfig)

	// 开始监听ws消息
	startWebSocket()

	// 等等ws数据
	time.Sleep(10 * time.Second)

	// 启动HTTP 服务
	startHTTPServer(&globalConfig)

	// 阻塞主进程
	for {
		time.Sleep(24 * time.Hour)
	}
}
