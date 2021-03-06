package commands

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/jnericks/obibot/internal/clients/iex"
	"github.com/leekchan/accounting"
)

var defaultStockSymbols = []string{
	"SPY", "QQQ", "BND",
}

func GetStock(iexClient iex.Client, formatter func(*iex.GetStockQuotesResponse) (*Output, error)) Func {
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
			symbols = defaultStockSymbols
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

		price := s.LatestPrice
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

		change := s.Change
		changePrefix := "+"
		if change < 0 {
			changePrefix = ""
		}

		sb.WriteString(fmt.Sprintf("%s: %s %s (%s%%)",
			s.Symbol,
			ac.FormatMoney(price),
			changePrefix+ac.FormatMoney(change),
			fmt.Sprintf("%+.2f", s.ChangePercent*100),
		))
	}

	return &Output{
		Response: sb.String(),
		Markdown: false,
	}, nil
}
