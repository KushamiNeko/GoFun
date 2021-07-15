package data

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/KushamiNeko/GoFun/Fun/Chart/utils"
	"github.com/KushamiNeko/GoFun/Fun/Model/finance"
)

type alphaVantage struct{}

func (a alphaVantage) Read(from, to time.Time, symbol string, freq Frequency) (*finance.TimeSeries, error) {

	//https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=MSFT&outputsize=full&apikey=demo&datatype=csv
	const root = `https://www.alphavantage.co/query?`

	url := fmt.Sprintf(
		"%sapikey=%s&symbol=%s&function=%s&interval=60min&outputsize=full&datatype=csv",
		root,
		os.Getenv("ALPHA_VANTAGE"),
		symbol,
		a.frequency(Daily),
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

	buffer := bytes.NewBuffer(body)

	//timestamp,open,high,low,close,volume
	//2019-09-05,2960.6001,2985.8601,2960.6001,2977.4299,769816370
	//2019-09-05 11:50:17

	const pattern = `(\d{4}-\d{2}-\d{2}(?:\s*\d{2}:\d{2}:\d{2})*),([0-9.]+),([0-9.]+),([0-9.]+),([0-9.]+),(\d+)`

	regex := regexp.MustCompile(pattern)

	match := regex.FindAllStringSubmatch(buffer.String(), -1)

	ml := len(match)

	if len(strings.Split(buffer.String(), "\n")) == 0 || ml <= 0 {
		return nil, fmt.Errorf("empty time series")
	}

	ts, vs := initTimeValues(ml)

	for i, m := range match {

		timeFormat := `2006-01-02 15:04:05`
		regex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
		if regex.MatchString(m[1]) {
			timeFormat = `2006-01-02`
		}

		err = setTimeValuesRow(timeFormat, ml-1-i, ts, vs, m[1], m[2], m[3], m[4], m[5], m[6], "0")
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

func (a alphaVantage) frequency(freq Frequency) string {

	switch freq {
	case Hourly:
		return "TIME_SERIES_INTRADAY"
	case Daily:
		return "TIME_SERIES_DAILY"
	case Weekly:
		return "TIME_SERIES_WEEKLY"
	case Monthly:
		return "TIME_SERIES_MONTHLY"
	default:
		panic(fmt.Sprintf("unknown frequency: %s", freq))
	}

}
