package estimator

import (
	"fmt"
	"math"
	"testing"
)

func TestMeanAllSourcesWorking(t *testing.T) {
	source1 := testSource{priceToReturn: 22000}
	source2 := testSource{priceToReturn: 22500}
	priceEstimator := NewBTCPriceEstimator(&source1, &source2)

	meanPrice, err := priceEstimator.EstimateMean()
	if err != nil {
		t.Errorf("Unexpected error %+v", err)
	}

	if !almostEqual(meanPrice, 22250) {
		t.Errorf("Unexpected mean %+v", meanPrice)
	}
}

func TestMeanOnlyOneSourceWorking(t *testing.T) {
	source1 := testSource{shouldError: true}
	source2 := testSource{priceToReturn: 22500}
	priceEstimator := NewBTCPriceEstimator(&source1, &source2)

	meanPrice, err := priceEstimator.EstimateMean()
	if err != nil {
		t.Errorf("Unexpected error %+v", err)
	}

	if !almostEqual(meanPrice, 22500) {
		t.Errorf("Unexpected mean %+v", meanPrice)
	}
}

func TestMeanNoSourcesWorking(t *testing.T) {
	source1 := testSource{shouldError: true}
	source2 := testSource{shouldError: true}
	priceEstimator := NewBTCPriceEstimator(&source1, &source2)

	_, err := priceEstimator.EstimateMean()
	if err == nil {
		t.Errorf("Unexpected error %+v", err)
	}
}

const float64EqualityThreshold = 1e-9

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}

type testSource struct {
	shouldError   bool
	priceToReturn float64
}

func (ts *testSource) FetchBTCPrice() (float64, error) {
	if ts.shouldError {
		return 0, fmt.Errorf("Error fetching BTC price")
	}

	return ts.priceToReturn, nil
}
