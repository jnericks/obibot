package cmc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
)

type (
	GetCryptocurrencyQuotesParams struct {
		Symbols []string `validate:"required"`
	}

	GetCryptocurrencyQuotesResponse struct {
		Quotes []CryptocurrencyQuote
		Error  string
	}
)

func (c *client) GetCryptocurrencyQuotes(ctx context.Context, params GetCryptocurrencyQuotesParams) (*GetCryptocurrencyQuotesResponse, error) {
	if err := c.validate.Struct(params); err != nil {
		return nil, fmt.Errorf("validating get crypto quote params: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodGet, c.baseURL+"/cryptocurrency/quotes/latest", params.Symbols)
	if err != nil {
		return nil, fmt.Errorf("creating cmc get latest quote request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing cmc get latest quote response: %w", err)
	}

	switch resp.StatusCode {
	case http.StatusBadRequest:
		var data struct {
			Status Status `json:"status"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, fmt.Errorf("decoding cmc get latest quote response: %w", err)
		}

		errMsg := data.Status.ErrorMessage
		if errMsg == "" {
			errMsg = fmt.Sprintf("error retrieving crypto data for %v", params.Symbols)
		}

		return &GetCryptocurrencyQuotesResponse{
			Quotes: nil,
			Error:  errMsg,
		}, nil

	case http.StatusOK:
		var data struct {
			Data map[string]CryptocurrencyQuote `json:"data"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, fmt.Errorf("decoding cmc get latest quote response: %w", err)
		}

		cryptos := make([]CryptocurrencyQuote, 0, len(data.Data))
		for _, c := range data.Data {
			cryptos = append(cryptos, c)
		}
		if len(cryptos) == 0 {
			return &GetCryptocurrencyQuotesResponse{
				Quotes: nil,
				Error:  fmt.Sprintf("no crypto quotes found for any of %v", params.Symbols),
			}, nil
		}

		// sort by percent change desc
		sort.Slice(cryptos, func(i, j int) bool {
			return cryptos[i].Quote.USD.PercentChange24H > cryptos[j].Quote.USD.PercentChange24H
		})

		return &GetCryptocurrencyQuotesResponse{
			Quotes: cryptos,
			Error:  "",
		}, nil

	default: // unknown error
		return nil, fmt.Errorf("cmc get latest quote response (%d %s)", resp.StatusCode, http.StatusText(resp.StatusCode))
	}
}
