package iex

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// GetCryptoParams is the params for get crypto api.
type GetCryptoParams struct {
	// Symbol of the cryptocurrency.
	Symbol string `json:"symbol"`
}

// GetCryptoResponse is the response of the get crypto price api.
// https://iexcloud.io/docs/api/#cryptocurrency-price
type GetCryptoResponse struct {
	// Symbol of the cryptocurrency.
	Symbol string `json:"symbol"`

	// LatestPrice is the latest price of the cryptocurrency.
	LatestPrice string `json:"latestPrice"`

	// PreviousClose is the price of the cryptocurrency 24 hours ago.
	PreviousClose *string `json:"previousClose"`
}

func (r GetCryptoResponse) PriceSummary() string {
	// attempt to calculate change in latest
	latest, err := strconv.ParseFloat(r.LatestPrice, 64)
	if err != nil {
		return fmt.Sprintf("%s: %s", r.Symbol, r.LatestPrice)
	}
	out := fmt.Sprintf("%s: %.2f", r.Symbol, latest) // better formatting

	if r.PreviousClose == nil || *r.PreviousClose == "" {
		return out
	}

	previous, err := strconv.ParseFloat(*r.PreviousClose, 64)
	if err != nil {
		return out
	}

	change := latest - previous
	changePercent := change / previous * 100

	return fmt.Sprintf("%s %+.2f (%+.2f%%)", out, change, changePercent)
}

func (c *client) urlGetCrypto(symbol string) string {
	return fmt.Sprintf("%s/crypto/%s/quote", c.baseURL, symbol)
}

func (c *client) GetCrypto(ctx context.Context, params GetCryptoParams) (*GetCryptoResponse, error) {
	if err := c.validate.Struct(params); err != nil {
		return nil, fmt.Errorf("validating get crypto params: %w", err)
	}

	// force 'USDT' suffix to simplify api (can just request 'BTC' or 'ETH')
	symbol := strings.ToUpper(params.Symbol)
	if !strings.HasSuffix(symbol, "USDT") {
		symbol += "USDT"
	}

	req, err := c.newRequest(ctx, http.MethodGet, c.urlGetCrypto(symbol))
	if err != nil {
		return nil, fmt.Errorf("creating iex get quote request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing iex get crypto response: %w", err)
	}

	switch resp.StatusCode {
	case http.StatusNotFound:
		return nil, ErrInvalidSymbol{Symbol: symbol}
	case http.StatusOK:
		var out GetCryptoResponse
		if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
			return nil, fmt.Errorf("decoding iex get crypto response: %w", err)
		}
		return &out, nil
	default: // unknown error
		return nil, errors.New("getting crypto price")
	}
}
