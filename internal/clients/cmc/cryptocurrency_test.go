package cmc_test

import (
	"context"
	_ "embed"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jnericks/obibot/internal/clients/cmc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/cryptocurrency_quotes_latest.json
var cryptocurrencyQuotesLatest []byte

//go:embed testdata/cryptocurrency_quotes_latest_bad_request.json
var cryptocurrencyQuotesLatestBadRequest []byte

func TestCryptocurrency(t *testing.T) {
	const apiKey = "fake-apiKey"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/cryptocurrency/quotes/latest", r.URL.Path)
		assert.Equal(t, apiKey, r.Header.Get("X-CMC_PRO_API_KEY"))

		_, err := w.Write(cryptocurrencyQuotesLatest)
		require.NoError(t, err)
	}))

	client, err := cmc.NewClient(http.DefaultClient, server.URL, apiKey)
	require.NoError(t, err)

	resp, err := client.GetCryptocurrencyQuotes(context.Background(), cmc.GetCryptocurrencyQuotesParams{
		Symbols: []string{"BTC", "ETH", "DOGE"},
	})
	require.NoError(t, err)

	assert.Len(t, resp.Quotes, 3)
	assert.Empty(t, resp.Error)
}

func TestCryptocurrencyBadRequest(t *testing.T) {
	const apiKey = "fake-apiKey"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.URL.Path, "/cryptocurrency/quotes/latest")
		assert.Equal(t, apiKey, r.Header.Get("X-CMC_PRO_API_KEY"))

		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write(cryptocurrencyQuotesLatestBadRequest)
		require.NoError(t, err)
	}))

	client, err := cmc.NewClient(http.DefaultClient, server.URL, apiKey)
	require.NoError(t, err)

	resp, err := client.GetCryptocurrencyQuotes(context.Background(), cmc.GetCryptocurrencyQuotesParams{
		Symbols: []string{"BTC", "YOOO"},
	})
	require.NoError(t, err)

	assert.Empty(t, resp.Quotes)
	assert.Equal(t, `Invalid value for "symbol": "YOOO"`, resp.Error)
}
