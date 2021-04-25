package iex

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type GetQuoteParams struct {
	Symbol string `validate:"required"`
}

func (c *client) urlGetQuote(symbol string) string {
	return fmt.Sprintf("%s/stock/%s/quote", c.baseURL, symbol)
}

func (c *client) GetQuote(ctx context.Context, params GetQuoteParams) (QuoteResponse, error) {
	if err := c.validate.Struct(params); err != nil {
		return QuoteResponse{}, fmt.Errorf("validating get quote params: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodGet, c.urlGetQuote(params.Symbol))
	if err != nil {
		return QuoteResponse{}, fmt.Errorf("creating iex quote request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return QuoteResponse{}, fmt.Errorf("executing iex quote response: %w", err)
	}

	var out QuoteResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return QuoteResponse{}, fmt.Errorf("decoding iex quote response: %w", err)
	}

	return out, nil
}
