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

func TestCryptocurrency(t *testing.T) {
	const apiKey = "fake-apiKey"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.URL.Path, "/cryptocurrency/quotes/latest")

		assert.Equal(t, apiKey, r.Header.Get("X-CMC_PRO_API_KEY"))

		_, err := w.Write(cryptocurrencyQuotesLatest)
		require.NoError(t, err)
	}))

	client, err := cmc.NewClient(http.DefaultClient, server.URL, apiKey)
	require.NoError(t, err)

	resp, err := client.GetLatestQuote(context.Background(), cmc.GetLatestQuoteParams{
		Symbols: []string{"BTC", "ETH", "DOGE"},
	})
	require.NoError(t, err)

	require.Len(t, resp.Data, 3)
}
