package utils

import (
	"okxdata/config"
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
