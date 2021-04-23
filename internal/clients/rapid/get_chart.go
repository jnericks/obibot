package rapid

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type GetChartParams struct {
	Symbol string `validate:"required"`
}

func (c *client) GetChart(ctx context.Context, params GetChartParams) (Chart, error) {
	if err := c.validate.Struct(params); err != nil {
		return Chart{}, err
	}

	req, err := c.newRequest(ctx, http.MethodGet, urlGetChart)
	if err != nil {
		return Chart{}, err
	}
	q := req.URL.Query()
	q.Set("symbol", params.Symbol)
	req.URL.RawQuery = q.Encode()

	resp, err := c.http.Do(req)
	if err != nil {
		return Chart{}, err
	}

	var r struct {
		Chart struct {
			Results []Chart   `json:"result"`
			Error   *APIError `json:"error,omitempty"`
		} `json:"chart"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return Chart{}, err
	}
	if apiError := r.Chart.Error; apiError != nil {
		return Chart{}, fmt.Errorf("api error: (%s) %s", apiError.Code, apiError.Description)
	}
	if len(r.Chart.Results) == 0 {
		return Chart{}, fmt.Errorf("no charts found")
	}

	return r.Chart.Results[0], nil
}
