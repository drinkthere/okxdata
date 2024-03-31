package config

import (
	"encoding/json"
	"os"

	"go.uber.org/zap/zapcore"
)

type Config struct {
	// 日志配置
	LogLevel zapcore.Level
	LogPath  string

	// Okx配置
	OkxAPIKey    string
	OkxSecretKey string
	OkxPassword  string

	InstIDs      []string
	KeepPricesMs int64
	MinAccuracy  float64 // 价格最小精度
}

func LoadConfig(filename string) *Config {
	config := new(Config)
	reader, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer reader.Close()

	// 加载配置
	decoder := json.NewDecoder(reader)
	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}

	return config
}
