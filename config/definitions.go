package config

type (
	Exchange       string
	InstrumentType string
)

const (
	BinanceExchange = Exchange("Binance")
	OkxExchange     = Exchange("Okx")
	BybitExchange   = Exchange("Bybit")
	CryptoExchange  = Exchange("Crypto")

	UnknownInstrument  = InstrumentType("UNKNOWN")
	SpotInstrument     = InstrumentType("SPOT")
	FuturesInstrument  = InstrumentType("FUTURES")
	DeliveryInstrument = InstrumentType("DELIVERY")
	SwapInstrument     = InstrumentType("SWAP")
	LinearInstrument   = InstrumentType("LINEAR")
)
