package iex

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
)

type (
	GetStockQuotesParams struct {
		Symbols []string `validate:"required"`
	}

	GetStockQuotesResponse struct {
		Quotes []StockQuote
		Error  string
	}
)

// StockQuote is the iex quote api response.
// https://iexcloud.io/docs/api/#quote
type StockQuote struct {
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

	// For US stocks, indicates if the market is in normal market hours. Will be false during extended hours trading.
	IsUSMarketOpen bool `json:"isUSMarketOpen"`
}

func (r StockQuote) PriceSummary() string {
	return fmt.Sprintf("%s: %.2f %+.2f (%+.2f%%)", r.Symbol, r.LatestPrice, r.Change, r.ChangePercent*100)
}

func (c *client) GetStockQuotes(ctx context.Context, params GetStockQuotesParams) (*GetStockQuotesResponse, error) {
	if err := c.validate.Struct(params); err != nil {
		return nil, fmt.Errorf("validating get stock market batch params: %w", err)
	}

	req, err := c.newStockMarketBatchRequest(ctx, params.Symbols)
	if err != nil {
		return nil, fmt.Errorf("creating iex get stock market batch request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing iex get stock market batch response: %w", err)
	}

	switch resp.StatusCode {
	case http.StatusNotFound:
		return &GetStockQuotesResponse{
			Quotes: nil,
			Error:  fmt.Sprintf("no quotes found for any of %v", params.Symbols),
		}, nil

	case http.StatusOK:
		var data map[string]struct {
			Quote StockQuote `json:"quote"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, fmt.Errorf("decoding iex get stock market batch response: %w", err)
		}

		stocks := make([]StockQuote, 0, len(data))
		for _, q := range data {
			stocks = append(stocks, q.Quote)
		}
		if len(stocks) == 0 {
			return &GetStockQuotesResponse{
				Quotes: nil,
				Error:  fmt.Sprintf("no stock quotes found for any of %v", params.Symbols),
			}, nil
		}

		// sort by percent change desc
		sort.Slice(stocks, func(i, j int) bool {
			return stocks[i].ChangePercent > stocks[j].ChangePercent
		})

		return &GetStockQuotesResponse{
			Quotes: stocks,
			Error:  "",
		}, nil

	default: // unknown error
		return nil, fmt.Errorf("iex get stock market batch response (%d %s)", resp.StatusCode, http.StatusText(resp.StatusCode))
	}
}
