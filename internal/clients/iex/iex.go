package iex

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-playground/validator"
)

type ErrInvalidSymbol struct {
	Symbol string
}

func (e ErrInvalidSymbol) Error() string {
	return fmt.Sprintf("symbol '%s' not found", e.Symbol)
}

type Client interface {
	GetStockQuotes(context.Context, GetStockQuotesParams) (*GetStockQuotesResponse, error)
}

type client struct {
	baseURL  string
	apiToken string
	validate *validator.Validate
	http     *http.Client
}

func (c *client) newStockMarketBatchRequest(ctx context.Context, symbols []string) (*http.Request, error) {
	if len(symbols) == 0 {
		return nil, errors.New("no symbols")
	}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/stock/market/batch", c.baseURL), nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Set("token", c.apiToken)
	q.Set("types", "quote")
	q.Set("symbols", strings.Replace(strings.Join(symbols, ","), " ", "", -1))
	req.URL.RawQuery = q.Encode()
	return req.WithContext(ctx), nil
}

func NewClient(httpClient *http.Client, baseURL, apiToken string) (Client, error) {
	validate := validator.New()
	if err := validate.Struct(struct {
		APIToken string `validate:"required"`
	}{
		APIToken: apiToken,
	}); err != nil {
		return nil, err
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("could not parse base url: %w", err)
	}

	return &client{
		baseURL:  u.String(),
		apiToken: apiToken,
		validate: validate,
		http:     httpClient,
	}, nil
}
