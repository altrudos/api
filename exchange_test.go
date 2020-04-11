package charityhonor

import "testing"

func TestExchangeToUSD(t *testing.T) {
	usd, err := ExchangeToUSD(100, "HKD")
	if err != nil {
		t.Error(err)
	}
	if usd > 100 {
		t.Error("Should probably be less")
	}

	usd, err = ExchangeToUSD(10, "USD")
	if usd != 10 {
		t.Error("Changing from USD to USD should keep same number")
	}

	if exchangeRatesLastFetched.IsZero() {
		t.Error("Exchange last fetched should be updated")
	}
}
