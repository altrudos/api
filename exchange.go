package charityhonor

import (
	"net/http"
	"time"

	"github.com/monstercat/golib/request"
)

var defaultExchangeRate = map[string]float64{
	"CAD": 1.4047115119,
	"HKD": 7.7536578633,
	"ISK": 143.4618569983,
	"PHP": 50.5558111714,
	"DKK": 6.8700653354,
	"HUF": 326.4562436735,
	"CZK": 24.7621238612,
	"GBP": 0.805788166,
	"RON": 4.4474095887,
	"SEK": 10.0722370479,
	"IDR": 15867.4979295114,
	"INR": 76.311309469,
	"BRL": 5.1491672035,
	"RUB": 74.2523235484,
	"HRK": 7.009754302,
	"JPY": 108.8892978743,
	"THB": 32.8195454127,
	"CHF": 0.9715652894,
	"EUR": 0.9202171713,
	"MYR": 4.3375356584,
	"BGN": 1.7997607435,
	"TRY": 6.7390264102,
	"CNY": 7.058893899,
	"NOK": 10.3195914236,
	"NZD": 1.668169688,
	"ZAR": 18.0715008742,
	"USD": 1,
	"MXN": 23.9551854238,
	"SGD": 1.4244041594,
	"AUD": 1.6052268335,
	"ILS": 3.5813932088,
	"KRW": 1216.9780068096,
	"PLN": 4.1949019969,
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
