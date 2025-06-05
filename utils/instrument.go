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

func ConvertBinanceFuturesInstIDToBinanceSpotInstID(binanceFuturesInstID string) string {
	return binanceFuturesInstID
}

func ConvertOkxSwapInstIDToOkxSpotInstID(okxSwapInstID string) string {
	return strings.Replace(okxSwapInstID, "-SWAP", "", -1)
}

func ConvertBybitLinearInstIDToBybitSpotInstID(bybitLinearInstID string) string {
	// BTC-USDT-SWAP => BTC-USDT
	return bybitLinearInstID
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

func GetCryptoInstTypeFromInstID(instID string) config.InstrumentType {
	if strings.HasSuffix(instID, "USD-PERP") {
		return config.SwapInstrument
	}
	if strings.HasSuffix(instID, "_USD") || strings.HasSuffix(instID, "_USDT") {
		return config.SpotInstrument
	}
	return config.UnknownInstrument
}
