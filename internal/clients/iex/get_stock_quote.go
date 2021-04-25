package iex

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type GetStockQuoteParams struct {
	Symbol string `validate:"required"`
}

// GetStockQuoteResponse is the iex quote api response.
// https://iexcloud.io/docs/api/#quote
type GetStockQuoteResponse struct {
	// Symbol refers to the stock ticker.
	Symbol string `json:"symbol"`

	// CompanyName refers to the company name.
	CompanyName string `json:"companyName"`

	// PreviousClose refers to the previous trading day closing price.
	PreviousClose float64 `json:"previousClose"`

	// LatestPrice refers to the latest relevant price of the security which is derived from multiple sources.
	LatestPrice float64 `json:"latestPrice"`

	// Change refers to the change in price between LatestPrice and PreviousClose.
	Change float64 `json:"change"`

	// ChangePercent refers to the percent change in price between LatestPrice and PreviousClose.
	//
	// For example, a 5% change would be represented as 0.05.
	ChangePercent float64 `json:"changePercent"`
}

func (r GetStockQuoteResponse) PriceSummary() string {
	return fmt.Sprintf("%s: %.2f %+.2f (%+.2f%%)", r.Symbol, r.LatestPrice, r.Change, r.ChangePercent*100)
}

func (c *client) urlGetQuote(symbol string) string {
	return fmt.Sprintf("%s/stock/%s/quote", c.baseURL, symbol)
}

func (c *client) GetStockQuote(ctx context.Context, params GetStockQuoteParams) (*GetStockQuoteResponse, error) {
	if err := c.validate.Struct(params); err != nil {
		return nil, fmt.Errorf("validating get quote params: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodGet, c.urlGetQuote(params.Symbol))
	if err != nil {
		return nil, fmt.Errorf("creating iex get quote request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing iex get quote response: %w", err)
	}

	var out GetStockQuoteResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decoding iex get quote response: %w", err)
	}

	return &out, nil
}
