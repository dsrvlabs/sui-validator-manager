package pricefeeder

import (
	"encoding/json"
	"io"
	"net/http"
)

type CoinGeckoPrice struct {
	Sui struct {
		USD         float64 `json:"usd"`
		LastUpdated int64   `json:"last_updated_at"`
	} `json:"sui"`
}

type Client interface {
	QueryPrice() (float64, error)
}

type dummyPriceFeeder struct{}

func (f *dummyPriceFeeder) QueryPrice() (float64, error) {
	// Token price will be fixed during Wave 3 as 50.0 USD
	base := "https://api.coingecko.com"

	path := "/api/v3/simple/price?"

	path += "ids=sui"
	path += "&vs_currencies=usd"
	path += "&include_market_cap=false"
	path += "&include_24hr_vol=false"
	path += "&include_24hr_change=false"
	path += "&include_last_updated_at=true"
	path += "&precision=2"

	url := base + path

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	cli := http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		return 0, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	respData := &CoinGeckoPrice{}
	err = json.Unmarshal(data, respData)
	if err != nil {
		return 0, err
	}

	return respData.Sui.USD, nil
}

// NewClient create new PriceFeeder client.
func NewClient() Client {
	return &dummyPriceFeeder{}
}
