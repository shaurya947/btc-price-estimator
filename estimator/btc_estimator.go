package estimator

import "fmt"

type BTCPriceEstimator struct {
	// exported so consumer can inspect and/or query each source
	Sources []Source
}

func NewBTCPriceEstimator(sources ...Source) *BTCPriceEstimator {
	return &BTCPriceEstimator{sources}
}

// Returns the mean price of BTC across all sources. As long as at least 1
// source is working, this function won't error. If all sources return an error,
// then this function will return an error and the float should be discarded.
func (b *BTCPriceEstimator) EstimateMean() (float64, error) {
	var prices []float64
	priceChan := make(chan *float64, len(b.Sources))
	for _, source := range b.Sources {
		go fetchPriceFromSource(source, priceChan)
	}

	for range b.Sources {
		maybePrice := <-priceChan
		if maybePrice != nil {
			prices = append(prices, *maybePrice)
		}
	}

	if len(prices) == 0 {
		return 0, fmt.Errorf("All sources errored out")
	}

	var meanPrice float64
	for _, price := range prices {
		meanPrice += price
	}
	meanPrice /= float64(len(prices))
	return meanPrice, nil
}

func fetchPriceFromSource(source Source, priceChan chan *float64) {
	price, err := source.FetchBTCPrice()
	if err != nil {
		priceChan <- nil
		return
	}
	priceChan <- &price
}
