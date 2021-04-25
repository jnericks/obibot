package iex

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetCryptoPriceParams is the params for get crypto api.
type GetCryptoPriceParams struct {
	// Symbol of the cryptocurrency.
	Symbol string `json:"symbol"`
}

// GetCryptoPriceResponse is the response of the get crypto price api.
// https://iexcloud.io/docs/api/#cryptocurrency-price
type GetCryptoPriceResponse struct {
	// Symbol of the cryptocurrency.
	Symbol string `json:"symbol"`

	// Price of the cryptocurrency.
	Price string `json:"price"`
}

func (r GetCryptoPriceResponse) PriceSummary() string {
	return fmt.Sprintf("%s: %s", r.Symbol, r.Price)
}

func (c *client) urlGetCrypto(symbol string) string {
	return fmt.Sprintf("%s/crypto/%s/price", c.baseURL, symbol)
}

func (c *client) GetCryptoPrice(ctx context.Context, params GetCryptoPriceParams) (*GetCryptoPriceResponse, error) {
	if err := c.validate.Struct(params); err != nil {
		return nil, fmt.Errorf("validating get crypto params: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodGet, c.urlGetCrypto(params.Symbol))
	if err != nil {
		return nil, fmt.Errorf("creating iex get quote request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing iex get crypto response: %w", err)
	}

	var out GetCryptoPriceResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decoding iex get crypto response: %w", err)
	}

	return &out, nil
}
