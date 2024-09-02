package utils

import (
	"okxdata/config"
	"strings"
)

func ConvertBinanceDeliveryInstIDToFuturesInstID(binanceDeliveryInstID string) string {
	// BTCUSD_PERP => BTCUSDT
	// BTCUSD_240927 => BTCUSDT

	parts := strings.Split(binanceDeliveryInstID, "_")
	baseQuote := parts[0] // 获取前半部分，例如 "BTCUSD"

	return strings.Replace(baseQuote, "USD", "USDT", -1)
}

func ConvertBinanceDeliveryInstIDToSpotInstID(binanceDeliveryInstID string) string {
	// BTCUSD_PERP => BTCUSDT
	// BTCUSD_240927 => BTCUSDT

	parts := strings.Split(binanceDeliveryInstID, "_")
	baseQuote := parts[0] // 获取前半部分，例如 "BTCUSD"

	return strings.Replace(baseQuote, "USD", "USDT", -1)
}

func ConvertBinanceDeliveryInstIDToOkxSwapInstID(binanceDeliveryInstID string) string {
	// BTCUSD_PERPETUAL => BTC-USDT-SWAP
	// BTCUSD_240927 => BTC-USDT-SWAP

	parts := strings.Split(binanceDeliveryInstID, "_")
	baseQuote := parts[0] // 获取前半部分，例如 "BTCUSD"

	return strings.Replace(baseQuote, "USD", "-USDT-SWAP", -1)
}

func ConvertBinanceDeliveryInstIDToOkxSpotInstID(binanceDeliveryInstID string) string {
	// BTCUSD_PERP => BTC-USDT
	// BTCUSD_240927 => BTC-USDT
	parts := strings.Split(binanceDeliveryInstID, "_")
	baseQuote := parts[0] // 获取前半部分，例如 "BTCUSD"

	return strings.Replace(baseQuote, "USD", "-USDT", -1)
}

func RemoveDuplicateInstIDs(ids []string) []string {
	uniqueIDs := make(map[string]bool)
	var result []string

	for _, id := range ids {
		if _, ok := uniqueIDs[id]; !ok {
			uniqueIDs[id] = true
			result = append(result, id)
		}
	}

	return result
}

func ConvertToStdInstType(exchange config.Exchange, instType string) config.InstrumentType {
	if exchange == config.OkxExchange {
		switch instType {
		case "SWAP":
			return config.SwapInstrument
		case "SPOT":
			return config.SpotInstrument
		}
	}
	return config.UnknownInstrument
}

func ConvertToBinanceInstID(okxFuturesInstID string) string {
	// BTC-USDT-SWAP => BTCUSDT
	return strings.Replace(okxFuturesInstID, "-USDT-SWAP", "USDT", -1)
}
