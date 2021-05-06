package commands

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/jnericks/obibot/internal/clients/cmc"
	"github.com/olekukonko/tablewriter"
)

func GetCrypto(cmcClient cmc.Client) Func {
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

		sort.Strings(symbols)
		resp, err := cmcClient.GetLatestQuote(ctx, cmc.GetLatestQuoteParams{
			Symbols: symbols,
		})
		if err != nil {
			return nil, err
		}
		if resp == nil || len(resp.Data) == 0 {
			return nil, fmt.Errorf("no data for symbols %v", symbols)
		}

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
		for _, c := range resp.Data {
			var prc string
			if p := c.Quote.USD.Price; p < 10 {
				prc = fmt.Sprintf("%.4f", p)
			} else {
				prc = fmt.Sprintf("%.2f", p)
			}
			t.Append([]string{
				fmt.Sprintf("%s (%s)", c.Name, c.Symbol),
				prc,
				fmt.Sprintf("%+.2f%%", c.Quote.USD.PercentChange1H),
				fmt.Sprintf("%+.2f%%", c.Quote.USD.PercentChange24H),
				fmt.Sprintf("%+.2f%%", c.Quote.USD.PercentChange7D),
			})
		}
		t.Render()

		return &Output{
			Response: "```\n" + buf.String() + "```",
			Markdown: true,
		}, nil
	}
}
