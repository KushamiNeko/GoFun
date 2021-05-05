package futures

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/KushamiNeko/GoFun/Model/finance"
)

type Contract struct {
	symbol string
	series *finance.TimeSeries

	format SymbolFormat
}

func NewContract(symbol string, series *finance.TimeSeries, format SymbolFormat) *Contract {
	return &Contract{
		symbol: symbol,
		series: series,
		format: format,
	}
}

func (c *Contract) ContractYear() int {
	var regex *regexp.Regexp

	switch c.format {
	case BarchartSymbolFormat:
		regex = regexp.MustCompile(barchartContractPattern)
	case QuandlSymbolFormat:
		regex = regexp.MustCompile(quandlContractPattern)
	default:
		panic("unknown symbol format")
	}

	match := regex.FindAllStringSubmatch(c.symbol, -1)
	if len(match) == 0 {
		panic("invalid symbol for specified symbol format")
	}

	var year string
	switch c.format {
	case BarchartSymbolFormat:
		year = fmt.Sprintf("20%s", match[0][3])
	case QuandlSymbolFormat:
		year = match[0][3]
	default:
		panic("unknown symbol format")

	}

	y, err := strconv.ParseInt(year, 10, 32)
	if err != nil {
		panic(err)
	}

	if int(y) > time.Now().Year() {
		y -= 100
	}

	return int(y)
}

func (c *Contract) ContractMonth() time.Month {
	var regex *regexp.Regexp

	switch c.format {
	case BarchartSymbolFormat:
		regex = regexp.MustCompile(barchartContractPattern)
	case QuandlSymbolFormat:
		regex = regexp.MustCompile(quandlContractPattern)
	default:
		panic("unknown symbol format")
	}

	match := regex.FindAllStringSubmatch(c.symbol, -1)
	if len(match) == 0 {
		panic("invalid symbol for specified symbol format")
	}

	return MonthCode(match[0][2]).Month()
}

func (c *Contract) Symbol() string {
	return c.symbol
}

func (c *Contract) Series() *finance.TimeSeries {
	return c.series
}
