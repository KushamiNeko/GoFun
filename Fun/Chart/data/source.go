package data

import (
	"math"
	"strconv"
	"time"

	"github.com/KushamiNeko/GoFun/Model/finance"
)

type Source int

const (
	Yahoo Source = iota
	StockCharts

	InvestingCom

	Barchart
	Quandl

	AlphaVantage

	TradeStation

	Continuous
)

type Frequency int

const (
	Hourly Frequency = iota
	Daily
	Weekly
	Monthly
)

func (f Frequency) String() string {
	switch f {
	case Hourly:
		return "h"
	case Daily:
		return "d"
	case Weekly:
		return "w"
	case Monthly:
		return "m"
	default:
		panic("unknown frequency")
	}
}

func ParseFrequency(freq string) Frequency {
	switch freq {
	case "h":
		return Hourly
	case "d":
		return Daily
	case "w":
		return Weekly
	case "m":
		return Monthly
	default:
		panic("unknown frequency")
	}
}

func initTimeValues(length int) ([]time.Time, map[string][]float64) {
	times := make([]time.Time, length)

	values := make(map[string][]float64)

	values["open"] = make([]float64, length)
	values["high"] = make([]float64, length)
	values["low"] = make([]float64, length)
	values["close"] = make([]float64, length)
	values["volume"] = make([]float64, length)
	values["openinterest"] = make([]float64, length)

	return times, values
}

func setTimeValuesRow(timeFormat string, index int, ts []time.Time, vs map[string][]float64, t, o, h, l, c, v, oi string) error {
	var err error

	ts[index], err = time.Parse(timeFormat, t)
	if err != nil {
		return err
	}

	if o == "null" {
		vs["open"][index] = math.NaN()
	} else {
		vs["open"][index], err = strconv.ParseFloat(o, 64)
		if err != nil {
			return err
		}
	}

	if h == "null" {
		vs["high"][index] = math.NaN()
	} else {
		vs["high"][index], err = strconv.ParseFloat(h, 64)
		if err != nil {
			return err
		}
	}

	if l == "null" {
		vs["low"][index] = math.NaN()
	} else {
		vs["low"][index], err = strconv.ParseFloat(l, 64)
		if err != nil {
			return err
		}
	}

	if c == "null" {
		vs["close"][index] = math.NaN()
	} else {
		vs["close"][index], err = strconv.ParseFloat(c, 64)
		if err != nil {
			return err
		}
	}

	if v != "" && v != "0" && v != "null" {
		vs["volume"][index], err = strconv.ParseFloat(v, 64)
		if err != nil {
			return err
		}
	}

	if oi != "" && oi != "0" && o != "null" {
		vs["openinterest"][index], err = strconv.ParseFloat(oi, 64)
		if err != nil {
			return err
		}
	}

	return nil
}

type DataSource interface {
	Read(from, to time.Time, symbol string, freq Frequency) (*finance.TimeSeries, error)
}

func NewDataSource(src Source) DataSource {
	switch src {
	case StockCharts:
		return stockCharts{}
	case Yahoo:
		return yahooFinance{}
	case Barchart:
		return barchart{}
	case Quandl:
		return quandl{}
	case InvestingCom:
		return investingCom{}
	case AlphaVantage:
		return alphaVantage{}
	case Continuous:
		return continuousContract{}
	default:
		panic("unknown data source")
	}
}
