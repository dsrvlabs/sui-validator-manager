package pricefeeder

type Client interface {
	QueryPrice() (float64, error)
}

type dummyPriceFeeder struct{}

func (f *dummyPriceFeeder) QueryPrice() (float64, error) {
	// Token price will be fixed during Wave 3 as 50.0 USD
	return 50.0, nil
}

// NewClient create new PriceFeeder client.
func NewClient() Client {
	return &dummyPriceFeeder{}
}
