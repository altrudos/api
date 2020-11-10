package altrudos

import (
	"net/http"
	"time"

	"github.com/monstercat/golib/request"
)

// These 5 are the ones JustGiving supports
var defaultExchangeRate = map[string]float64{
	"CAD": 1.4047115119,
	"GBP": 0.805788166,
	"EUR": 0.9202171713,
	"USD": 1,
	"AUD": 1.6052268335,
}
var exchangeRates = map[string]float64{}
var exchangeRatesLastFetched = time.Time{}

// Update exchange rate every 24 hours
const ExchangeRateRefreshRate = time.Hour * -12

func getExchangeRates() map[string]float64 {
	if exchangeRatesLastFetched.Before(time.Now().Add(ExchangeRateRefreshRate)) {
		var err error
		exchangeRates, err = FetchExchangeRates()
		if err != nil {
			exchangeRatesLastFetched = time.Now().Add(ExchangeRateRefreshRate / 4)
			exchangeRates = make(map[string]float64)
			for k, v := range defaultExchangeRate {
				exchangeRates[k] = v
			}
		} else {
			exchangeRatesLastFetched = time.Now()
		}
	}

	return exchangeRates
}

func FetchExchangeRates() (map[string]float64, error) {
	type Resp struct {
		Rates map[string]float64 `json:"rates"`
	}
	resp := Resp{}
	params := request.Params{
		Method: http.MethodGet,
		Url:    "https://api.exchangeratesapi.io/latest?base=USD",
	}
	if err := request.Request(&params, nil, &resp); err != nil {
		return nil, err
	}
	return resp.Rates, nil
}

func ExchangeToUSD(amount int, currency string) (int, error) {
	if currency == "USD" {
		return amount, nil
	}
	rates := getExchangeRates()
	if rate, ok := rates[currency]; ok {
		return int(float64(amount) / rate), nil
	}
	return 0, ErrInvalidCurrency
}
