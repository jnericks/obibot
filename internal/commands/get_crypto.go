package commands

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/jnericks/obibot/internal/clients/cmc"
)

func GetCrypto(cmcClient cmc.Client, formatter func(*cmc.GetLatestQuoteResponse) (*Output, error)) Func {
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

		resp, err := cmcClient.GetLatestQuote(ctx, cmc.GetLatestQuoteParams{
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
		if len(resp.Data) == 0 {
			return nil, fmt.Errorf("no data for symbols %v", symbols)
		}

		return formatter(resp)
	}
}

func FormatCryptoAsFlat(resp *cmc.GetLatestQuoteResponse) (*Output, error) {
	var sb strings.Builder
	for i, c := range resp.Data {
		if i > 0 {
			sb.WriteByte('\n')
		}

		price := c.Quote.USD.Price
		percentChange := c.Quote.USD.PercentChange24H
		change := price / (1 + percentChange)
		format := "#,###.##"
		if -10 < price && price < 10 {
			format = "#,###.####"
		}

		sPrice := "$" + humanize.FormatFloat(format, price)
		sPercentChange := fmt.Sprintf("%+.2f%%", percentChange)
		sChange := humanize.FormatFloat(format, change)
		if change < 0 {
			sChange = "-$" + sChange
		} else {
			sChange = "+$" + sChange
		}

		sb.WriteString(fmt.Sprintf("%s: %s %s (%s)",
			c.Symbol,
			sPrice,
			sChange,
			sPercentChange,
		))
	}

	return &Output{
		Response: sb.String(),
		Markdown: false,
	}, nil
}
