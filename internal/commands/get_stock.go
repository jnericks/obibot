package commands

import (
	"context"
	"errors"

	"github.com/jnericks/obibot/internal/clients/iex"
)

func GetStock(iexClient iex.Client) Func {
	return func(ctx context.Context, input Input) (*Output, error) {
		if len(input.Args) < 1 {
			return nil, errors.New("expecting stock symbol as input")
		}

		resp, err := iexClient.GetStock(ctx, iex.GetStockParams{
			Symbol: input.Args[0],
		})
		if err != nil {
			if ierr, ok := err.(iex.ErrInvalidSymbol); ok {
				return &Output{
					Response: ierr.Error(),
				}, nil
			}
			return nil, err
		}

		return &Output{
			Response: resp.PriceSummary(),
		}, nil
	}
}
