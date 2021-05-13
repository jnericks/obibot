package commands_test

import (
	"testing"

	"github.com/jnericks/obibot/internal/clients/cmc"
	"github.com/jnericks/obibot/internal/commands"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	cryptoBTC = cmc.CryptocurrencyQuote{
		Name:   "Bitcoin",
		Symbol: "BTC",
		Rank:   1,
		Quote: cmc.USDQuote{
			USD: cmc.Quote{
				Price:            57200.27405442796,
				Volume24H:        69278660657.27121,
				PercentChange1H:  -0.03518281,
				PercentChange24H: 6.74517701,
				PercentChange7D:  4.49438325,
				PercentChange30D: -2.7693705,
				PercentChange60D: 16.55866453,
				PercentChange90D: 54.61716182,
				MarketCap:        1069689054628.2766,
			},
		},
	}
	cryptoETH = cmc.CryptocurrencyQuote{
		Name:   "Ethereum",
		Symbol: "ETH",
		Rank:   2,
		Quote: cmc.USDQuote{
			USD: cmc.Quote{
				Price:            3519.2771620316908,
				Volume24H:        48507691889.498695,
				PercentChange1H:  0.58517401,
				PercentChange24H: 8.36818046,
				PercentChange7D:  28.55227351,
				PercentChange30D: 67.58398754,
				PercentChange60D: 111.37776398,
				PercentChange90D: 119.69062823,
				MarketCap:        407371242898.9393,
			},
		},
	}
	cryptoDOGE = cmc.CryptocurrencyQuote{
		Name:   "Dogecoin",
		Symbol: "DOGE",
		Rank:   4,
		Quote: cmc.USDQuote{
			USD: cmc.Quote{
				Price:            0.64800566229809,
				Volume24H:        42246109296.67634,
				PercentChange1H:  1.46170998,
				PercentChange24H: 19.33389632,
				PercentChange7D:  102.52328806,
				PercentChange30D: 989.67831719,
				PercentChange60D: 1168.78441148,
				PercentChange90D: 1073.39162554,
				MarketCap:        83900690210.49292,
			},
		},
	}
)

func TestFormatGetCryptocurrencyQuotesResponse(t *testing.T) {
	a, err := commands.FormatGetCryptocurrencyQuotesResponse(&cmc.GetCryptocurrencyQuotesResponse{
		Quotes: []cmc.CryptocurrencyQuote{
			cryptoBTC,
			cryptoETH,
			cryptoDOGE,
		},
		Error: "",
	})
	require.NoError(t, err)

	e := &commands.Output{
		Response: `BTC: 57200.27 +7385.28 (+6.75%)
ETH: 3519.28 +375.66 (+8.37%)
DOGE: 0.6480 +0.0319 (+19.33%)`,
		Markdown: false,
	}

	assert.Equal(t, e, a)
}
