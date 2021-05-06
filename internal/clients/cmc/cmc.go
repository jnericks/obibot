package cmc

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-playground/validator"
)

type Client interface {
	GetLatestQuote(context.Context, GetLatestQuoteParams) (*GetLatestQuoteResponse, error)
}

type client struct {
	baseURL  string
	apiKey   string
	validate *validator.Validate
	http     *http.Client
}

func (c *client) newRequest(ctx context.Context, method, url string, symbols []string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-CMC_PRO_API_KEY", c.apiKey)

	q := req.URL.Query()
	q.Set("symbol", strings.Join(symbols, ","))
	req.URL.RawQuery = q.Encode()

	return req.WithContext(ctx), nil
}

func NewClient(httpClient *http.Client, baseURL, apiKey string) (Client, error) {
	validate := validator.New()
	if err := validate.Struct(struct {
		APIKey string `validate:"required"`
	}{
		APIKey: apiKey,
	}); err != nil {
		return nil, err
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("could not parse base url: %w", err)
	}

	return &client{
		baseURL:  u.String(),
		apiKey:   apiKey,
		validate: validate,
		http:     httpClient,
	}, nil
}
