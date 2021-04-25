package commands

import (
	"context"
	"errors"

	"github.com/jnericks/obibot/internal/clients/iex"
)

func GetCryptoPrice(iexClient iex.Client) Func {
	return func(ctx context.Context, input Input) (Output, error) {
		if len(input.Args) < 1 {
			return Output{}, errors.New("expecting symbol as input")
		}

		resp, err := iexClient.GetCryptoPrice(ctx, iex.GetCryptoPriceParams{
			Symbol: input.Args[0],
		})
		if err != nil {
			return Output{}, err
		}

		return Output{
			Response: resp.PriceSummary(),
		}, nil
	}
}
