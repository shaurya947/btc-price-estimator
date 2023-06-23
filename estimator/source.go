package estimator

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// exported in case consumer wishes to provide custom sources
type Source interface {
	FetchBTCPrice() (float64, error)
}

type BinanceAPI struct{}

// https://binance-docs.github.io/apidocs/spot/en/#symbol-price-ticker
func (b *BinanceAPI) FetchBTCPrice() (float64, error) {
	resp, err := http.Get("https://api.binance.com/api/v3/ticker/price?symbol=BTCUSDT")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	priceData := struct {
		Symbol, Price string
	}{}
	err = decoder.Decode(&priceData)
	if err != nil {
		return 0, err
	}

	floatPrice, err := strconv.ParseFloat(priceData.Price, 64)
	if err != nil {
		return 0, err
	}

	return floatPrice, nil
}

type BitfinexAPI struct{}

// https://docs.bitfinex.com/reference/rest-public-ticker
func (b *BitfinexAPI) FetchBTCPrice() (float64, error) {
	resp, err := http.Get("https://api-pub.bitfinex.com/v2/ticker/tBTCUSD")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var tickerData []float64
	err = decoder.Decode(&tickerData)
	if err != nil {
		return 0, err
	}

	return tickerData[6], nil
}
