package utils

import (
	"okxdata/config"
	"strings"
)

func ConvertToStdInstType(exchange config.Exchange, instType string) config.InstrumentType {
	if exchange == config.OkxExchange {
		switch instType {
		case "SWAP":
			return config.FuturesInstrument
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
