package iex

import "fmt"

// QuoteResponse is the iex quote api response.
// https://iexcloud.io/docs/api/#quote
type QuoteResponse struct {
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

func (r QuoteResponse) PriceSummary() string {
	return fmt.Sprintf("%s: %.2f %+.2f (%+.2f%%)", r.Symbol, r.LatestPrice, r.Change, r.ChangePercent*100)
}
