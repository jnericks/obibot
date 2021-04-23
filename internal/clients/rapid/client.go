package rapid

import (
	"context"
	"net/http"

	"github.com/go-playground/validator"
)

const (
	urlBase     = "https://apidojo-yahoo-finance-v1.p.rapidapi.com/stock/v2"
	urlGetChart = urlBase + "/get-chart"
)

type Client interface {
	GetChart(ctx context.Context, params GetChartParams) (Chart, error)
}

func NewClient(httpClient *http.Client, apiKey, apiHost string) (Client, error) {
	validate := validator.New()
	if err := validate.Struct(struct {
		APIKey  string `validate:"required"`
		APIHost string `validate:"required"`
	}{
		APIKey:  apiKey,
		APIHost: apiHost,
	}); err != nil {
		return nil, err
	}

	return &client{
		apiKey:   apiKey,
		apiHost:  apiHost,
		validate: validate,
		http:     httpClient,
	}, nil
}

type client struct {
	apiKey   string
	apiHost  string
	validate *validator.Validate
	http     *http.Client
}

func (c *client) newRequest(ctx context.Context, method, url string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-rapidapi-key", c.apiKey)
	req.Header.Set("x-rapidapi-host", c.apiHost)
	return req.WithContext(ctx), nil
}
