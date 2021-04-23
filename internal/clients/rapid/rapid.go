package rapid

import "fmt"

type APIError struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

type Chart struct {
	Meta       ChartMeta       `json:"meta"`
	Timestamp  []int           `json:"timestamp"`
	Indicators ChartIndicators `json:"indicators"`
}

type ChartMeta struct {
	Currency             string        `json:"currency"`
	Symbol               string        `json:"symbol"`
	ExchangeName         string        `json:"exchangeName"`
	InstrumentType       string        `json:"instrumentType"`
	FirstTradeDate       int           `json:"firstTradeDate"`
	RegularMarketTime    int           `json:"regularMarketTime"`
	GMTOffset            int           `json:"gmtoffset"`
	Timezone             string        `json:"timezone"`
	ExchangeTimezoneName string        `json:"exchangeTimezoneName"`
	RegularMarketPrice   float64       `json:"regularMarketPrice"`
	ChartPreviousClose   float64       `json:"chartPreviousClose"`
	PriceHint            int           `json:"priceHint"`
	CurrentTradingPeriod TradingPeriod `json:"currentTradingPeriod"`
	DataGranularity      string        `json:"dataGranularity"`
	Range                string        `json:"range"`
	ValidRanges          []string      `json:"validRanges"`
}

func (m ChartMeta) PriceChangeAmount() float64 {
	return m.RegularMarketPrice - m.ChartPreviousClose
}

func (m ChartMeta) PriceChantPercent() float64 {
	return m.PriceChangeAmount() / m.ChartPreviousClose
}

func (m ChartMeta) PriceSummary() string {
	return fmt.Sprintf("%s: %.2f %+.2f (%+.2f%%)", m.Symbol, m.RegularMarketPrice, m.PriceChangeAmount(), m.PriceChantPercent()*100)
}

type TimePeriod struct {
	Timezone  string `json:"timezone"`
	Start     int    `json:"start"`
	End       int    `json:"end"`
	GMTOffset int    `json:"gmtoffset"`
}

type TradingPeriod struct {
	Pre     TimePeriod `json:"pre"`
	Regular TimePeriod `json:"regular"`
	Post    TimePeriod `json:"post"`
}

type ChartQuote struct {
	High   []float64 `json:"high"`
	Close  []float64 `json:"close"`
	Low    []float64 `json:"low"`
	Volume []int     `json:"volume"`
	Open   []float64 `json:"open"`
}

type ChartAdjustedClose struct {
	Values []float64 `json:"adjclose"`
}

type ChartIndicators struct {
	Quote          []ChartQuote         `json:"quote"`
	AdjustedCloses []ChartAdjustedClose `json:"adjclose"`
}
