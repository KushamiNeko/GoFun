package data

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/KushamiNeko/GoFun/Model/finance"
)

type quandl struct{}

func (q quandl) Read(from, to time.Time, symbol string, freq Frequency) (*finance.TimeSeries, error) {
	switch {
	case (strings.HasPrefix(symbol, "nk225") ||
		strings.HasPrefix(symbol, "nk300")):

		return q.readOSE(from, to, symbol, freq)

	case (strings.HasPrefix(symbol, "fesx") ||
		strings.HasPrefix(symbol, "fvs") ||
		strings.HasPrefix(symbol, "fgbl")):

		return q.readEUREX(from, to, symbol, freq)

	default:
		return q.readSCF(from, to, symbol, freq)
	}
}

func (q quandl) readOSE(from, to time.Time, symbol string, freq Frequency) (*finance.TimeSeries, error) {

	const root = "https://www.quandl.com/api/v3/datasets/OSE"

	url := fmt.Sprintf(
		"%s/%s.csv?api_key=%s",
		root,
		symbol,
		os.Getenv("QUANDL"),
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	//Date,Open,High,Low,Last,Change,Volume,Sett Price,Open Interest
	//2020-01-15,23950.0,23955.0,23910.0,23920.0,5.0,2203.0,23915.0,31278.0
	const pattern = `(\d{4}-\d{2}-\d{2})\s*,\s*([0-9.]+)\s*,\s*([0-9.]+)\s*,\s*([0-9.]+)\s*,\s*([0-9.]+)\s*,\s*([-0-9.]+)\s*,\s*([0-9.]+)\s*,\s*([0-9.]+)\s*,\s*([0-9.]+)`

	regex := regexp.MustCompile(pattern)

	match := regex.FindAllStringSubmatch(string(body), -1)
	ml := len(match)

	if ml <= 0 {
		return nil, fmt.Errorf("empty time series")
	}

	ts, vs := initTimeValues(len(match))

	for i, m := range match {

		d := m[1]
		o := m[2]
		h := m[3]
		l := m[4]
		c := m[5]
		v := m[7]
		oi := m[9]

		err = setTimeValuesRow("2006-01-02", ml-1-i, ts, vs, d, o, h, l, c, v, oi)
		if err != nil {
			return nil, err
		}
	}

	series := finance.NewTimeSeries(ts, vs)
	series.TimeSlice(from, to)

	return series, nil
}
func (q quandl) readEUREX(from, to time.Time, symbol string, freq Frequency) (*finance.TimeSeries, error) {

	const root = "https://www.quandl.com/api/v3/datasets/EUREX"

	url := fmt.Sprintf(
		"%s/%s.csv?api_key=%s",
		root,
		symbol,
		os.Getenv("QUANDL"),
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	//Date,Open,High,Low,Sett Price,Volume,Open Interest
	const pattern = `(\d{4}-\d{2}-\d{2})\s*,\s*([0-9.]+)\s*,\s*([0-9.]+)\s*,\s*([0-9.]+)\s*,\s*([0-9.]+)\s*,\s*([0-9.]+)\s*,\s*([0-9.]+)\s*`

	regex := regexp.MustCompile(pattern)

	match := regex.FindAllStringSubmatch(string(body), -1)
	ml := len(match)

	if ml <= 0 {
		return nil, fmt.Errorf("empty time series")
	}

	ts, vs := initTimeValues(len(match))

	for i, m := range match {

		d := m[1]
		o := m[2]
		h := m[3]
		l := m[4]
		c := m[5]
		v := m[6]
		oi := m[7]

		err = setTimeValuesRow("2006-01-02", ml-1-i, ts, vs, d, o, h, l, c, v, oi)
		if err != nil {
			return nil, err
		}
	}

	series := finance.NewTimeSeries(ts, vs)
	series.TimeSlice(from, to)

	return series, nil
}

func (q quandl) readSCF(from, to time.Time, symbol string, freq Frequency) (*finance.TimeSeries, error) {

	const root = "https://www.quandl.com/api/v3/datatables/SCF/PRICES.csv"

	url := fmt.Sprintf(
		"%s?api_key=%s&quandl_code=%s",
		root,
		os.Getenv("QUANDL"),
		//"CME_MD1_FW",
		//"CME_MD1_FR",
		"CME_MD1_OR",
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	//quandl_code,name,exchange,symbol,depth,method,date,open,high,low,settle,volume,prev_day_open_interest,front_contract
	//CME_MD1_OR,"CME S&P 400 Midcap Index Futures #1 (MD1) - Front Month - Backwards Ratio Adjusted Prices, Roll on Open Interest Switch",CME,MD,1,OR,2020-01-09,2057.7,2064.7,2052.6,2059.1,12414,76899,MDH2020
	const pattern = `\w+\s*,\s*\"[^"]+\"\s*,\s*\w+\s*,\s*\w+\s*,\s*\d\s*,\s*\w+\s*,\s*(\d{4}-\d{2}-\d{2})\s*,\s*([0-9.]+)\s*,\s*([0-9.]+)\s*,\s*([0-9.]+)\s*,\s*([0-9.]+)\s*,\s*([0-9]+)\s*,\s*([0-9]+)\s*,\s*\w+\s*`

	regex := regexp.MustCompile(pattern)

	match := regex.FindAllStringSubmatch(string(body), -1)
	ml := len(match)

	if ml <= 0 {
		return nil, fmt.Errorf("empty time series")
	}

	ts, vs := initTimeValues(len(match))

	for i, m := range match {

		d := m[1]
		o := m[2]
		h := m[3]
		l := m[4]
		c := m[5]
		v := m[6]
		oi := m[7]

		err = setTimeValuesRow("2006-01-02", ml-1-i, ts, vs, d, o, h, l, c, v, oi)
		if err != nil {
			return nil, err
		}
	}

	series := finance.NewTimeSeries(ts, vs)
	series.TimeSlice(from, to)

	return series, nil
}
