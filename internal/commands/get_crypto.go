package commands

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jnericks/obibot/internal/clients/cmc"
)

func GetCrypto(cmcClient cmc.Client, formatter func(*cmc.GetCryptocurrencyQuotesResponse) (*Output, error)) Func {
	return func(ctx context.Context, input Input) (*Output, error) {
		if len(input.Args) < 1 {
			return nil, errors.New("expecting crypto symbol as input")
		}

		symbols := make([]string, 0, len(input.Args))
		for _, a := range input.Args {
			a = strings.ToUpper(strings.Replace(a, " ", "", -1))
			if a == "" {
				continue
			}
			symbols = append(symbols, strings.Split(a, ",")...)
		}

		resp, err := cmcClient.GetCryptocurrencyQuotes(ctx, cmc.GetCryptocurrencyQuotesParams{
			Symbols: symbols,
		})
		if err != nil {
			return nil, err
		}

		if resp == nil {
			return nil, errors.New("server error")
		}
		if resp.Error != "" {
			return &Output{
				Response: resp.Error,
				Markdown: false,
			}, nil
		}

		return formatter(resp)
	}
}

func FormatGetCryptocurrencyQuotesResponse(resp *cmc.GetCryptocurrencyQuotesResponse) (*Output, error) {
	var sb strings.Builder
	for i, c := range resp.Quotes {
		if i > 0 {
			sb.WriteByte('\n')
		}

		price := c.Quote.USD.Price
		changePercent := c.Quote.USD.PercentChange24H
		change := price / (1 + changePercent)
		priceFormat := "%.2f"
		changeFormat := "%+.2f"
		changePercentFormat := "%+.2f%%"
		if -10 < price && price < 10 {
			priceFormat = "%.4f"
			changeFormat = "+%.4f"
		}

		sb.WriteString(fmt.Sprintf("%s: %s %s (%s)",
			c.Symbol,
			fmt.Sprintf(priceFormat, price),
			fmt.Sprintf(changeFormat, change),
			fmt.Sprintf(changePercentFormat, changePercent),
		))
	}

	return &Output{
		Response: sb.String(),
		Markdown: false,
	}, nil
}
