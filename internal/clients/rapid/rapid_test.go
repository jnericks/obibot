package rapid_test

import (
	"testing"

	rapid2 "github.com/jnericks/obibot/internal/clients/rapid"
	"github.com/stretchr/testify/assert"
)

func TestMeta(t *testing.T) {
	m := rapid2.ChartMeta{
		Symbol:             "AMZN",
		RegularMarketPrice: 3328.38,
		ChartPreviousClose: 3372.01,
	}
	assert.Equal(t, "AMZN 3328.38 -43.63 (-1.29%)", m.PriceSummary())
}
