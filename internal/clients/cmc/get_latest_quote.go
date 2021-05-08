package cmc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
)

type (
	GetLatestQuoteParams struct {
		Symbols []string
	}

	GetLatestQuoteResponse struct {
		Data  []Cryptocurrency
		Error string
	}
)

func (c *client) GetLatestQuote(ctx context.Context, params GetLatestQuoteParams) (*GetLatestQuoteResponse, error) {
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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cmc get latest quote response (%d %s)", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	var data struct {
		Data   map[string]Cryptocurrency `json:"data"`
		Status Status                    `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("decoding cmc get latest quote response: %w", err)
	}

	cryptos := make([]Cryptocurrency, 0, len(data.Data))
	for _, c := range data.Data {
		cryptos = append(cryptos, c)
	}

	sort.Slice(cryptos, func(i, j int) bool {
		return cryptos[i].Rank < cryptos[j].Rank
	})
	return &GetLatestQuoteResponse{
		Data:  cryptos,
		Error: data.Status.ErrorMessage,
	}, nil
}
