package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/shaurya947/btc-price-estimator/estimator"
)

// Notes:
// It's generally inadvisable to use floats when working with currency. However,
// since this is just an example exercise, and the exchanges also return floats,
// we are using floats for now. In a production system that was dealing with
// currencies, we would decide ahead of time on the precision we need (example
// 1/100th or 1/1000th of a cent) and create a custom Currency type. We would
// then add methods on this type to create instances and do math operations.
//
// Binance API uses the BTC/USDT pair. This example assumes that USDT and USD are
// interchangeable. This is obviously not ideal, but using other APIs like
// coinbase for example required an account for keys, and I wanted to keep
// the setup simple for this example.
//
// User input is expected to come through the "price" url param. We can handle
// the dollar sign, currency commas, as well as decimal points. Any other
// formatting will return an error.
//
// In case one of the dex price fetches fails for any reason, we just use the
// other one. If both fail, we return an error.
//
// The BTCPriceEstimator struct is thread-safe as-is because it doesn't maintain
// any dynamic state (besides the static sources). If we were to implement
// caching or some other dynamic state within it, we'd probably want to use
// mutexes or similar to ensure thread safety, and test using Go's race detector.

func main() {
	port := "8080"
	priceEstimator := estimator.NewBTCPriceEstimator(
		&estimator.BinanceAPI{}, &estimator.BitfinexAPI{},
	)
	http.HandleFunc("/", priceDifferenceBTC(priceEstimator))
	log.Printf("Starting server on port %s", port)
	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func priceDifferenceBTC(priceEstimator *estimator.BTCPriceEstimator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		priceInputStr := r.URL.Query().Get("price")
		if len(priceInputStr) == 0 {
			http.Error(w, "No price supplied", http.StatusBadRequest)
			return
		}

		priceInput, err := getFloatPrice(priceInputStr)
		if err != nil {
			http.Error(w, "Invalid price input", http.StatusBadRequest)
			return
		}

		priceChan := make(chan *float64)
		go fetchMeanPrice(priceEstimator, priceChan)
		select {
		case maybePrice := <-priceChan:
			if maybePrice == nil {
				http.Error(w, "Internal error, try again later", http.StatusInternalServerError)
				return
			}
			difference := math.Abs(*maybePrice - priceInput)
			fmt.Fprintf(w, "$%.2f", difference)
		case <-ctxTimeout.Done():
			http.Error(w, "Timeout", http.StatusInternalServerError)
		}
	}
}

func fetchMeanPrice(priceEstimator *estimator.BTCPriceEstimator, priceChan chan *float64) {
	meanPrice, err := priceEstimator.EstimateMean()
	if err != nil {
		priceChan <- nil
	}
	priceChan <- &meanPrice
}

func getFloatPrice(priceInputStr string) (float64, error) {
	if priceInputStr[0] == '$' {
		priceInputStr = priceInputStr[1:]
	}

	priceInputStr = strings.Replace(priceInputStr, ",", "", -1)
	return strconv.ParseFloat(priceInputStr, 64)
}
