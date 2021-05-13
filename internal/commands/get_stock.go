package commands

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jnericks/obibot/internal/clients/iex"
)

func GetStock(iexClient iex.Client, formatter func(*iex.GetStockQuotesResponse) (*Output, error)) Func {
	return func(ctx context.Context, input Input) (*Output, error) {
		if len(input.Args) < 1 {
			return nil, errors.New("expecting stock symbol as input")
		}

		symbols := make([]string, 0, len(input.Args))
		for _, a := range input.Args {
			a = strings.ToUpper(strings.Replace(a, " ", "", -1))
			if a == "" {
				continue
			}
			symbols = append(symbols, strings.Split(a, ",")...)
		}

		resp, err := iexClient.GetStockQuotes(ctx, iex.GetStockQuotesParams{
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

func FormatGetStockQuotesResponse(resp *iex.GetStockQuotesResponse) (*Output, error) {
	var sb strings.Builder
	for i, s := range resp.Quotes {
		if i > 0 {
			sb.WriteByte('\n')
		}

		priceFormat := "%.2f"
		changeFormat := "%+.2f"
		changePercentFormat := "%+.2f%%"
		if -10 < s.LatestPrice && s.LatestPrice < 10 {
			priceFormat = "%.4f"
			changeFormat = "+%.4f"
		}

		sb.WriteString(fmt.Sprintf("%s: %s %s (%s)",
			s.Symbol,
			fmt.Sprintf(priceFormat, s.LatestPrice),
			fmt.Sprintf(changeFormat, s.Change),
			fmt.Sprintf(changePercentFormat, s.ChangePercent),
		))
	}

	return &Output{
		Response: sb.String(),
		Markdown: false,
	}, nil
}
