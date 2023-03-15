package pricefeeder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryPrice(t *testing.T) {
	cli := NewClient()

	price, err := cli.QueryPrice()

	assert.Nil(t, err)
	assert.Equal(t, 50.0, price)
}
