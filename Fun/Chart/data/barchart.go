package data

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/KushamiNeko/GoFun/Chart/utils"
	"github.com/KushamiNeko/GoFun/Model/finance"
)

type barchart struct{}

func (b barchart) Read(from, to time.Time, symbol string, freq Frequency) (*finance.TimeSeries, error) {
	//if symbol == "np" || symbol == "fx" || symbol == "vi" || symbol == "dv" {
	//return b.readPremierInteractive(from, to, symbol, freq)
	//} else {
	//return b.readOnDemand(from, to, symbol, freq)
	//}
	//var series *finance.TimeSeries
	//var err error

	//series, err = b.readPremierHistorical(from, to, symbol, freq)
	//if err == nil {
	//goto re
	//}

	//series, err = b.readPremierInteractive(from, to, symbol, freq)
	//if err == nil {
	//goto re
	//}

	//series, err = b.readOnDemand(from, to, symbol, freq)
	//if err == nil {
	//goto re
	//}

	//return nil, fmt.Errorf("invalid barchart data format\n")

	//re:
	//return series, nil

	root := filepath.Join(os.Getenv("HOME"), "Documents/data_source/barchart")

	//path := filepath.Join(root, fmt.Sprintf("%s_%s.csv", symbol, freq))
	path := filepath.Join(root, fmt.Sprintf("%s.csv", symbol))

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("unknown symbol: %s", symbol)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(content)

	bl := len(strings.Split(buffer.String(), "\n"))
	if bl <= 2 {
		return nil, errEmptyFile
	}

	reader := csv.NewReader(buffer)
	match := make([][]string, 0, bl)
	for {
		record, err := reader.Read()
		if errors.Is(err, csv.ErrFieldCount) || errors.Is(err, io.EOF) {
			break
		}

		match = append(match, record)
	}

	ml := len(match)
	if ml <= 2 {
		return nil, errDataFormat
	}

	match = match[1:]
	ml -= 1

	ts, vs := initTimeValues(ml)

	for i, m := range match {
		err := setTimeValuesRow(`01/02/2006`, ml-1-i, ts, vs, m[0], m[1], m[2], m[3], m[4], m[6], "0")
		if err != nil {
			err := setTimeValuesRow(`01/02/06`, ml-1-i, ts, vs, m[0], m[1], m[2], m[3], m[4], m[6], "0")
			if err != nil {
				return nil, err
			}
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

	return series, nil
}

func (b barchart) readPremierInteractive(from, to time.Time, symbol string, freq Frequency) (*finance.TimeSeries, error) {

	root := filepath.Join(os.Getenv("HOME"), "Documents/data_source/barchart")

	//path := filepath.Join(root, fmt.Sprintf("%s_%s.csv", symbol, freq))
	path := filepath.Join(root, fmt.Sprintf("%s.csv", symbol))

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("unknown symbol: %s", symbol)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(content)
	c := continuousContract{}
	ts, err := c.barchartInteractive(buffer)
	if err != nil {
		return nil, err
	}

	switch freq {
	case Hourly:
		panic("unimplemented")
	case Daily:
	case Weekly:
		ts = utils.DailyToWeekly(ts)
	case Monthly:
		panic("unimplemented")
	default:
		panic("unknown frequency")
	}

	return ts, nil

	//"Date Time","Symbol","Open","High","Low","Close","Change"
	//2000-01-03,ESH00,1549.5,1556.5,1512.75,1527.25,0
	//const regex string = `\s*(\d{4}-\d{2}-\d{2})\s*,\s*([\w\d]+)\s*,\s*([0-9.]+)\s*,\s*([0-9.]+)\s*,\s*([0-9.]+)\s*,\s*([0-9.]+)\s*,\s*([-0-9.]+)\s*`

	//re := regexp.MustCompile(regex)

	//match := re.FindAllStringSubmatch(buffer.String(), -1)

	//ml := len(match)

	//if ml <= 0 {
	//return nil, fmt.Errorf("empty time series")
	//}

	//ts, vs := initTimeValues(ml)

	//const timeFormat = `2006-01-02`

	//for i, m := range match {
	//err = setTimeValuesRow(timeFormat, i, ts, vs, m[1], m[3], m[4], m[5], m[6], "", "")
	//if err != nil {
	//return nil, err
	//}
	//}

	//series := finance.NewTimeSeries(ts, vs)
	//series.TimeSlice(from, to)

	//return series, nil
}

func (b barchart) readPremierHistorical(from, to time.Time, symbol string, freq Frequency) (*finance.TimeSeries, error) {
	root := filepath.Join(os.Getenv("HOME"), "Documents/data_source/barchart")

	//path := filepath.Join(root, fmt.Sprintf("%s_%s.csv", symbol, freq))
	path := filepath.Join(root, fmt.Sprintf("%s.csv", symbol))

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("unknown symbol: %s", symbol)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(content)

	c := continuousContract{}
	ts, err := c.barchartHistorical(buffer)
	if err != nil {
		return nil, err
	}

	switch freq {
	case Hourly:
		panic("unimplemented")
	case Daily:
	case Weekly:
		ts = utils.DailyToWeekly(ts)
	case Monthly:
		panic("unimplemented")
	default:
		panic("unknown frequency")
	}

	return ts, nil

	//bl := len(strings.Split(buffer.String(), "\n"))
	//if bl <= 2 {
	//return nil, errEmptyFile
	//}

	//reader := csv.NewReader(buffer)
	//match := make([][]string, 0, bl)
	//for {
	//record, err := reader.Read()
	//if errors.Is(err, csv.ErrFieldCount) || errors.Is(err, io.EOF) {
	//break
	//}

	//match = append(match, record)
	//}

	//ml := len(match)
	//if ml <= 2 {
	//return nil, errDataFormat
	//}

	//match = match[1:]
	//ml -= 1

	//ts, vs := initTimeValues(ml)

	//for i, m := range match {
	////err := setTimeValuesRow(timeFormat, ml-1-i, ts, vs, m[1], m[2], m[3], m[4], m[5], m[7], m[8])
	//err := setTimeValuesRow(`01/02/2006`, ml-1-i, ts, vs, m[0], m[1], m[2], m[3], m[4], m[6], m[7])
	//if err != nil {
	//err := setTimeValuesRow(`01/02/06`, ml-1-i, ts, vs, m[0], m[1], m[2], m[3], m[4], m[6], m[7])
	//if err != nil {
	//return nil, err
	//}
	//}
	//}

	//series := finance.NewTimeSeries(ts, vs)

	//Symbol,Time,Open,High,Low,Last,Change,Volume,"Open Int"
	//const regex string = `(?:\s*([a-zA-Z0-9]+)\s*,)*\s*(\d{2}/\d{2}/\d{2})\s*,\s*([0-9.]+)\s*,\s*([0-9.]+)\s*,\s*([0-9.]+)\s*,\s*([0-9.]+)\s*,\s*([-0-9.]+)\s*,\s*([0-9]+)\s*(?:,\s*([0-9]+))*\s*`

	//re := regexp.MustCompile(regex)

	//match := re.FindAllStringSubmatch(buffer.String(), -1)

	//ts, vs := initTimeValues(len(match))

	//ml := len(match)

	//if ml <= 0 {
	//return nil, fmt.Errorf("empty time series")
	//}

	//const timeFormat = `01/02/06`

	//for i, m := range match {
	//err = setTimeValuesRow(timeFormat, ml-1-i, ts, vs, m[2], m[3], m[4], m[5], m[6], m[8], m[9])
	//if err != nil {
	//return nil, err
	//}
	//}

	//series := finance.NewTimeSeries(ts, vs)
	//series.TimeSlice(from, to)

	//return series, nil
}

func (b barchart) readOnDemand(from, to time.Time, symbol string, freq Frequency) (*finance.TimeSeries, error) {

	const root = `https://ondemand.websol.barchart.com/getHistory.csv?`

	const timeFormat = "20060102"

	d := b.preload(freq)
	s := from.Add(-d)
	e := to

	url := fmt.Sprintf(
		"%sapikey=%s&symbol=%s&type=%s&startDate=%s&endDate=%s&interval=%s&volume=%s&backAdjust=%s&contractRoll=%s",
		root,
		os.Getenv("BARCHART"),
		symbol,
		//b.frequency(freq),
		b.frequency(Daily),
		s.Format(timeFormat),
		e.Format(timeFormat),
		b.interval(freq),
		"total",
		"true",
		"combined",
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

	c := continuousContract{}
	ts, err := c.barchartOnDemand(buffer)
	if err != nil {
		return nil, err
	}

	switch freq {
	case Hourly:
		panic("unimplemented")
	case Daily:
	case Weekly:
		ts = utils.DailyToWeekly(ts)
	case Monthly:
		panic("unimplemented")
	default:
		panic("unknown frequency")
	}

	return ts, nil

	//"esu19","2018-01-02T00:00:00-06:00","2018-01-02","2692.5","2713.25","2691.75","2710.25","1003194","3051095"
	//const pattern = `"([a-zA-Z0-9]+)","(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}-\d{2}:\d{2})","(\d{4}-\d{2}-\d{2})","([0-9.]+)","([0-9.]+)","([0-9.]+)","([0-9.]+)","(\d+)"(?:,"(\d+)")*`

	//regex := regexp.MustCompile(pattern)

	//match := regex.FindAllStringSubmatch(string(body), -1)

	//ml := len(match)

	//if ml <= 0 {
	//return nil, fmt.Errorf("empty time series")
	//}

	//ts, vs := initTimeValues(ml)

	//for i, m := range match {

	//d := m[2]
	//o := m[4]
	//h := m[5]
	//l := m[6]
	//c := m[7]
	//v := m[8]
	//oi := m[9]

	//regex := regexp.MustCompile(`(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2})-\d{2}:\d{2}`)
	//d = regex.FindAllStringSubmatch(d, -1)[0][1]

	//if oi == "" {
	//oi = "0"
	//}

	//err = setTimeValuesRow(`2006-01-02T15:04:05`, i, ts, vs, d, o, h, l, c, v, oi)
	//if err != nil {
	//return nil, err
	//}
	//}

	//series := finance.NewTimeSeries(ts, vs)
	//series.TimeSlice(from, to)

	//return series, nil
}

func (b barchart) frequency(freq Frequency) string {
	switch freq {
	case Hourly:
		return "nearbyMinutes"
	case Daily:
		return "dailyNearest"
	case Weekly:
		return "weeklyNearest"
	case Monthly:
		return "monthlyNearest"
	default:
		panic(fmt.Sprintf("unknown frequency: %s", freq))
	}
}

func (b barchart) interval(freq Frequency) string {
	switch freq {
	case Hourly:
		return "60"
	default:
		return "1"
	}
}

func (b barchart) preload(freq Frequency) time.Duration {

	const preload = 30

	switch freq {
	case Hourly:
		return 3 * 24 * time.Hour
	case Daily:
		return preload * 24 * time.Hour
	case Weekly:
		return preload * 7 * 24 * time.Hour
	case Monthly:
		return preload * 30 * 24 * time.Hour
	default:
		panic(fmt.Sprintf("unknown frequency: %s", freq))
	}
}
