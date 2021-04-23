package commands

import (
	"context"
	"errors"

	"github.com/jnericks/obibot/internal/clients/rapid"
)

func GetStockPrice(rapidClient rapid.Client) Func {
	return func(ctx context.Context, input Input) (Output, error) {
		if len(input.Args) < 1 {
			return Output{}, errors.New("expecting symbol as input")
		}

		chart, err := rapidClient.GetChart(ctx, rapid.GetChartParams{
			Symbol: input.Args[0],
		})
		if err != nil {
			return Output{}, err
		}

		return Output{
			Response: chart.Meta.PriceSummary(),
		}, nil
	}
}
