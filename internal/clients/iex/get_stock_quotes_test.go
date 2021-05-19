package iex_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jnericks/obibot/internal/clients/iex"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetStockQuotes(t *testing.T) {
	type testConfig struct {
		apiToken   string
		statusCode int
		fileSuffix string
		params     iex.GetStockQuotesParams
	}

	getStockQuotes := func(t *testing.T, cfg testConfig) (*iex.GetStockQuotesResponse, error) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, cfg.apiToken, r.URL.Query().Get("token"))
			assert.Equal(t, "/stock/market/batch", r.URL.Path)

			if cfg.statusCode == http.StatusNotFound {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			filename := fmt.Sprintf("testdata/stock_market_batch_%s.json", cfg.fileSuffix)
			b, err := testData.ReadFile(filename)
			require.NoError(t, err)

			_, err = w.Write(b)
			require.NoError(t, err)
		}))
		defer server.Close()

		client, err := iex.NewClient(http.DefaultClient, server.URL, cfg.apiToken)
		require.NoError(t, err)

		return client.GetStockQuotes(context.Background(), cfg.params)
	}

	t.Run("symbol=tsla", func(t *testing.T) {
		cfg := testConfig{
			apiToken:   "apiToken-tsla",
			statusCode: http.StatusOK,
			fileSuffix: "tsla",
			params: iex.GetStockQuotesParams{
				Symbols: []string{
					"tsla",
				},
			},
		}
		resp, err := getStockQuotes(t, cfg)
		require.NoError(t, err)

		assert.Len(t, resp.Quotes, 1)
		assert.Empty(t, resp.Error)

		assert.Equal(t, iex.StockQuote{
			Symbol:         "TSLA",
			CompanyName:    "Tesla Inc",
			PreviousClose:  617.2,
			LatestPrice:    601.08,
			Change:         -16.12,
			ChangePercent:  -0.02612,
			IsUSMarketOpen: true,
		}, resp.Quotes[0])
	})

	t.Run("symbol=tsla-amzn", func(t *testing.T) {
		cfg := testConfig{
			apiToken:   "apiToken-tsla-amzn",
			statusCode: http.StatusOK,
			fileSuffix: "tsla-amzn",
			params: iex.GetStockQuotesParams{
				Symbols: []string{
					"tsla",
					"amzn",
				},
			},
		}
		resp, err := getStockQuotes(t, cfg)
		require.NoError(t, err)

		assert.Len(t, resp.Quotes, 2)
		assert.Empty(t, resp.Error)

		assert.Equal(t, iex.StockQuote{
			Symbol:         "AMZN",
			CompanyName:    "Amazon.com Inc.",
			PreviousClose:  3223.91,
			LatestPrice:    3144.005,
			Change:         -79.905,
			ChangePercent:  -0.02479,
			IsUSMarketOpen: true,
		}, resp.Quotes[0])

		assert.Equal(t, iex.StockQuote{
			Symbol:         "TSLA",
			CompanyName:    "Tesla Inc",
			PreviousClose:  617.2,
			LatestPrice:    601.14,
			Change:         -16.06,
			ChangePercent:  -0.02602,
			IsUSMarketOpen: true,
		}, resp.Quotes[1])
	})

	t.Run("symbol=tsla-amzn-sq", func(t *testing.T) {
		cfg := testConfig{
			apiToken:   "apiToken-tsla-amzn-sq",
			statusCode: http.StatusOK,
			fileSuffix: "tsla-amzn-sq",
			params: iex.GetStockQuotesParams{
				Symbols: []string{
					"tsla",
					"amzn",
					"SQ",
				},
			},
		}
		resp, err := getStockQuotes(t, cfg)
		require.NoError(t, err)

		assert.Len(t, resp.Quotes, 3)
		assert.Empty(t, resp.Error)

		assert.Equal(t, iex.StockQuote{
			Symbol:         "AMZN",
			CompanyName:    "Amazon.com Inc.",
			PreviousClose:  3223.91,
			LatestPrice:    3140.68,
			Change:         -83.23,
			ChangePercent:  -0.02582,
			IsUSMarketOpen: true,
		}, resp.Quotes[0])

		assert.Equal(t, iex.StockQuote{
			Symbol:         "TSLA",
			CompanyName:    "Tesla Inc",
			PreviousClose:  617.2,
			LatestPrice:    600.36,
			Change:         -16.84,
			ChangePercent:  -0.02728,
			IsUSMarketOpen: true,
		}, resp.Quotes[1])

		assert.Equal(t, iex.StockQuote{
			Symbol:         "SQ",
			CompanyName:    "Square Inc - Class A",
			PreviousClose:  220.65,
			LatestPrice:    208.79,
			Change:         -11.86,
			ChangePercent:  -0.05375,
			IsUSMarketOpen: true,
		}, resp.Quotes[2])
	})

	t.Run("symbol=tsla-xyz", func(t *testing.T) {
		cfg := testConfig{
			apiToken:   "apiToken-tsla-xyz",
			statusCode: http.StatusOK,
			fileSuffix: "tsla-xyz",
			params: iex.GetStockQuotesParams{
				Symbols: []string{
					"tsla",
					"xyz", // no data
				},
			},
		}
		resp, err := getStockQuotes(t, cfg)
		require.NoError(t, err)

		assert.Len(t, resp.Quotes, 1)
		assert.Empty(t, resp.Error)

		assert.Equal(t, iex.StockQuote{
			Symbol:         "TSLA",
			CompanyName:    "Tesla Inc",
			PreviousClose:  617.2,
			LatestPrice:    599.8,
			Change:         -17.4,
			ChangePercent:  -0.02819,
			IsUSMarketOpen: true,
		}, resp.Quotes[0])
	})

	t.Run("symbol=abc,xyz", func(t *testing.T) {
		cfg := testConfig{
			apiToken:   "apiToken-abc-xyz",
			statusCode: http.StatusNotFound,
			fileSuffix: "",
			params: iex.GetStockQuotesParams{
				Symbols: []string{
					"abc", // no data
					"xyz", // no data
				},
			},
		}
		resp, err := getStockQuotes(t, cfg)
		require.NoError(t, err)

		assert.Empty(t, resp.Quotes)
		assert.Equal(t, "no quotes found for any of [abc xyz]", resp.Error)
	})

	t.Run("symbol=game", func(t *testing.T) {
		cfg := testConfig{
			apiToken:   "apiToken-game",
			statusCode: http.StatusOK,
			fileSuffix: "game",
			params: iex.GetStockQuotesParams{
				Symbols: []string{
					"game",
				},
			},
		}
		resp, err := getStockQuotes(t, cfg)
		require.NoError(t, err)

		assert.Empty(t, resp.Quotes)
		assert.Equal(t, "no quotes found for any of [game]", resp.Error)
	})
}
