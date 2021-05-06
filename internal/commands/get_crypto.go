package commands

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/jnericks/obibot/internal/clients/cmc"
	"github.com/olekukonko/tablewriter"
)

func GetCrypto(cmcClient cmc.Client, formatter func([]cmc.Cryptocurrency) (*Output, error)) Func {
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
		if resp == nil || len(resp.Data) == 0 {
			return nil, fmt.Errorf("no data for symbols %v", symbols)
		}

		return formatter(resp.Data)
	}
}

func formatPrice(v float64) string {
	if v < 10 {
		return "$" + humanize.FormatFloat("#,###.####", v)
	}
	return "$" + humanize.FormatFloat("#,###.##", v)
}

func formatPercent(v float64) string {
	return fmt.Sprintf("%.2f%%", v)
}

func FormatCryptoAsFlat(data []cmc.Cryptocurrency) (*Output, error) {
	var sb strings.Builder
	for i, c := range data {
		if i > 0 {
			sb.WriteByte('\n')
		}
		q := c.Quote.USD
		sb.WriteString(fmt.Sprintf("%s: %s (%s, %s, %s)",
			c.Symbol,
			formatPrice(q.Price),
			formatPercent(q.PercentChange1H),
			formatPercent(q.PercentChange24H),
			formatPercent(q.PercentChange7D),
		))
	}

	return &Output{
		Response: sb.String(),
		Markdown: false,
	}, nil
}

func FormatCryptoAsMarkdownTable(data []cmc.Cryptocurrency) (*Output, error) {
	var buf bytes.Buffer
	t := tablewriter.NewWriter(&buf)

	t.SetHeaderAlignment(tablewriter.ALIGN_RIGHT)
	t.SetHeader([]string{"", "Price", "1h", "24h", "7d"})

	t.SetColumnAlignment([]int{
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_RIGHT,
	})
	for _, c := range data {
		q := c.Quote.USD
		t.Append([]string{
			fmt.Sprintf("%s (%s)", c.Name, c.Symbol),
			formatPrice(q.Price),
			formatPercent(q.PercentChange1H),
			formatPercent(q.PercentChange24H),
			formatPercent(q.PercentChange7D),
		})
	}
	t.Render()
	return &Output{
		Response: "```\n" + buf.String() + "```",
		Markdown: true,
	}, nil
}
