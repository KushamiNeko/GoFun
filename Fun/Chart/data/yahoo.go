package data

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/KushamiNeko/GoFun/Chart/utils"
	"github.com/KushamiNeko/GoFun/Model/finance"
)

type yahooFinance struct{}

func (y yahooFinance) Read(from, to time.Time, symbol string, freq Frequency) (*finance.TimeSeries, error) {

	root := filepath.Join(os.Getenv("HOME"), "Documents/data_source/yahoo")

	path := filepath.Join(root, fmt.Sprintf("%s.csv", symbol))

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("unknown symbol: %s", symbol)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(content)

	reader := csv.NewReader(buffer)
	rs, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(rs) <= 1 {
		return nil, fmt.Errorf("empty time series")
	}

	ts, vs := initTimeValues(len(rs[1:]))

	const timeFormat = `2006-01-02`

	for i, m := range rs[1:] {
		err = setTimeValuesRow(timeFormat, i, ts, vs, m[0], m[1], m[2], m[3], m[4], m[6], "0")
		if err != nil {
			return nil, err
		}
	}

	series := finance.NewTimeSeries(ts, vs)

	switch freq {
	case Hourly:
		panic("unimplemented")
	case Daily:
	case Weekly:
		series = utils.DailyToWeekly(series)
	case Monthly:
		panic("unimplemented")
	default:
		panic("unknown frequency")
	}

	series.TimeSlice(from, to)

	return series, nil
}

func (y yahooFinance) frequency(freq Frequency) string {
	switch freq {
	case Hourly:
		return "1d"
	case Daily:
		return "1d"
	case Weekly:
		return "1wk"
	case Monthly:
		return "1mo"
	default:
		panic(fmt.Sprintf("unknown frequency: %s", freq))
	}
}
