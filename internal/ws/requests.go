package ws

import (
	"encoding/json"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"log"
)

func tickerRequest(platform string, upbitMarkets []string) []byte {
	switch platform {
	case "upbit":
		request := []map[string]interface{}{
			{"ticket": uuid.New().String()},
			{"type": "ticker", "isOnlyRealtime": true, "codes": upbitMarkets},
			{"format": "SIMPLE"},
		}
		jsonRequest, err := json.Marshal(request)
		if err != nil {
			log.Fatal("Failed to create json request to websocket", zap.Error(err))
		}
		return jsonRequest
	case "bithumb":
		request := `{
			"type" : "ticker",
			"symbols" : ["BTC_KRW", "ETH_KRW"],
			"tickTypes" : ["30M", "1H", "12H", "24H", "MID"]
		}`
		return []byte(request)
	default:
		log.Fatal("No such platform registered")
		return []byte("")
	}
}

func tradeRequest(platform string, upbitMarkets []string) []byte {
	switch platform {
	case "upbit":
		request := []map[string]interface{}{
			{"ticket": uuid.New().String()},
			{"type": "trade", "isOnlyRealtime": true, "codes": upbitMarkets},
			{"format": "SIMPLE"},
		}
		jsonRequest, err := json.Marshal(request)
		if err != nil {
			log.Fatal("Failed to create json request to websocket", zap.Error(err))
		}
		return jsonRequest
	case "bithumb":
		request := `{
		  "type" : "transaction", 
		  "symbols" : ["BTC_KRW" , "ETH_KRW"]
		}`
		return []byte(request)
	default:
		log.Fatal("No such platform registered")
		return []byte("")
	}
}
