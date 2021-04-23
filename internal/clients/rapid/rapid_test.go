package rapid_test

import (
	"testing"

	"github.com/jnericks/obibot/internal/clients/rapid"
	"github.com/stretchr/testify/assert"
)

func TestMeta(t *testing.T) {
	a := rapid.ChartMeta{
		Symbol:             "AAA",
		RegularMarketPrice: 76.34,
		ChartPreviousClose: 100.38,
	}
	assert.Equal(t, "AAA: 76.34 -24.04 (-23.95%)", a.PriceSummary())

	b := rapid.ChartMeta{
		Symbol:             "BBB",
		RegularMarketPrice: 123.67,
		ChartPreviousClose: 100.23,
	}
	assert.Equal(t, "BBB: 123.67 +23.44 (+23.39%)", b.PriceSummary())
}
