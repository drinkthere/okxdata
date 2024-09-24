package main

import (
	"encoding/json"
	"log"
	"net/http"
	"okxdata/config"
	"okxdata/utils/logger"
	"strconv"
	"strings"
	"time"
)

func startHTTPServer(cfg *config.Config) {
	go func() {
		http.HandleFunc("/api/v1/getInstrumentPriceList", getInstPriceListHandler())
		logger.Info("Server listening on %s", cfg.Port)
		port := strconv.Itoa(cfg.Port)
		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()
}

func getInstPriceListHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		exchange := r.URL.Query().Get("exchange")
		// 兼容老服务的请求
		if exchange == "" {
			exchange = "Okx"
		}
		fmtExchange := config.Exchange(exchange)

		instType := r.URL.Query().Get("instType")
		instType = strings.ToUpper(instType)
		// 兼容老服务的请求
		if instType == "" {
			if fmtExchange == config.BinanceExchange {
				instType = "FUTURES"
			} else if fmtExchange == config.OkxExchange {
				instType = "SWAP"
			} else if fmtExchange == config.BybitExchange {
				instType = "LINEAR"
			}
		}
		fmtInstType := config.InstrumentType(instType)

		instID := r.URL.Query().Get("instID")
		instID = strings.ToUpper(instID)

		timePeriodMsStr := r.URL.Query().Get("timePeriodMs")
		timePeriodMs, err := strconv.ParseInt(timePeriodMsStr, 10, 64)
		if err != nil {
			logger.Info("[Main] TimePeriodMsStr:%s is invalid", timePeriodMsStr)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		logger.Info("exchange=%s, instType=%s, instID=%s, timePeriodMs=%d", fmtExchange, fmtInstType, instID, timePeriodMs)

		priceList := globalContext.PriceComposite.GetPriceList(fmtExchange, fmtInstType, instID)

		if priceList == nil {
			w.WriteHeader(http.StatusNotFound)
			logger.Info("[Main] Can Not Find %s's Price List", instID)
			return
		}

		var prices []float64
		cutoffTime := time.Now().Add(-time.Duration(timePeriodMs) * time.Millisecond)
		for _, price := range *priceList {
			if price.UpdateTime.After(cutoffTime) {
				prices = append(prices, price.Value) // 将价格添加到 relevantPrices 中
			}
		}
		response := struct {
			Code int       `json:"code"`
			Data []float64 `json:"data"`
		}{
			Code: http.StatusOK,
		}
		response.Data = prices

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
