package data

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/KushamiNeko/GoFun/Chart/utils"
	"github.com/KushamiNeko/GoFun/Model/finance"
)

type stockCharts struct{}

func (s stockCharts) Read(from, to time.Time, symbol string, freq Frequency) (*finance.TimeSeries, error) {
	root := filepath.Join(os.Getenv("HOME"), "Documents/data_source/stockcharts")

	path := filepath.Join(root, fmt.Sprintf("%s.txt", symbol))

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("unknown symbol: %s", symbol)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(content)

	const regex string = `(\w+)\s*(\d{2}-\d{2}-\d{4})\s*([0-9.]+)\s*([0-9.]+)\s*([0-9.]+)\s*([0-9.]+)\s*([0-9.]+)\s*`

	re := regexp.MustCompile(regex)

	match := re.FindAllStringSubmatch(buffer.String(), -1)

	ml := len(match)

	if ml <= 0 {
		return nil, fmt.Errorf("empty time series")
	}

	ts, vs := initTimeValues(ml)

	const timeFormat = `01-02-2006`

	for i, m := range match {
		err = setTimeValuesRow(timeFormat, ml-1-i, ts, vs, m[2], m[3], m[4], m[5], m[6], m[7], "0")
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
