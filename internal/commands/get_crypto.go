package commands

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/jnericks/obibot/internal/clients/cmc"
	"github.com/leekchan/accounting"
)

var defaultCryptoSymbols = []string{
	"BTC", "ETH", "DOGE",
}

func GetCrypto(cmcClient cmc.Client, formatter func(*cmc.GetCryptocurrencyQuotesResponse) (*Output, error)) Func {
	return func(ctx context.Context, input Input) (*Output, error) {
		var symbols []string
		for _, arg := range input.Args {
			for _, s := range strings.Split(strings.Replace(arg, " ", "", -1), ",") {
				if s == "" {
					continue
				}
				symbols = append(symbols, strings.ToUpper(s))
			}
		}

		if len(symbols) == 0 {
			symbols = defaultCryptoSymbols
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
		percent := c.Quote.USD.PercentChange24H
		change := price - (price / (1 + (percent / 100)))

		lc := accounting.LocaleInfo["USD"]
		ac := accounting.Accounting{
			Symbol:    lc.ComSymbol,
			Precision: 2,
			Thousand:  lc.ThouSep,
			Decimal:   lc.DecSep,
		}
		if math.Abs(price) < 2 {
			ac.Precision = 4
		}

		changePrefix := "+"
		if change < 0 {
			changePrefix = ""
		}

		sb.WriteString(fmt.Sprintf("%s: %s %s (%s%%)",
			c.Symbol,
			ac.FormatMoney(price),
			changePrefix+ac.FormatMoney(change),
			fmt.Sprintf("%+.2f", percent),
		))
	}

	return &Output{
		Response: sb.String(),
		Markdown: false,
	}, nil
}
