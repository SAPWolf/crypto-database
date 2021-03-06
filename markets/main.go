package main

import (
	"log"
	"os"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/ivopetiz/go-binance/binance"
	"github.com/jyap808/go-bittrex"
	"github.com/jyap808/go-cryptopia"
	"github.com/jyap808/go-poloniex"
	"github.com/shopspring/decimal"
)

func decimalToFloat(val decimal.Decimal) float64 {
	var valFloat float64
	valFloat, _ = val.Float64()

	return valFloat
}

func main() {
	//create your file with desired read/write permissions
	f, err := os.OpenFile("/log/altdb_coin.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	//defer to close when you're done with it, not because you think it's idiomatic!
	defer f.Close()
	//set output of logs to f
	log.SetOutput(f)

	// BITTREX
	bittrex := bittrex.New(API_KEY, API_SECRET)

	// POLONIEX
	poloniex := poloniex.New(API_KEY, API_SECRET)

	// CRYPTOPIA
	cryptopia := cryptopia.New(API_KEY, API_SECRET)

	// BINANCE
	binance := binance.New(API_KEY, API_SECRET)

	// Create a new HTTPClient
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     server_url,
		//Addr:	  "http://localhost:8086",
		Username: username,
		Password: password,
	})
	if err != nil {
		log.Println(err)
	}

	for true {
		start_time := time.Now()

		// Create a new point batch
		bp, err := client.NewBatchPoints(client.BatchPointsConfig{
			Database:  markets_DB,
			Precision: "s",
		})
		if err != nil {
			log.Println(err)
		}

		// BINANCE PART
		if _binance {
			marketSummaries, err := binance.GetAll24Hr()
			if err != nil {
				log.Println(err)
			}
			// Result from: GET /api/v1/ticker/24hr
			for _, coin := range marketSummaries {
				// Create a point and add to batch
				tags := map[string]string{"MarketName": coin.Symbol}
				fields := map[string]interface{}{
					"Ask":        coin.AskPrice,
					"BaseVolume": coin.Volume,
					"Bid":        coin.BidPrice,
					"High":       coin.HighPrice,
					"Last":       coin.LastPrice,
					"Low":        coin.LowPrice,
					"Volume":     coin.Volume,
				}
				//log.Println(err, marketSummaries)
				pt, err := client.NewPoint("binance", tags, fields, time.Now())
				if err != nil {
					log.Println(err)
				}
				bp.AddPoint(pt)
			}
		}

		// BITTREX PART
		if _bittrex {
			marketSummaries, err := bittrex.GetMarketSummaries()
			if err != nil {
				log.Println(err)
			}
			for _, coin := range marketSummaries {
				// Create a point and add to batch
				tags := map[string]string{"MarketName": coin.MarketName}
				fields := map[string]interface{}{
					"Ask":            coin.Ask,
					"BaseVolume":     coin.BaseVolume,
					"Bid":            coin.Bid,
					"High":           coin.High,
					"Last":           coin.Last,
					"Low":            coin.Low,
					"OpenBuyOrders":  float64(coin.OpenBuyOrders),
					"OpenSellOrders": float64(coin.OpenSellOrders),
					"PrevDay":        coin.PrevDay,
					"Volume":         coin.Volume,
				}
				//log.Println(err, marketSummaries)
				pt, err := client.NewPoint("bittrex", tags, fields, time.Now())
				if err != nil {
					log.Println(err)
				}
				bp.AddPoint(pt)
			}
		}

		// POLONIEX PART
		if _poloniex {
			tickers, err := poloniex.GetTickers()
			if err != nil {
				log.Println(err)
			}
			for key, ticker := range tickers {
				// Create a point and add to batch
				tags := map[string]string{"MarketName": key}
				fields := map[string]interface{}{
					"Ask":        decimalToFloat(ticker.LowestAsk),
					"BaseVolume": decimalToFloat(ticker.BaseVolume),
					"Bid":        decimalToFloat(ticker.HighestBid),
					"High":       decimalToFloat(ticker.High24Hr),
					"Last":       decimalToFloat(ticker.Last),
					"Low":        decimalToFloat(ticker.Low24Hr),
					//"OpenBuyOrders": float64(ticker.OpenBuyOrders),
					//"OpenSellOrders": float64(ticker.OpenSellOrders),
					//"PrevDay": ticker.PrevDay,
					"Volume": decimalToFloat(ticker.QuoteVolume),
				}
				//log.Println(err, marketSummaries)
				pt, err := client.NewPoint("poloniex", tags, fields, time.Now())
				if err != nil {
					log.Println(err)
				}
				bp.AddPoint(pt)
			}
		}

		// CRYPTOPIA PART
		if _cryptopia {
			markets, err := cryptopia.GetMarkets()
			if err != nil {
				log.Println(err)
			}
			for _, market := range markets {
				// Create a point and add to batch

				tags := map[string]string{"MarketName": market.Label}
				fields := map[string]interface{}{
					"Ask":            market.AskPrice,
					"Bid":            market.BidPrice,
					"Last":           market.LastPrice,
					"Low":            market.Low,
					"High":           market.High,
					"BaseVolume":     market.BaseVolume,
					"BaseBuyVolume":  market.BaseBuyVolume,
					"BaseSellVolume": market.BaseSellVolume,
					"Change":         market.Change,
					"Open":           market.Open,
					"Close":          market.Close,
					//"OpenBuyOrders": float64(market.OpenBuyOrders),
					//"OpenSellOrders": float64(market.OpenSellOrders),
					//"PrevDay": market.PrevDay,
					"Volume":     market.Volume,
					"BuyVolume":  market.BuyVolume,
					"SellVolume": market.SellVolume,
				}
				//log.Println(err, marketSummaries)
				pt, err := client.NewPoint("cryptopia", tags, fields, time.Now())
				if err != nil {
					log.Println(err)
				}
				bp.AddPoint(pt)

			}
		}

		// Write the batch
		if err := c.Write(bp); err != nil {
			log.Println(err)
		}

		elapsed := time.Since(start_time)
		if elapsed < interval {
			time.Sleep(interval - elapsed)

		}
	}
}
