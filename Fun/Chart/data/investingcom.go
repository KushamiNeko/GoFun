package data

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/KushamiNeko/GoFun/Chart/utils"
	"github.com/KushamiNeko/GoFun/Model/finance"
)

type investingCom struct{}

func (s investingCom) Read(from, to time.Time, symbol string, freq Frequency) (*finance.TimeSeries, error) {
	root := filepath.Join(os.Getenv("HOME"), "Documents/data_source/investing.com")

	path := filepath.Join(root, fmt.Sprintf("%s.csv", symbol))

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("unknown symbol: %s", symbol)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(content)

	//"Date","Price","Open","High","Low","Vol.","Change %"
	//"Jan 07, 2020","3,744","3,748","3,772","3,736","-","0.13%"
	const regex string = `\"(\w{3}\s*\d{2}\s*,\s*\d{4})\"\s*,\s*\"([0-9.,]+)\"\s*,\s*\"([0-9.,]+)\"\s*,\s*\"([0-9.,]+)\"\s*,\s*\"([0-9.,]+)\"\s*,\s*\"([-0-9.,KM]+)\"\s*,\s*\"[-0-9.,%]+\"`

	re := regexp.MustCompile(regex)

	match := re.FindAllStringSubmatch(buffer.String(), -1)

	ml := len(match)

	if ml <= 0 {
		return nil, fmt.Errorf("empty time series")
	}

	ts, vs := initTimeValues(ml)

	const timeFormat = `Jan 02, 2006`

	for i, m := range match {
		t := m[1]
		o := strings.ReplaceAll(m[3], ",", "")
		h := strings.ReplaceAll(m[4], ",", "")
		l := strings.ReplaceAll(m[5], ",", "")
		c := strings.ReplaceAll(m[2], ",", "")

		v := m[6]
		if v == "" || v == "-" {
			v = "0"
		} else {
			f, err := strconv.ParseFloat(v[:len(v)-1], 10)
			if err != nil {
				return nil, err
			}

			switch v[len(v)-1:] {
			case "K":
				f = f * 1000
			case "M":
				f = f * 1000000
			}

			v = fmt.Sprintf("%f", f)
		}

		err = setTimeValuesRow(timeFormat, ml-1-i, ts, vs, t, o, h, l, c, v, "0")
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
