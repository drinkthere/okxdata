package main

import (
	"fmt"
	binanceSpot "github.com/dictxwang/go-binance"
	binanceDelivery "github.com/dictxwang/go-binance/delivery"
	binanceFutures "github.com/dictxwang/go-binance/futures"
	"github.com/drinkthere/okx/events/public"
	"github.com/hirokisan/bybit/v2"
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
	if len(globalContext.InstrumentComposite.OkxSwapInstIDs) > 0 {
		// 监听okx行情信息
		okxFuturesTickerChan := make(chan *public.Tickers)
		okxSpotTickerChan := make(chan *public.Tickers)
		message.StartOkxMarketWs(&globalConfig, &globalContext, okxFuturesTickerChan, okxSpotTickerChan)
		// 开启Okx行情数据收集整理
		message.StartGatherOkxFuturesTicker(okxFuturesTickerChan, &globalConfig, &globalContext)
		message.StartGatherOkxSpotTicker(okxSpotTickerChan, &globalConfig, &globalContext)
	}

	if len(globalContext.InstrumentComposite.BinanceFuturesInstIDs) > 0 {
		// 监听binance行情信息并收集整理
		binanceFuturesTickerChan := make(chan *binanceFutures.WsBookTickerEvent)
		binanceDeliveryTickerChan := make(chan *binanceDelivery.WsBookTickerEvent)
		binanceSpotTickerChan := make(chan *binanceSpot.WsBookTickerEvent)
		message.StartBinanceMarketWs(&globalConfig, &globalContext, binanceFuturesTickerChan,
			binanceDeliveryTickerChan, binanceSpotTickerChan)
		message.StartGatherBinanceDeliveryBookTicker(binanceDeliveryTickerChan, &globalConfig, &globalContext)
		message.StartGatherBinanceFuturesBookTicker(binanceFuturesTickerChan, &globalConfig, &globalContext)
		message.StartGatherBinanceSpotBookTicker(binanceSpotTickerChan, &globalConfig, &globalContext)
	}

	if len(globalContext.InstrumentComposite.BybitLinearInstIDs) > 0 {
		bybitLinearTickerChan := make(chan *bybit.V5WebsocketPublicTickerResponse)
		bybitSpotTickerChan := make(chan *bybit.V5WebsocketPublicTickerResponse)
		message.StartBybitMarketWs(&globalConfig, &globalContext, bybitLinearTickerChan, bybitSpotTickerChan)
		message.StartGatherBybitLinearTicker(bybitLinearTickerChan, &globalConfig, &globalContext)
		message.StartGatherBybitSpotTicker(bybitSpotTickerChan, &globalConfig, &globalContext)
	}

	if len(globalContext.InstrumentComposite.CryptoSwapInstIDs) > 0 {
		swapTickerChan := make(chan *public.Tickers)

		message.StartCryptoMarketWs(&globalConfig, &globalContext, config.SwapInstrument, swapTickerChan)
		message.StartGatherCryptoSwapTicker(&globalConfig, &globalContext, swapTickerChan)
	}
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
