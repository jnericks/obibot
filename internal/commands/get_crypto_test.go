package commands_test

import (
	"context"
	"sync"
	"testing"

	"github.com/jnericks/obibot/internal/clients/cmc"
	"github.com/jnericks/obibot/internal/commands"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubCMCClient struct {
	getLatestQuoteMtx sync.Mutex
	getLatestQuote    func(context.Context, cmc.GetLatestQuoteParams) (*cmc.GetLatestQuoteResponse, error)
}

func (s *stubCMCClient) GetLatestQuote(ctx context.Context, params cmc.GetLatestQuoteParams) (*cmc.GetLatestQuoteResponse, error) {
	s.getLatestQuoteMtx.Lock()
	defer s.getLatestQuoteMtx.Unlock()
	if s.getLatestQuote == nil {
		return nil, nil
	}
	return s.getLatestQuote(ctx, params)
}

var (
	cryptoBTC = cmc.Cryptocurrency{
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
	cryptoETH = cmc.Cryptocurrency{
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
	cryptoDOGE = cmc.Cryptocurrency{
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

func TestGetCrypto(t *testing.T) {
	cmdFunc := commands.GetCrypto(&stubCMCClient{
		getLatestQuote: func(ctx context.Context, params cmc.GetLatestQuoteParams) (*cmc.GetLatestQuoteResponse, error) {
			assert.ElementsMatch(t, []string{"BTC", "DOGE", "ETH"}, params.Symbols)
			return &cmc.GetLatestQuoteResponse{
				Data: []cmc.Cryptocurrency{
					cryptoBTC,
					cryptoETH,
					cryptoDOGE,
				},
			}, nil
		},
	})

	output, err := cmdFunc(context.Background(), commands.Input{
		Args: []string{"btc,eth", "doge"},
	})

	require.NoError(t, err)
	e := "```" + `
+-----------------+----------+--------+---------+----------+
|                 |    PRICE |     1H |     24H |       7D |
+-----------------+----------+--------+---------+----------+
| Bitcoin (BTC)   | 57200.27 | -0.04% |  +6.75% |   +4.49% |
| Ethereum (ETH)  |  3519.28 | +0.59% |  +8.37% |  +28.55% |
| Dogecoin (DOGE) |   0.6480 | +1.46% | +19.33% | +102.52% |
+-----------------+----------+--------+---------+----------+
` + "```"

	assert.Equal(t, e, output.Response)
	assert.True(t, output.Markdown)
}
