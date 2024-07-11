package price

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var logger = logrus.New().WithField("module", "price")

type StraxPrice struct {
	Stratis struct {
		Cad float64 `json:"cad"`
		Cny float64 `json:"cny"`
		Eur float64 `json:"eur"`
		Jpy float64 `json:"jpy"`
		Rub float64 `json:"rub"`
		Usd float64 `json:"usd"`
		Gbp float64 `json:"gbp"`
		Aud float64 `json:"aud"`
	} `json:"stratis"`
}

var availableCurrencies = []string{"STRAX", "USD", "EUR", "GBP", "CNY", "RUB", "CAD", "AUD", "JPY"}
var straxPrice = new(StraxPrice)
var straxPriceMux = &sync.RWMutex{}

func Init(chainId uint64) {
	go updateStraxPrice(chainId)
}

func updateStraxPrice(chainId uint64) {
	errorRetrievingStraxPriceCount := 0
	for {
		fetchPrice(chainId, &errorRetrievingStraxPriceCount)
		time.Sleep(time.Minute)
	}
}

func fetchPrice(chainId uint64, errorRetrievingEthPriceCount *int) {
	if chainId != 105105 {
		straxPrice = &StraxPrice{
			Stratis: struct {
				Cad float64 "json:\"cad\""
				Cny float64 "json:\"cny\""
				Eur float64 "json:\"eur\""
				Jpy float64 "json:\"jpy\""
				Rub float64 "json:\"rub\""
				Usd float64 "json:\"usd\""
				Gbp float64 "json:\"gbp\""
				Aud float64 "json:\"aud\""
			}{
				Cad: 0,
				Cny: 0,
				Eur: 0,
				Jpy: 0,
				Rub: 0,
				Usd: 0,
				Gbp: 0,
				Aud: 0,
			},
		}
		return
	}

	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Get("https://api.coingecko.com/api/v3/simple/price?ids=stratis&vs_currencies=usd%2Ceur%2Crub%2Ccny%2Ccad%2Cjpy%2Cgbp%2Caud")
	if err != nil {
		*errorRetrievingEthPriceCount++
		if *errorRetrievingEthPriceCount <= 3 { // warn 3 times, before throwing errors starting with the fourth time
			logger.Warnf("error (%d) retrieving STRAX price: %v", *errorRetrievingEthPriceCount, err)
		} else {
			logger.Errorf("error (%d) retrieving STRAX price: %v", *errorRetrievingEthPriceCount, err)
		}
		return
	} else {
		*errorRetrievingEthPriceCount = 0
	}

	straxPriceMux.Lock()
	defer straxPriceMux.Unlock()
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&straxPrice)
	if err != nil {
		logger.Errorf("error decoding STRAX price json response to struct: %v", err)
		return
	}
}

func GetStratisPrice(currency string) float64 {
	straxPriceMux.RLock()
	defer straxPriceMux.RUnlock()

	switch currency {
	case "EUR":
		return straxPrice.Stratis.Eur
	case "USD":
		return straxPrice.Stratis.Usd
	case "RUB":
		return straxPrice.Stratis.Rub
	case "CNY":
		return straxPrice.Stratis.Cny
	case "CAD":
		return straxPrice.Stratis.Cad
	case "AUD":
		return straxPrice.Stratis.Aud
	case "JPY":
		return straxPrice.Stratis.Jpy
	case "GBP":
		return straxPrice.Stratis.Gbp
	default:
		return 1
	}
}

func GetAvailableCurrencies() []string {
	return availableCurrencies
}

func IsAvailableCurrency(currency string) bool {
	for _, c := range availableCurrencies {
		if c == currency {
			return true
		}
	}
	return false
}

func GetCurrencyLabel(currency string) string {
	switch currency {
	case "STRAX":
		return "Stratis"
	case "USD":
		return "United States Dollar"
	case "EUR":
		return "Euro"
	case "GBP":
		return "Pound Sterling"
	case "CNY":
		return "Chinese Yuan"
	case "RUB":
		return "Russian Ruble"
	case "CAD":
		return "Canadian Dollar"
	case "AUD":
		return "Australian Dollar"
	case "JPY":
		return "Japanese Yen"
	default:
		return ""
	}
}

func GetCurrencySymbol(currency string) string {
	switch currency {
	case "EUR":
		return "€"
	case "USD":
		return "$"
	case "RUB":
		return "₽"
	case "CNY":
		return "¥"
	case "CAD":
		return "C$"
	case "AUD":
		return "A$"
	case "JPY":
		return "¥"
	case "GBP":
		return "£"
	case "STRAX":
		return "STRAX"
	default:
		return ""
	}
}
