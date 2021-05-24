package commands_test

import (
	"testing"

	"github.com/jnericks/obibot/internal/clients/iex"
	"github.com/jnericks/obibot/internal/commands"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	stockAMZN = iex.StockQuote{
		Symbol:         "AMZN",
		CompanyName:    "Amazon.com Inc.",
		PreviousClose:  3223.91,
		LatestPrice:    3140.68,
		Change:         -83.23,
		ChangePercent:  -0.02582,
		IsUSMarketOpen: true,
	}
	stockTSLA = iex.StockQuote{
		Symbol:         "TSLA",
		CompanyName:    "Tesla Inc",
		PreviousClose:  617.2,
		LatestPrice:    600.36,
		Change:         -16.84,
		ChangePercent:  -0.02728,
		IsUSMarketOpen: true,
	}
	stockSQ = iex.StockQuote{
		Symbol:         "SQ",
		CompanyName:    "Square Inc - Class A",
		PreviousClose:  220.65,
		LatestPrice:    208.79,
		Change:         -11.86,
		ChangePercent:  -0.05375,
		IsUSMarketOpen: true,
	}
)

func TestFormatGetStockQuotesResponse(t *testing.T) {
	a, err := commands.FormatGetStockQuotesResponse(&iex.GetStockQuotesResponse{
		Quotes: []iex.StockQuote{
			stockTSLA,
			stockAMZN,
			stockSQ,
		},
		Error: "",
	})
	require.NoError(t, err)

	e := &commands.Output{
		Response: `TSLA: $600.36 -$16.84 (-2.73%)
AMZN: $3,140.68 -$83.23 (-2.58%)
SQ: $208.79 -$11.86 (-5.38%)`,
		Markdown: false,
	}

	assert.Equal(t, e, a)
}
