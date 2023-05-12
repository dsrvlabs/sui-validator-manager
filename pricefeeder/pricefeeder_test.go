package pricefeeder

import (
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestQueryPrice(t *testing.T) {
	httpmock.Activate()
	defer httpmock.Deactivate()

	url := "https://api.coingecko.com/api/v3/simple/price?ids=sui&vs_currencies=usd&include_market_cap=false&include_24hr_vol=false&include_24hr_change=false&include_last_updated_at=true&precision=2"

	fixture := "./fixtures/price.json"

	f, err := os.Open(fixture)
	defer f.Close()

	assert.Nil(t, err)

	data, err := io.ReadAll(f)
	assert.Nil(t, err)

	httpmock.RegisterResponder(
		http.MethodGet,
		url,
		httpmock.NewStringResponder(http.StatusOK, string(data)))

	cli := NewClient()

	price, err := cli.QueryPrice()

	assert.Nil(t, err)
	assert.Equal(t, 1.04, price)
}
