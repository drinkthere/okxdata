package main

import (
	"encoding/json"
	"log"
	"net/http"
	"okxdata/utils/logger"
	"strconv"
	"strings"
	"time"
)

func startHTTPServer() {
	go func() {
		http.HandleFunc("/api/v1/getInstrumentPriceList", getInstPriceListHandler())
		logger.Info("Server listening on :8888")
		err := http.ListenAndServe(":8888", nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()
}

func getInstPriceListHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		instID := r.URL.Query().Get("instID")
		instID = strings.ToUpper(instID)

		timePeriodMsStr := r.URL.Query().Get("timePeriodMs")
		timePeriodMs, err := strconv.ParseInt(timePeriodMsStr, 10, 64)
		if err != nil {
			logger.Info("[Main] TimePeriodMsStr:%s is invalid", timePeriodMsStr)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		priceList := globalContext.PriceComposite.GetPriceList(instID)
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
