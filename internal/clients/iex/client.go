package iex

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-playground/validator"
)

type Client interface {
	GetQuote(context.Context, GetQuoteParams) (QuoteResponse, error)
}

type client struct {
	baseURL  string
	apiToken string
	validate *validator.Validate
	http     *http.Client
}

func (c *client) newRequest(ctx context.Context, method, url string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Set("token", c.apiToken)
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
